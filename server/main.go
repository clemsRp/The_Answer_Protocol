package main

import (
	pr "tap/protocol"
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

var (
	requests = make(chan pr.Request)
	logs     = make(chan Log, 500)
	entering = make(chan *pr.Client)
	leaving  = make(chan *pr.Client)

	groups    map[string][]*pr.Client
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

	groups = make(map[string][]*pr.Client)
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
	clients := make(map[string]*pr.Client)

	for {
		select {
		// Print server logs
		case log := <-logs:
			writeLog(log.Level, log.Message, log.Datas)

		// Handle a client command
		case req := <-requests:
			handleRequest(clients, req)

		// Handle the start of the a client connection
		case cli := <-entering:
			clients[cli.Ip] = cli
			LogInfo("Start of connection", map[string]any{
				"ip":       cli.Ip,
				"duration": get_timestamp(),
			})

		// Handle the end of the a client connection
		case cli := <-leaving:
			if c, ok := clients[cli.Ip]; ok {
				delete(clients, cli.Ip)
				close(c.Ch)
			}
			LogInfo("End of connection", map[string]any{
				"ip":       cli.Ip,
				"duration": get_timestamp(),
				"player":   cli.Name,
			})
		}
	}
}

func handleRequest(clients map[string]*pr.Client, request pr.Request) {
	// Get commands
	req := strings.Split(request.Msg, " ")

	var res string
	var datas any
	var err error

	// Get "true" client
	activeCli, ok := clients[request.Cli.Ip]
	if !ok {
		activeCli = request.Cli
	}
	name := activeCli.Name
	if name == "" {
		name = activeCli.Ip
	}

	logCtx := make(map[string]any)
	logCtx["ip"] = request.Cli.Ip
	logCtx["player"] = name
	logCtx["cmd"] = req[0]
	logCtx["full_msg"] = request.Msg

	if len(request.Msg) > 1024 {
		res, datas, err = "", "", errors.New("ERR 413 REQUEST_ENTITY_TOO_LARGE")
	}

	// Force first command to be CONNECT
	if err == nil && strings.ToUpper(req[0]) != pr.CmdConnect && !request.Cli.Datas.Connected {
		res, datas, err = "", "", errors.New("ERR 401 UNAUTHORIZED")

	} else {
		// Handle the command type
		switch strings.ToUpper(req[0]) {
		case pr.CmdConnect:
			res, datas, err = handleCmdConnect(clients, request.Cli.Ip, req, logCtx)
		case pr.CmdQuit:
			res, datas, err = handleCmdQuit(clients, activeCli, req, logCtx)
		case pr.CmdWho:
			res, datas, err = handleCmdWho(clients, req, logCtx)
		case pr.CmdLook:
			res, datas, err = handleCmdLook(clients, activeCli, req, logCtx)
		case pr.CmdMove:
			res, datas, err = handleCmdMove(clients, activeCli, req, logCtx)
		case pr.CmdChat:
			res, datas, err = handleCmdChat(clients, activeCli, req, logCtx)
		case pr.CmdGroup:
			res, datas, err = handleCmdGroup(clients, activeCli, req, logCtx)
		case pr.CmdStatus:
			res, datas, err = handleCmdStatus(activeCli, req, logCtx)
		case pr.CmdTake:
			res, datas, err = handleCmdTake(activeCli, req, logCtx)
		case pr.CmdDrop:
			res, datas, err = handleCmdDrop(activeCli, req, logCtx)
		case pr.CmdInventory:
			res, datas, err = handleCmdInventory(activeCli, req, logCtx)
		case pr.CmdQuest:
			res, datas, err = handleCmdQuest(activeCli, req, logCtx)
		case pr.CmdQuests:
			res, datas, err = handleCmdQuests(req, logCtx)
		case pr.CmdTalk:
			res, datas, err = handleCmdTalk(activeCli, req, logCtx)
		case pr.CmdAttack:
			res, datas, err = handleCmdAttack(activeCli, req, logCtx)

		default:
			res, datas, err = "", "", errors.New("ERR 400 BAD_REQUEST")
		}
	}

	// Handle command errors and status codes
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

	// Send server response to the client
	activeCli.Ch <- pr.Response{Msg: res, Datas: datas, Req: request}
}
