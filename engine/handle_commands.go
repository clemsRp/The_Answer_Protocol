package engine

import (
	"errors"
	"log/slog"
	"strings"
	pr "tap/protocol"
)

func (e *Engine) validateSession(id string, cmd string) (string, error) {
	pseudo, exists := e.sessions[id]

	if !exists {
		return "", errors.New(pr.ErrInternalServer)
	}

	if cmd != pr.CmdConnect && pseudo == "" {
		return "", errors.New(pr.ErrNotConnected)
	}

	return pseudo, nil
}

func (e *Engine) handleCommands(request pr.ServerRequest) (string, any, error) {
	req := strings.SplitN(request.Msg, " ", 5)
	cmd := strings.ToUpper(req[0])

	// Define variables
	var pseudo string
	var res string
	var datas any
	var err error

	// validates session if player is connected
	pseudo, err = e.validateSession(request.Id, cmd)
	if err != nil {
		return "", nil, err
	}

	if cmd == pr.CmdConnect {
		return e.handleCmdConnect(request.Id, req)
	}
	player := e.players[pseudo]

	// Handle the command types for authentified players
	switch cmd {
	case pr.CmdQuit:
		res, datas, err = e.handleCmdQuit(player, req)
	case pr.CmdWho:
		res, datas, err = e.handleCmdWho(req)
	case pr.CmdLook:
		res, datas, err = e.handleCmdLook(player, req)
	case pr.CmdUnGrouped:
		res, datas, err = e.handleCmdUnGrouped(player, req)
	case pr.CmdGrouped:
		res, datas, err = e.handleCmdGrouped(player, req)
	case pr.CmdMove:
		res, datas, err = e.handleCmdMove(player, req)
	case pr.CmdChat:
		res, datas, err = e.handleCmdChat(player, req)
	case pr.CmdGroup:
		res, datas, err = e.handleCmdGroup(player, req)
	case pr.CmdStatus:
		res, datas, err = e.handleCmdStatus(player, req)
	case pr.CmdTake:
		res, datas, err = e.handleCmdTake(player, req)
	case pr.CmdDrop:
		res, datas, err = e.handleCmdDrop(player, req)
	case pr.CmdInventory:
		res, datas, err = e.handleCmdInventory(player, req)
	case pr.CmdQuest:
		res, datas, err = e.handleCmdQuest(player, req)
	case pr.CmdQuests:
		res, datas, err = e.handleCmdQuests(player, req)
	case pr.CmdTalk:
		res, datas, err = e.handleCmdTalk(player, req)
	case pr.CmdAttack:
		res, datas, err = e.handleCmdAttack(player, req)
	case pr.CmdFlee:
		res, datas, err = e.handleCmdFlee(player, req)

	default:
		res, datas, err = "", nil, errors.New(pr.ErrInvalidCommand)
	}

	// Log server response
	if err != nil {
		res = err.Error()
		slog.Error("Server response", "response", res, "command", request.Msg)

	} else if datas != "" {
		slog.Info("Server response", "response", res, "datas", datas, "command", request.Msg)

	} else {
		slog.Info("Server response", "response", res, "command", request.Msg)
	}

	return res, datas, err
}
