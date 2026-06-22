package main

import (
	"errors"
	"fmt"
)

const (
	CmdConnect   = "CONNECT"
	CmdLook      = "LOOK"
	CmdMove      = "MOVE"
	CmdChat      = "CHAT"
	CmdTake      = "TAKE"
	CmdDrop      = "DROP"
	CmdInventory = "INVENTORY"
	CmdTalk      = "TALK"
	CmdAttack    = "ATTACK"
	CmdStatus    = "STATUS"
	CmdQuest     = "QUEST"
	CmdQuests    = "QUESTS"
	CmdWho       = "WHO"
	CmdGroup     = "GROUP"
	CmdQuit      = "QUIT"
)

func handleCmdConnect(clients map[string]*Client, ip string, req []string) (string, error) {
	// Check for invalid command
	if len(req) != 2 {
		return "", errors.New("ERR Invalid name: shouldn't contain space character")
	}

	//Check for duplicated commands
	if clients[ip].connected {
		return "", errors.New("ERR User already connected")
	}

	// Check for name's presence
	already_present := false
	for cli := range clients {
		if clients[cli].name == req[1] {
			already_present = true
			break
		}
	}
	if already_present {
		return "", errors.New("ERR 201 NAME_IN_USE")
	}

	// Save user's name and connection state
	clients[ip].name = req[1]
	clients[ip].connected = true
	return "OK connected", nil
}

func handleCmdWho(clients map[string]*Client, req []string) (string, error) {
	// Check for invalid command
	if len(req) >= 2 {
		return "", errors.New("ERR Invalid command")
	}

	// Get nb of connected clients
	nb_clients := 0
	for cli := range clients {
		if clients[cli].connected {
			nb_clients++
		}
	}

	return fmt.Sprintf("OK players=%d", nb_clients), nil
}
