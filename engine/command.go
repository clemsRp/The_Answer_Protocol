package engine

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	pr "tap/protocol"
)

func (e *Engine) handleCmdConnect(ip string, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New("ERR 400 Invalid name: shouldn't contain space character")
	}

	if e.clients[ip].Datas.Connected {
		return "", "", errors.New("ERR 403 User already connected")
	}
	for _, cli := range e.clients {
		if cli.Name == req[1] {
			return "", "", errors.New("ERR 201 NAME_IN_USE")
		}
	}

	// Update client variables
	e.clients[ip].Name = req[1]
	e.clients[ip].Datas.Connected = true
	e.dialogues[req[1]] = make(map[string]int)
	e.inform_all(e.clients[ip], fmt.Sprintf("EVT STATS players=%d", e.get_nb_connected_players()))
	e.inform_room(e.clients[ip], e.clients[ip].Datas.Room, "EVT ROOM PRESENCE ENTER "+e.clients[ip].Name)

	return "OK connected", "", nil
}

func (e *Engine) handleCmdQuit(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New("ERR 400 BAD_REQUEST")
	}

	e.inform_room(cli, cli.Datas.Room, "EVT ROOM PRESENCE LEAVE "+cli.Name)
	e.inform_all(cli, fmt.Sprintf("EVT STATS players=%d", e.get_nb_connected_players()-1))
	return "OK bye", "", nil
}

func (e *Engine) handleCmdWho(req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New("ERR 400 BAD_REQUEST")
	}

	// Calculate number of players
	nb_clients := 0
	for cli := range e.clients {
		if e.clients[cli].Datas.Connected {
			nb_clients++
		}
	}

	return fmt.Sprintf("OK players=%d", nb_clients), "", nil
}

func (e *Engine) handleCmdLook(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New("ERR 400 BAD_REQUEST")
	}

	res := make(map[string]any)
	Room := make(map[string]any)
	players := make([]string, 0)

	// Get present players
	for ip := range e.clients {
		if e.clients[ip].Datas.Room == cli.Datas.Room && e.clients[ip].Datas.Connected {
			players = append(players, e.clients[ip].Name)
		}
	}

	// Format response
	currentRoom := e.world.Rooms[cli.Datas.Room]
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

func (e *Engine) handleCmdMove(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New("ERR 400 BAD_REQUEST")
	}

	direction := req[1]
	currentRoom := e.world.Rooms[cli.Datas.Room]

	// Check room is valid
	nextRoom, exists := currentRoom.Exits[direction]
	if !exists {
		return "", "", errors.New("ERR 301 NO_EXIT")
	}

	e.inform_room(cli, cli.Datas.Room, "EVT ROOM PRESENCE LEAVE "+cli.Name)

	// Change the player room variable
	if cli, ok := e.clients[cli.Ip]; ok {
		cli.Datas.Room = nextRoom
	}
	cli.Datas.Room = nextRoom

	e.inform_room(cli, cli.Datas.Room, "EVT ROOM PRESENCE ENTER "+cli.Name)

	return fmt.Sprintf("OK room=%s", nextRoom), "", nil
}

func (e *Engine) handleCmdChat(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) < 3 {
		return "", "", errors.New("ERR 400 BAD_REQUEST")
	}

	var chat string
	scope := strings.ToUpper(req[1])
	msg := strings.Join(req[2:], " ")

	// Check scope exist
	if !slices.Contains([]string{pr.GlobalChat, pr.RoomChat, pr.GroupChat}, scope) {
		return "", "", errors.New("ERR 400 BAD_REQUEST")
	}

	for ip := range e.clients {
		if e.clients[ip].Name != cli.Name {

			// Compare scopes with other player datas
			is_global := scope == pr.GlobalChat
			is_group := scope == pr.GroupChat && cli.Datas.Group != "" && cli.Datas.Group == e.clients[ip].Datas.Group
			is_room := scope == pr.RoomChat && cli.Datas.Room == e.clients[ip].Datas.Room

			// Send chat to player
			if is_global || is_group || is_room {
				chat = fmt.Sprintf("EVT %s CHAT %s %s", scope, cli.Name, msg)
				e.clients[ip].Ch <- pr.ServerResponse{Msg: chat}
			}
		}
	}
	return "OK", "", nil
}

func (e *Engine) handleCmdGroup(cli *pr.Client, req []string) (string, any, error) {
	var err error
	var res string
	scope := strings.ToUpper(req[1])

	// Handle invalid command
	contains_commands := []string{
		pr.CreateGroup,
		pr.LeaveGroup,
		pr.AcceptPromoteGroup,
	}
	contains := slices.Contains(contains_commands, scope)
	if (len(req) != 3 && !contains) || (len(req) > 2 && contains) {
		return "", "", errors.New("ERR 400 BAD_REQUEST")
	}

	var arg string
	if !contains {
		arg = req[2]
	}

	// Handle group scopes
	switch scope {
	case pr.CreateGroup:
		res, err = e.create_group(cli)
	case pr.InviteGroup:
		res, err = e.invite_user_in_group(cli, arg)
	case pr.JoinGroup:
		res, err = e.join_group(cli, arg)
	case pr.LeaveGroup:
		res, err = e.leave_group(cli)
	case pr.PromoteGroup:
		res, err = e.promote_user(cli, arg)
	case pr.AcceptPromoteGroup:
		res, err = e.accept_promotion(cli)
	case pr.DeclinePromoteGroup:
		res, err = e.decline_promotion(cli)
	default:
		return "", "", errors.New("ERR 400 Invalid scope")
	}

	if err != nil {
		return "", "", err
	}

	return res, "", nil
}

func (e *Engine) handleCmdStatus(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New("ERR 400 BAD_REQUEST")
	}

	// Format response
	res := make(map[string]any)
	res["status"] = cli.Datas.Status
	res["max_hp"] = cli.Datas.Max_hp
	res["hp"] = cli.Datas.Hp

	return "OK", res, nil
}

func (e *Engine) handleCmdTake(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New("ERR 400 BAD_REQUEST")
	}

	object := req[1]
	for obj_index, obj := range e.world.Rooms[cli.Datas.Room].Items {
		if obj == object {
			cli.Datas.Inventory = append(cli.Datas.Inventory, object)
			e.world.Rooms[cli.Datas.Room].Items = append(e.world.Rooms[cli.Datas.Room].Items[:obj_index], e.world.Rooms[cli.Datas.Room].Items[obj_index+1:]...)

			return "OK taken=" + object, "", nil
		}
	}

	return "", "", errors.New("ERR 404 ITEM_NOT_FOUND")
}

func (e *Engine) handleCmdDrop(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New("ERR 400 BAD_REQUEST")
	}

	object := req[1]
	for obj_index, obj := range cli.Datas.Inventory {
		if obj == object {
			cli.Datas.Inventory = append(cli.Datas.Inventory[:obj_index], cli.Datas.Inventory[obj_index+1:]...)
			e.world.Rooms[cli.Datas.Room].Items = append(e.world.Rooms[cli.Datas.Room].Items, object)

			return "OK dropped=" + object, "", nil
		}
	}

	return "", "", errors.New("ERR 404 ITEM_NOT_IN_INVENTORY")
}

func (e *Engine) handleCmdInventory(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New("ERR 400 BAD_REQUEST")
	}
	return "OK", cli.Datas.Inventory, nil
}

func (e *Engine) handleCmdQuest(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New("ERR 400 BAD_REQUEST")
	}

	npc := req[1]
	for npc_name, npc_datas := range e.world.Npcs {
		if npc_name == npc {
			for _, room_npc := range e.world.Rooms[cli.Datas.Room].Npcs {
				if room_npc == npc {
					if npc_datas.QuestId == "" || e.world.Quests[npc_datas.QuestId].Status == "unavailable" {
						return "", "", errors.New("ERR 406 NO_QUEST_AVAILABLE")
					}

					// Format response
					quest := e.world.Quests[npc_datas.QuestId]
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

func (e *Engine) handleCmdQuests(req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New("ERR 400 BAD_REQUEST")
	}

	res := make([]map[string]string, 0)
	for quest_id, quest := range e.world.Quests {
		// Format response
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
func (e *Engine) handleCmdTalk(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New("ERR 400 BAD_REQUEST")
	}

	npc := req[1]
	// Find npc
	for npc_name, npc_datas := range e.world.Npcs {
		if npc_name == npc {
			for _, room_npc := range e.world.Rooms[cli.Datas.Room].Npcs {
				if room_npc == npc {
					// Get npc dialogue
					_, ok := e.dialogues[cli.Name][npc_name]
					if !ok {
						e.dialogues[cli.Name][npc_name] = 0
					}

					npc_index := e.dialogues[cli.Name][npc_name]
					Datas := npc_datas.Dialogue[npc_index%len(npc_datas.Dialogue)]
					e.dialogues[cli.Name][npc_name]++

					return "OK", Datas, nil
				}
			}
		}
	}

	return "", "", errors.New("ERR 404 NPC_NOT_FOUND")
}

func (e *Engine) handleCmdAttack(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New("ERR 400 BAD_REQUEST")
	}

	npc := req[1]
	for npc_name, npc_datas := range e.world.Npcs {
		if npc_name == npc {
			for _, room_npc := range e.world.Rooms[cli.Datas.Room].Npcs {
				if room_npc == npc {
					// Format response
					Datas := make(map[string]any)
					Datas["attacker_hp"] = cli.Datas.Hp
					Datas["target_hp"] = npc_datas.Stats.Hp
					Datas["damage"] = 10
					Datas["status"] = "combat"

					return "OK", Datas, nil
				}
			}
		}
	}

	return "", "", errors.New("ERR 404 NPC_NOT_FOUND")
}
