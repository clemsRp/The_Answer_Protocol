package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type Datas struct {
	room       string
	inventory  []string
	group      string
	invitation []string
	hp         int
	max_hp     int
	status     string
	connected  bool
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
		datas: Datas{"start", []string{}, "", []string{}, 50, 50, "healthy", false},
	}

	cli.ch <- Response{"[INFO]: You are connected as " + who, "", Request{}}
	entering <- cli

	input := bufio.NewScanner(conn)
	for input.Scan() {
		requests <- Request{cli, input.Text()}
	}

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
