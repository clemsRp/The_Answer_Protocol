package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleClient(conn net.Conn) {
	responses := make(chan Response)
	go clientWriter(conn, responses)

	who := conn.RemoteAddr().String()
	cli := Client{conn, responses, who}

	cli.ch <- Response{"[INFO]: You are connected as " + who}
	entering <- cli

	input := bufio.NewScanner(conn)
	for input.Scan() {
		requests <- Request{cli, input.Text()}
	}

	leaving <- cli
	conn.Close()
}

func clientWriter(conn net.Conn, responses <-chan Response) {
	for res := range responses {
		fmt.Fprintln(conn, res.msg)
	}
}
