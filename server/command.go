package main

import (
	"errors"
	"fmt"
	"strings"
	pr "tap/protocol"
)

const (
	South = "south"
	North = "north"
	East  = "east"
	West  = "west"
)

func handleCmdConnect(clients map[string]*pr.Client, ip string, req []string, logCtx map[string]any) (string, any, error) {
	if len(req) != 2 {
		return "", "", errors.New("ERR 400 Invalid name: shouldn't contain space character")
	}
	if clients[ip].Datas.Connected {
		return "", "", errors.New("ERR 403 User already connected")
	}
	for _, cli := range clients {
		if cli.Name == req[1] {
			return "", "", errors.New("ERR 201 NAME_IN_USE")
		}
	}

	clients[ip].Name = req[1]
	clients[ip].Datas.Connected = true
	dialogues[req[1]] = make(map[string]int)
	inform_room(clients, clients[ip], clients[ip].Datas.Room, "EVT ROOM PRESENCE ENTER")

	logCtx["action"] = "player_connected"

	return "OK connected", "", nil
}

func handleCmdQuit(clients map[string]*pr.Client, cli *pr.Client, req []string, logCtx map[string]any) (string, any, error) {
	if len(req) != 1 {
		return "", "", errors.New("ERR 400 BAD_REQUEST Invalid command")
	}
	inform_room(clients, cli, cli.Datas.Room, "EVT ROOM PRESENCE LEAVE")
	return "OK bye", "", nil
}

func handleCmdWho(clients map[string]*pr.Client, req []string, logCtx map[string]any) (string, any, error) {
	if len(req) != 1 {
		return "", "", errors.New("ERR 400 BAD_REQUEST Invalid command")
	}
	nb_clients := 0
	for cli := range clients {
		if clients[cli].Datas.Connected {
			nb_clients++
		}
	}
	return fmt.Sprintf("OK players=%d", nb_clients), "", nil
}

func handleCmdLook(clients map[string]*pr.Client, cli *pr.Client, req []string, logCtx map[string]any) (string, any, error) {
	if len(req) != 1 {
		return "", "", errors.New("ERR 400 BAD_REQUEST Invalid command")
	}

	res := make(map[string]any)
	Room := make(map[string]any)
	players := make([]string, 0)

	for ip := range clients {
		if clients[ip].Datas.Room == cli.Datas.Room && clients[ip].Datas.Connected {
			players = append(players, clients[ip].Name)
		}
	}

	currentRoom := world.Rooms[cli.Datas.Room]
	res["npcs"] = currentRoom.Npcs
	res["items"] = currentRoom.Items
	res["players"] = players
	res["room"] = Room
	Room["id"] = "room." + cli.Datas.Room
	Room["exits"] = currentRoom.Exits
	Room["description"] = currentRoom.Description
	Room["Name"] = currentRoom.Name

	return "OK", res, nil
}

func handleCmdMove(clients map[string]*pr.Client, cli *pr.Client, req []string, logCtx map[string]any) (string, any, error) {
	if len(req) != 2 {
		return "", "", errors.New("ERR 400 BAD_REQUEST Invalid command")
	}

	direction := req[1]
	currentRoom := world.Rooms[cli.Datas.Room]

	nextRoom, exists := currentRoom.Exits[direction]
	if !exists {
		return "", "", errors.New("ERR 301 NO_EXIT")
	}

	inform_room(clients, cli, cli.Datas.Room, "EVT ROOM PRESENCE LEAVE")

	if cli, ok := clients[cli.Ip]; ok {
		cli.Datas.Room = nextRoom
	}
	cli.Datas.Room = nextRoom

	inform_room(clients, cli, cli.Datas.Room, "EVT ROOM PRESENCE ENTER")

	// Enrichissement du contexte de log
	logCtx["action"] = "player_moved"
	logCtx["prev_room"] = currentRoom.Name
	logCtx["new_room"] = nextRoom

	return fmt.Sprintf("OK room=%s", nextRoom), "", nil
}

func handleCmdChat(clients map[string]*pr.Client, cli *pr.Client, req []string, logCtx map[string]any) (string, any, error) {
	if len(req) < 3 {
		return "", "", errors.New("ERR 400 BAD_REQUEST Invalid command")
	}

	var chat string
	scope := strings.ToUpper(req[1])
	msg := strings.Join(req[2:], " ")

	for ip := range clients {
		if clients[ip].Name != cli.Name {
			is_global := scope == pr.GlobalChat
			is_group := scope == pr.GroupChat && cli.Datas.Group != "" && cli.Datas.Group == clients[ip].Datas.Group
			is_room := scope == pr.RoomChat && cli.Datas.Room == clients[ip].Datas.Room

			if is_global || is_group || is_room {
				chat = "[CHAT] " + cli.Name + ": " + msg
				clients[ip].Ch <- pr.Response{Msg: chat, Req: pr.Request{}}
			}
		}
	}
	return "OK", "", nil
}

func handleCmdGroup(clients map[string]*pr.Client, cli *pr.Client, req []string, logCtx map[string]any) (string, any, error) {
	var err error
	var res string
	scope := strings.ToUpper(req[1])

	if (len(req) != 3 && scope != pr.LeaveGroup) || (len(req) > 2 && scope == pr.LeaveGroup) {
		return "", "", errors.New("ERR 400 BAD_REQUEST Invalid command")
	}

	var arg string
	if scope != pr.LeaveGroup {
		arg = req[2]
	}

	switch scope {
	case pr.CreateGroup:
		res, err = create_group(cli, arg)
	case pr.InviteGroup:
		res, err = invite_user_in_group(clients, cli, arg)
	case pr.JoinGroup:
		res, err = join_group(cli, arg)
	case pr.LeaveGroup:
		res, err = leave_group(cli)
	default:
		return "", "", errors.New("ERR 400 Invalid scope")
	}

	if err != nil {
		return "", "", err
	}

	logCtx["action"] = "group_update"
	logCtx["group_scope"] = scope

	return res, "", nil
}

func handleCmdStatus(cli *pr.Client, req []string, logCtx map[string]any) (string, any, error) {
	if len(req) != 1 {
		return "", "", errors.New("ERR 400 BAD_REQUEST Invalid command")
	}
	res := make(map[string]any)
	res["status"] = cli.Datas.Status
	res["max_hp"] = cli.Datas.Max_hp
	res["hp"] = cli.Datas.Hp
	return "OK", res, nil
}

func handleCmdTake(cli *pr.Client, req []string, logCtx map[string]any) (string, any, error) {
	if len(req) != 2 {
		return "", "", errors.New("ERR 400 BAD_REQUEST Invalid command")
	}

	object := req[1]
	for obj_index, obj := range world.Rooms[cli.Datas.Room].Items {
		if obj == object {
			cli.Datas.Inventory = append(cli.Datas.Inventory, object)
			world.Rooms[cli.Datas.Room].Items = append(world.Rooms[cli.Datas.Room].Items[:obj_index], world.Rooms[cli.Datas.Room].Items[obj_index+1:]...)

			// Enrichissement du contexte de log
			logCtx["action"] = "item_taken"
			logCtx["item"] = object

			return "OK taken=" + object, "", nil
		}
	}

	return "", "", errors.New("ERR 404 ITEM_NOT_FOUND")
}

func handleCmdDrop(cli *pr.Client, req []string, logCtx map[string]any) (string, any, error) {
	if len(req) != 2 {
		return "", "", errors.New("ERR 400 BAD_REQUEST Invalid command")
	}

	object := req[1]
	for obj_index, obj := range cli.Datas.Inventory {
		if obj == object {
			cli.Datas.Inventory = append(cli.Datas.Inventory[:obj_index], cli.Datas.Inventory[obj_index+1:]...)
			world.Rooms[cli.Datas.Room].Items = append(world.Rooms[cli.Datas.Room].Items, object)

			// Enrichissement du contexte de log
			logCtx["action"] = "item_dropped"
			logCtx["item"] = object

			return "OK dropped=" + object, "", nil
		}
	}

	return "", "", errors.New("ERR 404 ITEM_NOT_IN_INVENTORY")
}

func handleCmdInventory(cli *pr.Client, req []string, logCtx map[string]any) (string, any, error) {
	if len(req) != 1 {
		return "", "", errors.New("ERR 400 BAD_REQUEST Invalid command")
	}
	return "OK", cli.Datas.Inventory, nil
}

func handleCmdQuest(cli *pr.Client, req []string, logCtx map[string]any) (string, any, error) {
	if len(req) != 2 {
		return "", "", errors.New("ERR 400 BAD_REQUEST Invalid command")
	}

	npc := req[1]
	for npc_name, npc_datas := range world.Npcs {
		if npc_name == npc {
			for _, room_npc := range world.Rooms[cli.Datas.Room].Npcs {
				if room_npc == npc {
					if npc_datas.QuestId == "" || world.Quests[npc_datas.QuestId].Status == "unavailable" {
						return "", "", errors.New("ERR 406 NO_QUEST_AVAILABLE")
					}
					quest := world.Quests[npc_datas.QuestId]
					Datas := make(map[string]any)
					Datas["status"] = quest.Status
					Datas["reward"] = quest.Reward
					Datas["description"] = quest.Description
					Datas["quest_id"] = npc_datas.QuestId
					return "OK", Datas, nil
				}
			}
		}
	}
	return "", "", errors.New("ERR 404 NPC_NOT_FOUND")
}

func handleCmdQuests(req []string, logCtx map[string]any) (string, any, error) {
	if len(req) != 1 {
		return "", "", errors.New("ERR 400 BAD_REQUEST Invalid command")
	}

	res := make([]map[string]string, 0)
	for quest_id, quest := range world.Quests {
		Datas := make(map[string]string)
		Datas["quest_id"] = quest_id
		Datas["status"] = quest.Status
		if quest.Status == "active" {
			Datas["progress"] = "1/3"
		}
		res = append(res, Datas)
	}
	return "OK", res, nil
}
func handleCmdTalk(cli *pr.Client, req []string, logCtx map[string]any) (string, any, error) {
	if len(req) != 2 {
		return "", "", errors.New("ERR 400 BAD_REQUEST Invalid command")
	}

	npc := req[1]
	for npc_name, npc_datas := range world.Npcs {
		if npc_name == npc {
			for _, room_npc := range world.Rooms[cli.Datas.Room].Npcs {
				if room_npc == npc {
					_, ok := dialogues[cli.Name][npc_name]
					if !ok {
						dialogues[cli.Name][npc_name] = 0
					}
					npc_index := dialogues[cli.Name][npc_name]
					Datas := npc_datas.Dialogue[npc_index%len(npc_datas.Dialogue)]
					dialogues[cli.Name][npc_name]++

					// Enrichissement du contexte de log
					logCtx["action"] = "dialogue_advanced"
					logCtx["npc"] = npc
					logCtx["dialogue_index"] = dialogues[cli.Name][npc_name]

					return "OK", Datas, nil
				}
			}
		}
	}

	return "", "", errors.New("ERR 404 NPC_NOT_FOUND")
}

func handleCmdAttack(cli *pr.Client, req []string, logCtx map[string]any) (string, any, error) {
	if len(req) != 2 {
		return "", "", errors.New("ERR 400 BAD_REQUEST Invalid command")
	}

	npc := req[1]
	for npc_name, npc_datas := range world.Npcs {
		if npc_name == npc {
			for _, room_npc := range world.Rooms[cli.Datas.Room].Npcs {
				if room_npc == npc {
					Datas := make(map[string]any)
					Datas["attacker_hp"] = cli.Datas.Hp
					Datas["target_hp"] = npc_datas.Stats.Hp
					Datas["damage"] = 10
					Datas["status"] = "combat"

					// Enrichissement du contexte de log
					logCtx["action"] = "attack"
					logCtx["target_npc"] = npc
					logCtx["damage"] = 10
					logCtx["attacker_hp"] = cli.Datas.Hp
					logCtx["target_hp"] = npc_datas.Stats.Hp

					return "OK", Datas, nil
				}
			}
		}
	}

	return "", "", errors.New("ERR 404 NPC_NOT_FOUND")
}
