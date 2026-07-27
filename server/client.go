package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	pr "tap/protocol"
	"time"
)

func handleClient(conn net.Conn) {
	// Send server responses to client terminal
	responses := make(chan pr.Response)
	go clientWriter(conn, responses)

	// Init client
	who := conn.RemoteAddr().String()

	cli := &pr.Client{
		Conn: conn,
		Ch:   responses,
		Ip:   who,
		Datas: pr.Datas{
			Room:          "entrance",
			Status:        "healthy",
			Promotion:     false,
			Hp:            100,
			Max_hp:        100,
			Connected:     false,
			Last_cmd_time: time.Now(),
		},
	}

	// Log player connection
	cli.Ch <- pr.Response{Msg: "OK hello proto=1"}
	entering <- cli

	// Handle player commands user buffers
	input := bufio.NewScanner(conn)
	buf := make([]byte, 0, 1024)
	input.Buffer(buf, 1024)

	for input.Scan() {
		// Handle command spams
		now := time.Now()
		if now.Sub(cli.Datas.Last_cmd_time) >= 1000*time.Millisecond {
			cli.Datas.Spam_warning = 0
			cli.Datas.Last_cmd_time = now
			requests <- pr.Request{Cli: cli, Msg: input.Text()}

			// Handle valid commands
		} else {
			cli.Datas.Spam_warning++
			LogWarn("Abuse detected", map[string]any{"warnings": cli.Datas.Spam_warning, "ip": cli.Ip, "player": cli.Name})

			if cli.Datas.Spam_warning > 3 {
				fmt.Fprintln(conn, "ERR 900 CONNECTION_CLOSED_DUE_TO_SPAM")
				cli.Ch <- pr.Response{Msg: "OK bye"}
				break
			}
		}
	}

	leaving <- cli
}

func clientWriter(conn net.Conn, responses <-chan pr.Response) {
	encoder := json.NewEncoder(conn)

	for res := range responses {
		if err := encoder.Encode(res); err != nil {
			return
		}
	}
}
