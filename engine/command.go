package engine

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	pr "tap/protocol"
)

func (e *Engine) handleCmdConnect(id string, req []string) (string, any, error) {
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidName)
	}

	// check if session is already connected
	if e.sessions[id] != "" {
		return "", "", errors.New(pr.ErrNameInUse)
	}
	pseudo := req[1]
	if _, exists := e.players[pseudo]; exists {
		return "", "", errors.New(pr.ErrNameInUse)
	}

	player, err := e.createNewPlayerInstance(pseudo, id)
	if err != nil {
		return "", "", err
	}
	e.players[pseudo] = player
	e.sessions[id] = pseudo
	e.dialogues[pseudo] = make(map[string]int)

	e.inform_all(player, fmt.Sprintf("%s%d", pr.EventStatsPlayers, len(e.players)))
	e.inform_room(player, player.room, fmt.Sprintf("%s %s ", pr.EventRoomPresenceEnter, player.name))

	return pr.ConnectCmdResponse, "", nil
}

func (e *Engine) playerQuits(player *Player) {
	delete(e.players, player.name)
	e.inform_room(player, player.room, fmt.Sprintf("%s %s ", pr.EventRoomPresenceLeave, player.name))
	e.inform_all(player, fmt.Sprintf("%s%d", pr.EventStatsPlayers, len(e.players)))
}

func (e *Engine) handleCmdQuit(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	for _, obj := range player.inventory {
		player.room.Items = append(player.room.Items, obj.Id)
		e.inform_room(player, player.room, fmt.Sprintf("%s %s", pr.EventItemDropped, obj.Id))
	}
	player.inventory = make([]*Item, 0)

	e.playerQuits(player)
	return pr.QuitCmdResponse, "", nil
}

func (e *Engine) handleCmdWho(req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	return fmt.Sprintf("%s%d", pr.WhoCmdResponse, len(e.players)), "", nil
}

func (e *Engine) handleCmdUnGrouped(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	// Get players
	users := make([]string, 0)
	for _, cli := range e.players {
		if cli.name != player.name && cli.group == "" && !slices.Contains(cli.invitations, player.group) {
			users = append(users, cli.name)
		}
	}

	return "OK", users, nil
}

func (e *Engine) handleCmdGrouped(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	// Get players
	users := make([]string, 0)
	for _, cli := range e.players {
		if cli.name != player.name && cli.group == player.group && player.group != "" {
			users = append(users, cli.name)
		}
	}

	return "OK", users, nil
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

	if len(players) == 0 {
		players = append(players, player.name)
	}

	// Format response

	exits := pr.ExitsData{
		North: player.room.Exits["north"],
		South: player.room.Exits["south"],
		East:  player.room.Exits["east"],
		West:  player.room.Exits["west"],
	}
	npcs := make([]string, 0)
	for _, npc := range player.room.Npcs {
		if !slices.Contains(player.DefeatedNpcs, npc) {
			npcs = append(npcs, npc)
		}
	}
	res := pr.LookCommandData{
		Id:          player.room.Id,
		Name:        player.room.Name,
		Description: player.room.Description,
		Exits:       exits,
		Players:     players,
		Items:       player.room.Items,
		Npcs:        npcs,
	}

	return "OK", res, nil
}

func (e *Engine) handleCmdMove(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}
	if player.inCombat {
		return "", "", errors.New(pr.ErrForbiddenInCombat)
	}

	direction := req[1]

	// Check room is valid
	nextRoomName, exists := player.room.Exits[direction]
	if !exists {
		return "", "", errors.New(pr.ErrNoExit)
	}
	nextRoom := e.world.Rooms[nextRoomName]
	e.inform_room(player, player.room, fmt.Sprintf("%s %s", pr.EventRoomPresenceLeave, player.name))

	// Change the player room variable
	player.room = nextRoom

	e.inform_room(player, player.room, fmt.Sprintf("%s %s", pr.EventRoomPresenceEnter, player.name))

	return fmt.Sprintf("%s%s", pr.MoveCmdResponse, nextRoom.Name), "", nil
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
				e.exchanger.ServerOutput <- pr.EngineResponse{Id: p.id, Msg: chat}
			}
		}
	}
	return "OK", "", nil
}

func (e *Engine) handleCmdGroup(player *Player, req []string) (string, any, error) {
	var err error
	var res string
	if len(req) <= 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)

	}
	scope := strings.ToUpper(req[1])

	// Handle invalid command
	contains_commands := []string{
		pr.CreateGroup,
		pr.LeaveGroup,
		pr.AcceptPromoteGroup,
		pr.DeclinePromoteGroup,
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
	case pr.KickGroup:
		res, err = e.kick_user_in_group(player, arg)
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
		Hp:     player.stats.Hp,
		MaxHp:  player.stats.HpMax,
		Status: player.stats.Status,
	}

	return "OK", res, nil
}

func (e *Engine) handleCmdTake(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}
	if player.inCombat {
		return "", "", errors.New(pr.ErrForbiddenInCombat)
	}

	object := req[1]
	for obj_index, item_name := range player.room.Items {
		item := e.world.Items[item_name]
		if item_name == object {
			player.inventory = append(player.inventory, item)
			player.room.Items = append(player.room.Items[:obj_index], player.room.Items[obj_index+1:]...)

			e.inform_room(player, player.room, fmt.Sprintf("%s %s", pr.EventItemTook, object))

			return pr.TakeCmdResponse + object, "", nil
		}
	}

	return "", "", errors.New(pr.ErrItemNotFound)
}

func (e *Engine) handleCmdDrop(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}
	if player.inCombat {
		return "", "", errors.New(pr.ErrForbiddenInCombat)
	}

	object := req[1]
	for obj_index, obj := range player.inventory {
		if obj.Id == object {
			player.inventory = append(player.inventory[:obj_index], player.inventory[obj_index+1:]...)
			player.room.Items = append(player.room.Items, object)

			e.inform_room(player, player.room, fmt.Sprintf("%s %s", pr.EventItemDropped, object))

			return pr.DropCmdResponse + object, "", nil
		}
	}

	return "", "", errors.New(pr.ErrItemNotInInventory)
}

func (e *Engine) handleCmdInventory(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	inventory := make([]string, 0)

	for _, inv := range player.inventory {
		inventory = append(inventory, inv.Id)
	}

	return "OK", inventory, nil
}

func (e *Engine) handleCmdQuest(player *Player, req []string) (string, any, error) {
	// Handle invalid command
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}
	if player.inCombat {
		return "", "", errors.New(pr.ErrForbiddenInCombat)
	}

	npc := req[1]
	for npc_name, npc_datas := range e.world.Npcs {
		if npc_name == npc {
			for _, room_npc := range player.room.Npcs {
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
	if player.inCombat {
		return "", "", errors.New(pr.ErrForbiddenInCombat)
	}

	npc := req[1]
	// Find npc
	for npc_name, npc_datas := range e.world.Npcs {
		if npc_name == npc {
			for _, room_npc := range player.room.Npcs {
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
	npc_copy, err := e.getValidTarget(player, npcName)
	if err != nil {
		return "", "", err
	}

	combat_session, exists := e.activeCombats[player.stats.CombatId]

	var res string
	var fullResponse *FullTurnResponse

	if !exists {
		combat_session, res, fullResponse = e.initiateCombat(player, npc_copy)
	} else {
		if !combat_session.isFighterTurn(player) {
			return "", nil, errors.New(pr.ErrNotYourTurnToPlay)
		}
		res, fullResponse = combat_session.processCombatTurn(player, npc_copy)
	}
	if len(combat_session.Players) > 1 {
		eventMsg, jsonErr := convertObjectToJson(pr.EventCombatUpdate, fullResponse)
		if jsonErr == nil {
			e.inform_combat_players(combat_session, player, eventMsg)
		} else {
			return "", nil, errors.New(pr.ErrInternalServer)
		}
	}

	if combat_session.State != StateOngoing {
		e.end_combat(combat_session)
	}

	return res, fullResponse, nil
}

func (e *Engine) handleCmdFlee(player *Player, req []string) (string, any, error) {
	if len(req) != 1 {
		return "", nil, errors.New(pr.ErrInvalidCommand)
	}
	cs, isInCombat := e.activeCombats[player.stats.CombatId]
	if !isInCombat {
		return "", nil, errors.New(pr.ErrNotInCombat)
	}
	err := cs.leaveCombat(player)
	if err != nil {
		return "", nil, err
	}
	combat_end := true
	for _, p := range cs.Players {
		if p.inCombat {
			combat_end = false
		}
	}
	if combat_end {
		cs.State = StateCancelled
	}

	msgEvent := fmt.Sprintf("%s%s", pr.EventCombatAllyLeaveCombat, player.name)
	e.inform_combat_players(cs, player, msgEvent)
	if len(cs.Players) > 1 {
		eventMsg, jsonErr := convertObjectToJson(pr.EventCombatUpdate, cs.TurnResponse)
		if jsonErr == nil {
			e.inform_combat_players(cs, player, eventMsg)
		} else {
			return "", nil, errors.New(pr.ErrInternalServer)
		}
	}
	return "OK", nil, nil

}
