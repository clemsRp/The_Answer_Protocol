package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type Log struct {
	msg string
}

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
	requests = make(chan Request)
	logs     = make(chan Log, 500)
	entering = make(chan *Client)
	leaving  = make(chan *Client)

	groups    map[string][]*Client
	dialogues map[string]map[string]int

	t_start = time.Now().Unix()

	world Map
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
	dialogues = make(map[string]map[string]int)

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
		case log := <-logs:
			fmt.Println(log.msg)

		case req := <-requests:
			handleRequest(clients, req)

		case cli := <-entering:
			clients[cli.ip] = cli

			fmt.Printf("[INFO]: Start of connection at %s /", cli.ip)

			timestamp := get_timestamp()
			print_timestamp(timestamp)

		case cli := <-leaving:
			if c, ok := clients[cli.ip]; ok {
				delete(clients, cli.ip)
				close(c.ch)
			}

			fmt.Printf("[INFO]: End of connection at %s /", cli.ip)

			timestamp := get_timestamp()
			print_timestamp(timestamp)
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

	// Force first command to be CONNECT
	if req[0] != CmdConnect && !request.cli.datas.connected {
		res, datas, err = "", "", errors.New("ERR CONNECT user first before doing any commands")

	} else {
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
			res, datas, err = handleCmdStatus(activeCli, req)
		case CmdTake:
			res, datas, err = handleCmdTake(activeCli, req)
		case CmdDrop:
			res, datas, err = handleCmdDrop(activeCli, req)
		case CmdInventory:
			res, datas, err = handleCmdInventory(activeCli, req)
		case CmdQuest:
			res, datas, err = handleCmdQuest(activeCli, req)
		case CmdQuests:
			res, datas, err = handleCmdQuests(req)
		case CmdTalk:
			res, datas, err = handleCmdTalk(activeCli, req)
		case CmdAttack:
			res, datas, err = handleCmdAttack(activeCli, req)

		default:
			res, datas, err = "", "", errors.New("Invalid command")
		}
	}

	// Handle command errors and log type
	msg_type := "INFO"
	if err != nil {
		res, datas = err.Error(), ""
		msg_type = "ERROR"
	}

	name := request.cli.name
	if name == "" {
		name = request.cli.ip
	}

	// Log command
	log := fmt.Sprintf("[%s]: '%s' (%s)", msg_type, request.msg, name)

	if err != nil {
		log += " -> " + err.Error()
	}

	logs <- Log{log}
	logs <- Log{fmt.Sprintf("[%s]: %s", msg_type, res)}

	// Return the response
	activeCli.ch <- Response{res, datas, request}
}
