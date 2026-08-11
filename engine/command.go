package engine

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	pr "tap/protocol"
)

func (e *Engine) handleCmdConnect(ip string, req []string) (string, any, error) {
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidName)
	}

	// check if session is already connected
	if e.sessions[ip] != "" {
		return "", "", errors.New(pr.ErrNameInUse)
	}
	pseudo := req[1]
	if _, exists := e.players[pseudo]; exists {
		return "", "", errors.New(pr.ErrNameInUse)
	}

	player := &Player{
		name:   pseudo,
		room:   "entrance",
		hp:     100,
		hpMax:  100,
		ip:     ip,
		status: "idle",
	}
	e.players[pseudo] = player
	e.sessions[ip] = pseudo
	e.dialogues[pseudo] = make(map[string]int)

	e.inform_all(player, fmt.Sprintf("EVT STATS players=%d", len(e.players)))
	e.inform_room(player, "entrance", "EVT ROOM PRESENCE ENTER "+player.name)

	return "OK connected", "", nil
}

func (e *Engine) playerQuits(player *Player) {
	delete(e.players, player.name)
	e.inform_room(player, player.room, "EVT ROOM PRESENCE LEAVE "+player.name)
	e.inform_all(player, fmt.Sprintf("EVT STATS players=%d", len(e.players)))
}

func (e *Engine) handleCmdQuit(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	e.playerQuits(player)
	return "OK bye", "", nil
}

func (e *Engine) handleCmdWho(req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	return fmt.Sprintf("OK players=%d", len(e.players)), "", nil
}

func (e *Engine) handleCmdUsers(req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	// Get players
	users := make([]string, 0)
	for _, cli := range e.players {
		if cli.group == "" {
			users = append(users, cli.name)
		}
	}

	return "OK", pr.UsersCommandData{Users: users}, nil
}

func (e *Engine) handleCmdLook(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	players := make([]string, 0)

	// Get present players
	for pseudo, p := range e.players {
		if player.room == p.room {
			players = append(players, pseudo)
		}
	}

	currentRoom := e.world.Rooms[player.room]
	// Format response

	exits := pr.ExitsData{
		North: currentRoom.Exits["north"],
		South: currentRoom.Exits["south"],
		East:  currentRoom.Exits["east"],
		West:  currentRoom.Exits["west"],
	}
	Room := pr.RoomData{
		Id:          "room." + player.room,
		Name:        currentRoom.Name,
		Description: currentRoom.Description,
		Exits:       exits,
		Players:     players,
		Items:       currentRoom.Items,
		Npcs:        currentRoom.Npcs,
	}
	res := pr.LookCommandData{
		Room: Room,
	}

	return "OK", res, nil
}

func (e *Engine) handleCmdMove(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	direction := req[1]
	currentRoom := e.world.Rooms[player.room]

	// Check room is valid
	nextRoom, exists := currentRoom.Exits[direction]
	if !exists {
		return "", "", errors.New(pr.ErrNoExit)
	}

	e.inform_room(player, player.room, "EVT ROOM PRESENCE LEAVE "+player.name)

	// Change the player room variable
	if player, ok := e.players[player.ip]; ok {
		player.room = nextRoom
	}
	player.room = nextRoom

	e.inform_room(player, player.room, "EVT ROOM PRESENCE ENTER "+player.name)

	return fmt.Sprintf("OK room=%s", nextRoom), "", nil
}

func (e *Engine) handleCmdChat(player *Player, req []string) (string, any, error) {
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

	for pseudo, p := range e.players {
		if pseudo != player.name {

			// Compare scopes with other player datas
			is_global := scope == pr.GlobalChat
			is_group := scope == pr.GroupChat && player.group != "" && p.group == player.group
			is_room := scope == pr.RoomChat && p.room == player.room

			// Send chat to player
			if is_global || is_group || is_room {
				chat = fmt.Sprintf("EVT %s CHAT %s %s", scope, player.name, msg)
				e.exchanger.ServerOutput <- pr.EngineResponse{Ip: p.ip, Msg: chat}
			}
		}
	}
	return "OK", "", nil
}

func (e *Engine) handleCmdGroup(player *Player, req []string) (string, any, error) {
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
		res, err = e.create_group(player)
	case pr.InviteGroup:
		res, err = e.invite_user_in_group(player, arg)
	case pr.JoinGroup:
		res, err = e.join_group(player, arg)
	case pr.LeaveGroup:
		res, err = e.leave_group(player)
	case pr.PromoteGroup:
		res, err = e.promote_user(player, arg)
	case pr.AcceptPromoteGroup:
		res, err = e.accept_promotion(player)
	case pr.DeclinePromoteGroup:
		res, err = e.decline_promotion(player)
	default:
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	if err != nil {
		return "", "", err
	}

	return res, "", nil
}

func (e *Engine) handleCmdStatus(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	// Format response
	res := pr.StatusCommandData{
		Hp:     player.hp,
		MaxHp:  player.hpMax,
		Status: player.status,
	}

	return "OK", res, nil
}

func (e *Engine) handleCmdTake(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	object := req[1]
	for obj_index, obj := range e.world.Rooms[player.room].Items {
		if obj == object {
			player.inventory = append(player.inventory, object)
			e.world.Rooms[player.room].Items = append(e.world.Rooms[player.room].Items[:obj_index], e.world.Rooms[player.room].Items[obj_index+1:]...)

			e.inform_room(player, player.room, "EVT ITEM TOOK "+object)

			return "OK taken=" + object, "", nil
		}
	}

	return "", "", errors.New(pr.ErrItemNotFound)
}

func (e *Engine) handleCmdDrop(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	object := req[1]
	for obj_index, obj := range player.inventory {
		if obj == object {
			player.inventory = append(player.inventory[:obj_index], player.inventory[obj_index+1:]...)
			e.world.Rooms[player.room].Items = append(e.world.Rooms[player.room].Items, object)

			e.inform_room(player, player.room, "EVT ITEM DROPPED "+object)

			return "OK dropped=" + object, "", nil
		}
	}

	return "", "", errors.New(pr.ErrItemNotInInventory)
}

func (e *Engine) handleCmdInventory(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}
	return "OK", player.inventory, nil
}

func (e *Engine) handleCmdQuest(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	npc := req[1]
	for npc_name, npc_datas := range e.world.Npcs {
		if npc_name == npc {
			for _, room_npc := range e.world.Rooms[player.room].Npcs {
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

func (e *Engine) handleCmdQuests(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}
	res := player.quests

	return "OK", res, nil
}
func (e *Engine) handleCmdTalk(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	npc := req[1]
	// Find npc
	for npc_name, npc_datas := range e.world.Npcs {
		if npc_name == npc {
			for _, room_npc := range e.world.Rooms[player.room].Npcs {
				if room_npc == npc {
					// Get npc dialogue
					_, ok := e.dialogues[player.name][npc_name]
					if !ok {
						e.dialogues[player.name][npc_name] = 0
					}

					npc_index := e.dialogues[player.name][npc_name]
					Datas := npc_datas.Dialogue[npc_index%len(npc_datas.Dialogue)]
					e.dialogues[player.name][npc_name]++

					return "OK", Datas, nil
				}
			}
		}
	}

	return "", "", errors.New(pr.ErrNpcNotFound)
}

func (e *Engine) handleCmdAttack(player *Player, req []string) (string, any, error) {
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}
	npcName := req[1]

	if player.status == "combat" {
		return e.processCombatTurn(player, npcName)
	}

	return e.initiateCombat(player, npcName)
}
