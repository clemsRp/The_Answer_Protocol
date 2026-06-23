package main

import (
	"bufio"
	"fmt"
	"net"
	"encoding/json"
)

type Datas struct {
	room      string
	inventory []string
	group     string
	hp        int
	max_hp    int
	status    string
}

type Client struct {
	conn      net.Conn
	ch        chan Response
	ip        string
	name      string
	connected bool
	datas     Datas
}

func handleClient(conn net.Conn) {
	responses := make(chan Response)
	go clientWriter(conn, responses)

	// Init player
	who := conn.RemoteAddr().String()
	cli := Client{conn, responses, who, "", false, Datas{"start", []string{}, "", 50, 50, "healthy"}}

	// Start messages
	cli.ch <- Response{"[INFO]: You are connected as " + who, "", Request{}}
	entering <- cli

	// Handle the input commands
	input := bufio.NewScanner(conn)
	for input.Scan() {
		requests <- Request{cli, input.Text()}
	}

	// End the player's session
	leaving <- cli
	conn.Close()
}

func clientWriter(conn net.Conn, responses <-chan Response) {
	// Write all the messages in the player terminal
	for res := range responses {
		fmt.Fprintln(conn, res.msg)

		// Handle json datas
		if res.datas != "" {
			jsonBytes, err := json.Marshal(res.datas)
			if err != nil {
				fmt.Fprintln(conn, "ERR Internal server error during JSON parsing")
			} else {
				fmt.Fprintln(conn, string(jsonBytes))
			}
		}
	}
}
