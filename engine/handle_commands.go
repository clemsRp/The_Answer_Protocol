package engine

import (
	"errors"
	"strings"
	pr "tap/protocol"
)

func (e *Engine) validateSession(ip string, cmd string) (string, error) {
	pseudo, exists := e.sessions[ip]

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

	// validates session if player is connected
	pseudo, err := e.validateSession(request.Ip, cmd)
	if err != nil {
		return "", nil, err
	}

	if cmd == pr.CmdConnect {
		return e.handleCmdConnect(request.Ip, req)
	}
	player := e.players[pseudo]

	// Handle the command types for authentified players
	switch cmd {
	case pr.CmdQuit:
		return e.handleCmdQuit(player, req)
	case pr.CmdWho:
		return e.handleCmdWho(req)
	case pr.CmdLook:
		return e.handleCmdLook(player, req)
	case pr.CmdMove:
		return e.handleCmdMove(player, req)
	case pr.CmdChat:
		return e.handleCmdChat(player, req)
	case pr.CmdGroup:
		return e.handleCmdGroup(player, req)
	case pr.CmdStatus:
		return e.handleCmdStatus(player, req)
	case pr.CmdTake:
		return e.handleCmdTake(player, req)
	case pr.CmdDrop:
		return e.handleCmdDrop(player, req)
	case pr.CmdInventory:
		return e.handleCmdInventory(player, req)
	case pr.CmdQuest:
		return e.handleCmdQuest(player, req)
	case pr.CmdQuests:
		return e.handleCmdQuests(player, req)
	case pr.CmdTalk:
		return e.handleCmdTalk(player, req)
	case pr.CmdAttack:
		return e.handleCmdAttack(player, req)

	default:
		return "", nil, errors.New(pr.ErrInvalidCommand)
	}

}
