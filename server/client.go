package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Datas struct {
	room          string
	status        string
	inventory     []string
	invitation    []string
	group         string
	hp            int
	max_hp        int
	connected     bool
	last_cmd_time time.Time
	spam_warning  int
}

type Client struct {
	conn  net.Conn
	ch    chan Response
	ip    string
	name  string
	datas Datas
}

func handleClient(conn net.Conn) {
	responses := make(chan Response)
	go clientWriter(conn, responses)

	who := conn.RemoteAddr().String()

	// Allocation sur le tas (Heap) via le symbole '&' pour figer l'instance en mémoire
	cli := &Client{
		conn:  conn,
		ch:    responses,
		ip:    who,
		name:  "",
		datas: Datas{"start", "healthy", []string{}, []string{}, "", 100, 100, false, time.Now(), 0},
	}

	cli.ch <- Response{"OK hello proto=1", "", Request{}}
	entering <- cli

	input := bufio.NewScanner(conn)
	for input.Scan() {
		now := time.Now()
		if now.Sub(cli.datas.last_cmd_time) < 500*time.Millisecond {
			cli.datas.spam_warning++
			LogWarn("Abuse detected: Command flooding", map[string]any{
				"ip":       cli.ip,
				"player":   cli.name,
				"warnings": cli.datas.spam_warning,
			})
			if cli.datas.spam_warning > 3 {
				LogWarn("Client disconnected due to spam", map[string]any{
					"ip":       cli.ip,
					"player":   cli.name,
					"warnings": cli.datas.spam_warning,
				})
				fmt.Fprintln(conn, "ERR 900 CONNECTION_CLOSED_DUE_TO_SPAM")
				break
			}
			continue
		}
		cli.datas.last_cmd_time = now
		cli.datas.spam_warning = 0
		requests <- Request{cli, input.Text()}
	}

	leaving <- cli
	conn.Close()
}

func clientWriter(conn net.Conn, responses <-chan Response) {
	// Write all the messages in the player terminal
	for res := range responses {
		fmt.Fprint(conn, res.msg)

		// Handle QUIT command
		if res.msg == "OK bye" {
			break
		}

		// Handle json datas
		if res.datas != "" {
			jsonBytes, err := json.Marshal(res.datas)
			if err != nil {
				fmt.Fprint(conn, " ERR Internal server error during JSON parsing")
			} else {
				fmt.Fprint(conn, " "+string(jsonBytes))
			}
		}

		fmt.Fprint(conn, "\n")
	}
}
