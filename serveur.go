package main

import (
	"fmt"
	"net"
	"strings"
)

var (
	entering = make(chan Client)
	leaving  = make(chan Client)
	requests = make(chan Request)
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("[ERROR]: An error occured during the serveur connection:", err)
		return
	}
	defer listener.Close()

	fmt.Println("[INFO]: Connection established on :8080.")

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[ERROR]: An error occured during a client's connection:", err)
			continue
		}
		go handleClient(conn)
	}
}

func broadcaster() {
	clients := make(map[Client]bool)

	for {
		select {
		case req := <-requests:
			handleRequest(req)

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

func handleRequest(request Request) {
	req := strings.Split(request.msg, " ")
	var res string

	switch req[0] {
	case CmdConnect:
		res = handleCmdConnect(req)

	default:
		res = "Error"
	}

	request.cli.ch <- Response{res}
}
