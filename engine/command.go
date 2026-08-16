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

	e.inform_all(player, fmt.Sprintf("EVT STATS players=%d", len(e.players)))
	e.inform_room(player, player.room, "EVT ROOM PRESENCE ENTER "+player.name)

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

	for _, obj := range player.inventory {
		player.room.Items = append(player.room.Items, obj.Id)
		e.inform_room(player, player.room, "EVT ITEM DROPPED "+obj.Id)
	}
	player.inventory = make([]*Item, 0)

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

	return "OK", pr.UnGroupedCommandData{UnGrouped: users}, nil
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

	return "OK", pr.GroupedCommandData{Grouped: users}, nil
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

	// Format response

	exits := pr.ExitsData{
		North: player.room.Exits["north"],
		South: player.room.Exits["south"],
		East:  player.room.Exits["east"],
		West:  player.room.Exits["west"],
	}
	Room := pr.RoomData{
		Id:          "room." + player.room.Name,
		Name:        player.room.Name,
		Description: player.room.Description,
		Exits:       exits,
		Players:     players,
		Items:       player.room.Items,
		Npcs:        player.room.Npcs,
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

	// Check room is valid
	nextRoomName, exists := player.room.Exits[direction]
	if !exists {
		return "", "", errors.New(pr.ErrNoExit)
	}
	nextRoom := e.world.Rooms[nextRoomName]
	e.inform_room(player, player.room, "EVT ROOM PRESENCE LEAVE "+player.name)

	// Change the player room variable
	player.room = nextRoom

	e.inform_room(player, player.room, "EVT ROOM PRESENCE ENTER "+player.name)

	return fmt.Sprintf("OK room=%s", nextRoom.Name), "", nil
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

	object := req[1]
	for obj_index, item_name := range player.room.Items {
		item := e.world.Items[item_name]
		if item_name == object {
			player.inventory = append(player.inventory, item)
			player.room.Items = append(player.room.Items[:obj_index], player.room.Items[obj_index+1:]...)

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
		if obj.Id == object {
			player.inventory = append(player.inventory[:obj_index], player.inventory[obj_index+1:]...)
			player.room.Items = append(player.room.Items, object)

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

// Dans la commande attack tout en bas, je voudrais:

// 1. d'abord verifier que la target est valide, si le Npc est dans la room pour l'instanciation du combat, ou si le Npc est dans les Fighters du combat session. ensuite recuperer le Npc_copy de cette fonction ou une erreur.

// 2. Dans l'instanciation, arranger le slice fighters par rapport  a l'initiative de chacun

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
	if !exists {
		return e.initiateCombat(player, npcName)
	}
	if !combat_session.isFighterTurn(player) {
		return "", "", errors.New(pr.ErrNotYourTurnToPlay)
	}
	return combat_session.processCombatTurn(player, npc_copy)
}
