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

	// Add player to npcs dialogues
	dialogues[req[1]] = make(map[string]int)

	// Inform player ROOM for the new player
	inform_room(clients, clients[ip], clients[ip].datas.room, "EVT ROOM PRESENCE ENTER")

	return "OK connected", "", nil
}

func handleCmdQuit(clients map[string]*Client, cli *Client, req []string) (string, any, error) {
	// Check for invalid command
	if len(req) != 1 {
		return "", "", errors.New("ERR Invalid command")
	}

	// Inform player ROOM for the quit of the player
	inform_room(clients, cli, cli.datas.room, "EVT ROOM PRESENCE LEAVE")

	return "OK bye", "", nil
}

func handleCmdWho(clients map[string]*Client, req []string) (string, any, error) {
	// Check for invalid command
	if len(req) != 1 {
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
	if len(req) != 1 {
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

	// Inform players of user LEAVING
	inform_room(clients, cli, cli.datas.room, "EVT ROOM PRESENCE LEAVE")
	
	// Move player
	if cli, ok := clients[cli.ip]; ok {
		cli.datas.room = nextRoom
	}
	cli.datas.room = nextRoom

	// Inform players of user ENTERING
	inform_room(clients, cli, cli.datas.room, "EVT ROOM PRESENCE ENTER")

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

func handleCmdStatus(cli *Client, req []string) (string, any, error) {
	// Check for invalid command
	if len(req) != 1 {
		return "", "", errors.New("ERR Invalid command")
	}

	// Format status response
	res := make(map[string]any)

	res["status"] = cli.datas.status
	res["max_hp"] = cli.datas.max_hp
	res["hp"] = cli.datas.hp

	return "OK", res, nil
}

func handleCmdTake(cli *Client, req []string) (string, any, error) {
	// Check for invalid command
	if len(req) != 2 {
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
	if len(req) != 2 {
		return "", "", errors.New("ERR Invalid command")
	}

	object := req[1]
	for obj_index, obj := range cli.datas.inventory {
		if obj == object {
			// Remove object to user inventory
			cli.datas.inventory = append(cli.datas.inventory[:obj_index], cli.datas.inventory[obj_index+1:]...)

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
	if len(req) != 1 {
		return "", "", errors.New("ERR Invalid command")
	}

	return "OK", cli.datas.inventory, nil
}

func handleCmdQuest(cli *Client, req []string) (string, any, error) {
	// Check for invalid command
	if len(req) != 2 {
		return "", "", errors.New("ERR Invalid command")
	}

	npc := req[1]

	// Check that npc exist
	for npc_name, npc_datas := range world.Npcs {
		if npc_name == npc {
			// Check npc room
			for _, room_npc := range world.Rooms[cli.datas.room].Npcs {
				if room_npc == npc {

					// Handle empty or validated quest
					if npc_datas.QuestId == "" || world.Quests[npc_datas.QuestId].Status == "unavailable" {
						return "", "", errors.New("ERR 406 NO_QUEST_AVAILABLE")
					}

					quest := world.Quests[npc_datas.QuestId]

					// Format datas
					datas := make(map[string]any)

					datas["status"] = quest.Status
					datas["reward"] = quest.Reward
					datas["description"] = quest.Description
					datas["quest_id"] = npc_datas.QuestId

					return "OK", datas, nil
				}
			}
		}
	}

	// Handle inexistant npc
	return "", "", errors.New("ERR 404 NPC_NOT_FOUND")
}

func handleCmdQuests(req []string) (string, any, error) {
	// Check for invalid command
	if len(req) != 1 {
		return "", "", errors.New("ERR Invalid command")
	}

	// Initialize datas variables
	res := make([]map[string]string, 0)
	var datas map[string]string

	for quest_id, quest := range world.Quests {
		// Format quest datas
		datas = make(map[string]string)

		datas["quest_id"] = quest_id
		datas["status"] = quest.Status

		// Handle active quest
		if quest.Status == "active" {
			datas["progress"] = "1/3"
		}

		// Add quest datas to res
		res = append(res, datas)
	}

	return "OK", res, nil
}

func handleCmdTalk(cli *Client, req []string) (string, any, error) {
	// Check for invalid command
	if len(req) != 2 {
		return "", "", errors.New("ERR Invalid command")
	}

	npc := req[1]

	// Check that npc exist
	for npc_name, npc_datas := range world.Npcs {
		if npc_name == npc {

			// Check npc room
			for _, room_npc := range world.Rooms[cli.datas.room].Npcs {
				if room_npc == npc {

					// Get the npc dialogue index
					_, ok := dialogues[cli.name][npc_name]
					if !ok {
						dialogues[cli.name][npc_name] = 0
					}
					npc_index := dialogues[cli.name][npc_name]

					// Get the npc text
					datas := npc_datas.Dialogue[npc_index%len(npc_datas.Dialogue)]

					// Update npc dialogue index
					dialogues[cli.name][npc_name]++

					return "OK", datas, nil
				}
			}
		}
	}

	// Handle inexistant npc
	return "", "", errors.New("ERR 404 NPC_NOT_FOUND")
}

func handleCmdAttack(cli *Client, req []string) (string, any, error) {
	// Check for invalid command
	if len(req) != 2 {
		return "", "", errors.New("ERR Invalid command")
	}

	npc := req[1]

	// Check that npc exist
	for npc_name, npc_datas := range world.Npcs {
		if npc_name == npc {

			// Check that npc is hostile
			/* if cond {
				return "", "", errors.New("ERR 405 NPC_NOT_HOSTILE")
			} */

			// Check npc room
			for _, room_npc := range world.Rooms[cli.datas.room].Npcs {
				if room_npc == npc {

					// Update values

					// Format datas
					datas := make(map[string]any)

					datas["attacker_hp"] = cli.datas.hp
					datas["target_hp"] = npc_datas.Stats.Hp
					datas["damage"] = 10
					datas["status"] = "combat"

					return "OK", datas, nil
				}
			}
		}
	}

	// Handle inexistant npc
	return "", "", errors.New("ERR 404 NPC_NOT_FOUND")
}
