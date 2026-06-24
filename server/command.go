package main

import (
	"errors"
	"fmt"
	"strings"
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

	GlobalChat = "GLOBAL"
	RoomChat   = "ROOM"
	GroupChat  = "GROUP"

	CreateGroup = "CREATE"
	InviteGroup = "INVITE"
	JoinGroup   = "JOIN"
	LeaveGroup  = "LEAVE"

	South = "south"
	North = "north"
	East  = "east"
	West  = "west"
)

func handleCmdConnect(clients map[string]*Client, ip string, req []string) (string, any, error) {
	// Check for invalid command
	if len(req) != 2 {
		return "", "", errors.New("ERR Invalid name: shouldn't contain space character")
	}

	//Check for duplicated commands
	if clients[ip].datas.connected {
		return "", "", errors.New("ERR User already connected")
	}

	// Check for name's presence
	for _, cli := range clients {
		if cli.name == req[1] {
			return "", "", errors.New("ERR 201 NAME_IN_USE")
		}
	}

	// Save user's name and connection state
	clients[ip].name = req[1]
	clients[ip].datas.connected = true
	return "OK connected", "", nil
}

func handleCmdQuit(clients map[string]*Client, cli *Client, req []string) (string, any, error) {
	// Check for invalid command
	if len(req) >= 2 {
		return "", "", errors.New("ERR Invalid command")
	}

	return "OK bye", "", nil
}

func handleCmdWho(clients map[string]*Client, req []string) (string, any, error) {
	// Check for invalid command
	if len(req) >= 2 {
		return "", "", errors.New("ERR Invalid command")
	}

	// Get nb of connected clients
	nb_clients := 0
	for cli := range clients {
		if clients[cli].datas.connected {
			nb_clients++
		}
	}

	return fmt.Sprintf("OK players=%d", nb_clients), "", nil
}

func handleCmdLook(clients map[string]*Client, cli *Client, req []string) (string, any, error) {
	// Check for invalid command
	if len(req) >= 2 {
		return "", "", errors.New("ERR Invalid command")
	}

	// Initialize the datas map
	res := make(map[string]any)
	room := make(map[string]any)
	players := make([]string, 0)

	// Get the room players
	for ip := range clients {
		if clients[ip].datas.room == cli.datas.room && clients[ip].datas.connected {
			players = append(players, clients[ip].name)
		}
	}

	currentRoom := world.Rooms[cli.datas.room]

	// Format the datas
	res["npcs"] = currentRoom.Npcs
	res["items"] = currentRoom.Items
	res["players"] = players
	res["room"] = room

	room["id"] = "room." + cli.datas.room
	room["exits"] = currentRoom.Exits
	room["description"] = currentRoom.Description
	room["name"] = currentRoom.Name

	return "OK", res, nil
}

func handleCmdMove(clients map[string]*Client, cli *Client, req []string) (string, any, error) {
	// Check for invalid command
	if len(req) != 2 {
		return "", "", errors.New("ERR Invalid command")
	}

	direction := req[1]
	currentRoom := world.Rooms[cli.datas.room]

	// Handle No exit errors
	nextRoom, exists := currentRoom.Exits[direction]
	if !exists {
		return "", "", errors.New("ERR 301 NO_EXIT")
	}

	// Move player
	if cli, ok := clients[cli.ip]; ok {
		cli.datas.room = nextRoom
	}
	cli.datas.room = nextRoom

	return fmt.Sprintf("OK room=%s", nextRoom), "", nil
}

func handleCmdChat(clients map[string]*Client, cli *Client, req []string) (string, any, error) {
	// Check for invalid command
	if len(req) < 3 {
		return "", "", errors.New("ERR Invalid command")
	}

	var chat string

	scope := req[1]
	msg := strings.Join(req[2:], " ")

	for ip := range clients {
		if clients[ip].name != cli.name {
			// Handle the scopes
			is_global := scope == GlobalChat
			is_group := scope == GroupChat && cli.datas.group != "" && cli.datas.group == clients[ip].datas.group
			is_room := scope == RoomChat && cli.datas.room == clients[ip].datas.room

			// Send chat message
			if is_global || is_group || is_room {
				chat = "[CHAT] " + cli.name + ": " + msg
				clients[ip].ch <- Response{chat, "", Request{}}
			}
		}
	}

	return "OK", "", nil
}

func handleCmdGroup(clients map[string]*Client, cli *Client, req []string) (string, any, error) {
	var err error
	var res string

	scope := req[1]

	// Check for invalid command
	if (len(req) != 3 && scope != LeaveGroup) || (len(req) > 2 && scope == LeaveGroup) {
		return "", "", errors.New("ERR Invalid command")
	}

	var arg string
	if scope != LeaveGroup {
		arg = req[2]
	}

	switch scope {
	case CreateGroup:
		res, err = create_group(cli, arg)
	case InviteGroup:
		res, err = invite_user_in_group(clients, cli, arg)
	case JoinGroup:
		res, err = join_group(cli, arg)
	case LeaveGroup:
		res, err = leave_group(cli)

	default:
		return "", "", errors.New("ERR Invalid scope")
	}

	if err != nil {
		return "", "", err
	}

	fmt.Println(groups)

	return res, "", nil
}

func handleCmdStatus(cli *Client) (string, any, error) {
	// Format status response
	res := make(map[string]any)

	res["status"] = cli.datas.status
	res["max_hp"] = cli.datas.max_hp
	res["hp"] = cli.datas.hp

	return "OK", res, nil
}

func handleCmdTake(cli *Client, req []string) (string, any, error) {
	// Check for invalid command
	if len(req) >= 3 {
		return "", "", errors.New("ERR Invalid command")
	}

	object := req[1]
	for obj_index, obj := range world.Rooms[cli.datas.room].Items {
		if obj == object {
			// Add object to user inventory
			cli.datas.inventory = append(cli.datas.inventory, object)

			// Remove object to map
			world.Rooms[cli.datas.room].Items = append(world.Rooms[cli.datas.room].Items[:obj_index], world.Rooms[cli.datas.room].Items[obj_index+1:]...)

			return "OK taken=" + object, "", nil
		}
	}

	// Handle Invalid object
	return "", "", errors.New("ERR 404 ITEM_NOT_FOUND")
}

func handleCmdDrop(cli *Client, req []string) (string, any, error) {
	// Check for invalid command
	if len(req) >= 3 {
		return "", "", errors.New("ERR Invalid command")
	}

	object := req[1]
	for obj_index, obj := range cli.datas.inventory {
		if obj == object {
			// Remove object to user inventory
			cli.datas.inventory = append(cli.datas.inventory[:obj_index], cli.datas.inventory[obj_index + 1:]...)

			// Add object to map
			world.Rooms[cli.datas.room].Items = append(world.Rooms[cli.datas.room].Items, object)

			return "OK dropped=" + object, "", nil
		}
	}

	// Handle Invalid object
	return "", "", errors.New("ERR 404 ITEM_NOT_IN_INVENTORY")
}

func handleCmdInventory(cli *Client, req []string) (string, any, error) {
	// Check for invalid command
	if len(req) >= 2 {
		return "", "", errors.New("ERR Invalid command")
	}

	return "OK", cli.datas.inventory, nil
}
