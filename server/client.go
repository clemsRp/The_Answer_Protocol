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
	buf := make([]byte, 0, 1024)
	input.Buffer(buf, 1024)
	for input.Scan() {
		now := time.Now()
		if now.Sub(cli.datas.last_cmd_time) >= 1000*time.Millisecond {
			cli.datas.spam_warning = 0
			cli.datas.last_cmd_time = now
			requests <- Request{cli, input.Text()}
		} else {
			cli.datas.spam_warning++
			LogWarn("Abuse detected", map[string]any{"warnings": cli.datas.spam_warning, "ip": cli.ip, "player": cli.name})

			if cli.datas.spam_warning > 3 {
				fmt.Fprintln(conn, "ERR 900 CONNECTION_CLOSED_DUE_TO_SPAM")
				cli.ch <- Response{"OK bye", Datas{}, Request{}}
				break
			}
		}
	}
	leaving <- cli
}

func clientWriter(conn net.Conn, responses <-chan Response) {
	encoder := json.NewEncoder(conn)

	for res := range responses {
		if err := encoder.Encode(res); err != nil {
			return
		}
	}
}
