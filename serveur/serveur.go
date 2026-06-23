package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type Request struct {
	cli *Client
	msg string
}

type Response struct {
	msg   string
	datas any
	req   Request
}

var (
	entering = make(chan *Client)
	leaving  = make(chan *Client)
	requests = make(chan Request)
	groups   map[string][]*Client
	world    Map
)

func main() {
	// Get the world
	var err error
	world, err = get_map("world.json")
	if err != nil {
		fmt.Println("[ERROR]:", err)
		return
	}

	groups = make(map[string][]*Client)

	// Start the serveur
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("[ERROR]: An error occured during the serveur connection:", err)
		return
	}
	defer listener.Close()
	fmt.Println("[INFO]: Connection established on :8080.")

	go broadcaster()

	// Handle new players
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
			clients[cli.ip] = cli

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
	var datas any
	var err error

	activeCli, ok := clients[request.cli.ip]
	if !ok {
		activeCli = request.cli
	}

	// Handle the command type
	switch req[0] {
	case CmdConnect:
		res, datas, err = handleCmdConnect(clients, request.cli.ip, req)
	case CmdQuit:
		res, datas, err = handleCmdQuit(clients, activeCli, req)
	case CmdWho:
		res, datas, err = handleCmdWho(clients, req)
	case CmdLook:
		res, datas, err = handleCmdLook(clients, activeCli, req)
	case CmdMove:
		res, datas, err = handleCmdMove(clients, activeCli, req)
	case CmdChat:
		res, datas, err = handleCmdChat(clients, activeCli, req)
	case CmdGroup:
		res, datas, err = handleCmdGroup(clients, activeCli, req)
	case CmdStatus:
		res, datas, err = handleCmdStatus(activeCli)
	case CmdTake:
		res, datas, err = handleCmdTake(activeCli, req)
	case CmdDrop:
		res, datas, err = handleCmdDrop(activeCli, req)

	default:
		res, datas, err = "", "", errors.New("Invalid command")
	}

	// Handle command errors
	if err != nil {
		res = err.Error()
	}

	// Return the response
	activeCli.ch <- Response{res, datas, request}
}
