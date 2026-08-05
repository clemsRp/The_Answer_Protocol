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
		return "", "", errors.New(pr.ErrInvalidName)
	}

	client, exists := e.clients[ip]
	if !exists {
		return "", "", errors.New(pr.ErrInternalServer)
	}

	if client.Datas.Connected {
		return "", "", errors.New(pr.ErrNameInUse)
	}

	for _, cli := range e.clients {
		if cli.Name == req[1] {
			return "", "", errors.New(pr.ErrNameInUse)
		}
	}
	// Update client variables
	client.Name = req[1]
	client.Datas.Connected = true
	e.dialogues[req[1]] = make(map[string]int)

	e.inform_all(client, fmt.Sprintf("EVT STATS players=%d", e.get_nb_connected_players()))
	e.inform_room(client, client.Datas.Room, "EVT ROOM PRESENCE ENTER "+client.Name)

	return "OK connected", "", nil
}

func (e *Engine) handleCmdQuit(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	e.inform_room(cli, cli.Datas.Room, "EVT ROOM PRESENCE LEAVE "+cli.Name)
	e.inform_all(cli, fmt.Sprintf("EVT STATS players=%d", e.get_nb_connected_players()-1))
	return "OK bye", "", nil
}

func (e *Engine) handleCmdWho(req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
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
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	players := make([]string, 0)

	// Get present players
	for ip := range e.clients {
		if e.clients[ip].Datas.Room == cli.Datas.Room && e.clients[ip].Datas.Connected {
			players = append(players, e.clients[ip].Name)
		}
	}

	currentRoom := e.world.Rooms[cli.Datas.Room]
	// Format response

	north, ok := currentRoom.Exits["north"]
	if !ok {
		north = ""
	}
	south, ok1 := currentRoom.Exits["south"]
	if !ok1 {
		south = ""
	}
	east, ok2 := currentRoom.Exits["east"]
	if !ok2 {
		east = ""
	}
	west, ok := currentRoom.Exits["west"]
	if !ok {
		west = ""
	}
	Exits := pr.ExitsData{
		North: north,
		South: south,
		West:  west,
		East:  east,
	}

	Room := pr.RoomData{
		Id:          "room." + cli.Datas.Room,
		Name:        currentRoom.Name,
		Description: currentRoom.Description,
		Exits:       Exits,
		Players:     players,
		Items:       currentRoom.Items,
		Npcs:        currentRoom.Npcs,
	}
	res := pr.LookCommandData{
		Room: Room,
	}

	return "OK", res, nil
}

func (e *Engine) handleCmdMove(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	direction := req[1]
	currentRoom := e.world.Rooms[cli.Datas.Room]

	// Check room is valid
	nextRoom, exists := currentRoom.Exits[direction]
	if !exists {
		return "", "", errors.New(pr.ErrNoExit)
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
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	var chat string
	scope := strings.ToUpper(req[1])
	msg := strings.Join(req[2:], " ")

	// Check scope exist
	if !slices.Contains([]string{pr.GlobalChat, pr.RoomChat, pr.GroupChat}, scope) {
		return "", "", errors.New(pr.ErrInvalidCommand)
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
		return "", "", errors.New(pr.ErrInvalidCommand)
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
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	if err != nil {
		return "", "", err
	}

	return res, "", nil
}

func (e *Engine) handleCmdStatus(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	// Format response
	res := pr.StatusCommandData{
		Hp:     cli.Datas.Hp,
		MaxHp:  cli.Datas.Max_hp,
		Status: cli.Datas.Status,
	}

	return "OK", res, nil
}

func (e *Engine) handleCmdTake(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	object := req[1]
	for obj_index, obj := range e.world.Rooms[cli.Datas.Room].Items {
		if obj == object {
			cli.Datas.Inventory = append(cli.Datas.Inventory, object)
			e.world.Rooms[cli.Datas.Room].Items = append(e.world.Rooms[cli.Datas.Room].Items[:obj_index], e.world.Rooms[cli.Datas.Room].Items[obj_index+1:]...)

			return "OK taken=" + object, "", nil
		}
	}

	return "", "", errors.New(pr.ErrItemNotFound)
}

func (e *Engine) handleCmdDrop(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	object := req[1]
	for obj_index, obj := range cli.Datas.Inventory {
		if obj == object {
			cli.Datas.Inventory = append(cli.Datas.Inventory[:obj_index], cli.Datas.Inventory[obj_index+1:]...)
			e.world.Rooms[cli.Datas.Room].Items = append(e.world.Rooms[cli.Datas.Room].Items, object)

			return "OK dropped=" + object, "", nil
		}
	}

	return "", "", errors.New(pr.ErrItemNotInInventory)
}

func (e *Engine) handleCmdInventory(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}
	return "OK", cli.Datas.Inventory, nil
}

func (e *Engine) handleCmdQuest(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	npc := req[1]
	for npc_name, npc_datas := range e.world.Npcs {
		if npc_name == npc {
			for _, room_npc := range e.world.Rooms[cli.Datas.Room].Npcs {
				if room_npc == npc {
					if npc_datas.QuestId == "" || e.world.Quests[npc_datas.QuestId].Status == "unavailable" {
						return "", "", errors.New(pr.ErrNoQuestAvailable)
					}

					// Format response
					quest := e.world.Quests[npc_datas.QuestId]
					res := pr.QuestData{

						Id:          npc_datas.QuestId,
						Reward:      quest.Reward,
						Description: quest.Description,
						Status:      quest.Status,
					}

					return "OK", res, nil
				}
			}
		}
	}
	return "", "", errors.New(pr.ErrNpcNotFound)
}

func (e *Engine) handleCmdQuests(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}
	res := cli.Datas.Quests

	return "OK", res, nil
}
func (e *Engine) handleCmdTalk(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
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

	return "", "", errors.New(pr.ErrNpcNotFound)
}

func (e *Engine) handleCmdAttack(cli *pr.Client, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
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

	return "", "", errors.New(pr.ErrNpcNotFound)
}
