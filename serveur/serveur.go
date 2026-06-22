package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type Request struct {
	cli Client
	msg string
}

type Response struct {
	msg string
	req Request
}

map_path := "world.json"

var (
	entering = make(chan Client)
	leaving  = make(chan Client)
	requests = make(chan Request)
	map		 = get_map(map_path)
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
	clients := make(map[string]*Client)

	for {
		select {
		case req := <-requests:
			handleRequest(clients, req)

		case cli := <-entering:
			clients[cli.ip] = &cli

		case cli := <-leaving:
			if c, ok := clients[cli.ip]; ok {
                delete(clients, cli.ip)
                close(c.ch)
            }
		}
	}
}

func handleRequest(clients map[string]*Client, request Request) {
	req := strings.Split(request.msg, " ")

	var res string
	var err error

	switch req[0] {
	case CmdConnect:
		res, err = handleCmdConnect(clients, request.cli.ip, req)
	case CmdWho:
		res, err = handleCmdWho(clients, req)

	default:
		res, err = "", errors.New("Invalid command")
	}

	if err != nil {
		res = err.Error()
	}
	request.cli.ch <- Response{res, request}
}
