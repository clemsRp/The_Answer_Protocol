package main

import (
	parser "tap/server/parser"

	"errors"
	"net"
	"strings"
	"time"
)

type Log struct {
	Timestamp string         `json:"timestamp"`
	Level     string         `json:"level"`
	Message   string         `json:"message"`
	Datas     map[string]any `json:"datas,omitempty"`
}

type Request struct {
	cli *Client
	msg string
}

type Response struct {
	Msg   string
	Datas any
	Req   Request
}

var (
	requests = make(chan Request)
	logs     = make(chan Log, 500)
	entering = make(chan *Client)
	leaving  = make(chan *Client)

	groups    map[string][]*Client
	dialogues map[string]map[string]int

	t_start = time.Now().Unix()

	world parser.Map
)

func main() {
	// Get the world
	var err error
	world, err = parser.Get_map("world.json")
	if err != nil {
		writeLog("ERROR", err.Error(), nil)
		return
	}

	groups = make(map[string][]*Client)
	dialogues = make(map[string]map[string]int)

	// Start the serveur
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		LogError("An error occured during the serveur connection:", nil)
		return
	}
	defer listener.Close()
	LogInfo("Connection established on :8080.", nil)

	go broadcaster()

	// Handle new players
	for {
		conn, err := listener.Accept()
		if err != nil {
			LogError("An error occurred during a client's connection", map[string]any{
				"error": err.Error(),
			})
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
			writeLog(log.Level, log.Message, log.Datas)

		case req := <-requests:
			handleRequest(clients, req)

		case cli := <-entering:
			clients[cli.ip] = cli
			LogInfo("Start of connection", map[string]any{
				"ip":       cli.ip,
				"duration": get_timestamp(),
			})

		case cli := <-leaving:
			if c, ok := clients[cli.ip]; ok {
				delete(clients, cli.ip)
				close(c.ch)
			}
			LogInfo("End of connection", map[string]any{
				"ip":       cli.ip,
				"duration": get_timestamp(),
				"player":   cli.name,
			})
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
	name := activeCli.name
	if name == "" {
		name = activeCli.ip
	}

	logCtx := make(map[string]any)
	logCtx["ip"] = request.cli.ip
	logCtx["player"] = name
	logCtx["cmd"] = req[0]
	logCtx["full_msg"] = request.msg

	if len(request.msg) > 1024 {
		res, datas, err = "", "", errors.New("ERR 413 REQUEST_ENTITY_TOO_LARGE")
	}

	// Force first command to be CONNECT
	if err == nil && strings.ToUpper(req[0]) != CmdConnect && !request.cli.datas.connected {
		res, datas, err = "", "", errors.New("ERR 401 UNAUTHORIZED")

	} else {
		// Handle the command type
		switch strings.ToUpper(req[0]) {
		case CmdConnect:
			res, datas, err = handleCmdConnect(clients, request.cli.ip, req, logCtx)
		case CmdQuit:
			res, datas, err = handleCmdQuit(clients, activeCli, req, logCtx)
		case CmdWho:
			res, datas, err = handleCmdWho(clients, req, logCtx)
		case CmdLook:
			res, datas, err = handleCmdLook(clients, activeCli, req, logCtx)
		case CmdMove:
			res, datas, err = handleCmdMove(clients, activeCli, req, logCtx)
		case CmdChat:
			res, datas, err = handleCmdChat(clients, activeCli, req, logCtx)
		case CmdGroup:
			res, datas, err = handleCmdGroup(clients, activeCli, req, logCtx)
		case CmdStatus:
			res, datas, err = handleCmdStatus(activeCli, req, logCtx)
		case CmdTake:
			res, datas, err = handleCmdTake(activeCli, req, logCtx)
		case CmdDrop:
			res, datas, err = handleCmdDrop(activeCli, req, logCtx)
		case CmdInventory:
			res, datas, err = handleCmdInventory(activeCli, req, logCtx)
		case CmdQuest:
			res, datas, err = handleCmdQuest(activeCli, req, logCtx)
		case CmdQuests:
			res, datas, err = handleCmdQuests(req, logCtx)
		case CmdTalk:
			res, datas, err = handleCmdTalk(activeCli, req, logCtx)
		case CmdAttack:
			res, datas, err = handleCmdAttack(activeCli, req, logCtx)

		default:
			res, datas, err = "", "", errors.New("Invalid command")
		}
	}
	statusCode := "200"
	if err != nil {
		res, datas = err.Error(), ""

		errParts := strings.Split(err.Error(), " ")
		if len(errParts) >= 2 && errParts[0] == "ERR" {
			statusCode = errParts[1]
		} else {
			statusCode = "500"
		}
	}

	logCtx["status_code"] = statusCode
	logCtx["response"] = res

	if err != nil {
		logCtx["error"] = err.Error()
		LogError("Command failed", logCtx)
	} else {
		LogInfo("Command success", logCtx)
	}

	activeCli.ch <- Response{res, datas, request}
}
