package engine

import (
	"errors"
	"strings"
	pr "tap/protocol"
)

func (e *Engine) handleCommands(request pr.ServerRequest) (*pr.Client, string, any, error) {
	var res string
	var datas any
	var err error

	// Split request
	req := strings.SplitN(request.Msg, " ", 5)

	// fmt.Println(req)

	// Get "true" client
	activeCli, ok := e.clients[request.Cli.Ip]
	if !ok {
		activeCli = request.Cli
	}
	name := activeCli.Name
	if name == "" {
		name = activeCli.Ip
	}

	if strings.ToUpper(req[0]) != pr.CmdConnect && !request.Cli.Datas.Connected {
		res, datas, err = "", nil, errors.New(pr.ErrNotConnected)

	} else {
		// Handle the command type
		switch strings.ToUpper(req[0]) {
		case pr.CmdConnect:
			res, datas, err = e.handleCmdConnect(request.Cli.Ip, req)
		case pr.CmdQuit:
			res, datas, err = e.handleCmdQuit(activeCli, req)
		case pr.CmdWho:
			res, datas, err = e.handleCmdWho(req)
		case pr.CmdLook:
			res, datas, err = e.handleCmdLook(activeCli, req)
		case pr.CmdMove:
			res, datas, err = e.handleCmdMove(activeCli, req)
		case pr.CmdChat:
			res, datas, err = e.handleCmdChat(activeCli, req)
		case pr.CmdGroup:
			res, datas, err = e.handleCmdGroup(activeCli, req)
		case pr.CmdStatus:
			res, datas, err = e.handleCmdStatus(activeCli, req)
		case pr.CmdTake:
			res, datas, err = e.handleCmdTake(activeCli, req)
		case pr.CmdDrop:
			res, datas, err = e.handleCmdDrop(activeCli, req)
		case pr.CmdInventory:
			res, datas, err = e.handleCmdInventory(activeCli, req)
		case pr.CmdQuest:
			res, datas, err = e.handleCmdQuest(activeCli, req)
		case pr.CmdQuests:
			res, datas, err = e.handleCmdQuests(activeCli, req)
		case pr.CmdTalk:
			res, datas, err = e.handleCmdTalk(activeCli, req)
		case pr.CmdAttack:
			res, datas, err = e.handleCmdAttack(activeCli, req)

		default:
			res, datas, err = "", nil, errors.New("ERR 400 BAD_REQUEST")
		}
	}

	return activeCli, res, datas, err
}
