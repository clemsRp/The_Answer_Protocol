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
	e.leave_group(player)
}

func (e *Engine) handleCmdQuit(player *Player, req []string) (string, any, error) {
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
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}
	return fmt.Sprintf("OK players=%d", len(e.players)), "", nil
}

func (e *Engine) handleCmdUnGrouped(player *Player, req []string) (string, any, error) {
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}
	users := make([]string, 0)
	for _, cli := range e.players {
		if cli.name != player.name && cli.group == "" && !slices.Contains(cli.invitations, player.group) {
			users = append(users, cli.name)
		}
	}
	return "OK", users, nil
}

func (e *Engine) handleCmdGrouped(player *Player, req []string) (string, any, error) {
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}
	users := make([]string, 0)
	for _, cli := range e.players {
		if cli.name != player.name && cli.group == player.group && player.group != "" {
			users = append(users, cli.name)
		}
	}
	return "OK", users, nil
}

func (e *Engine) handleCmdLook(player *Player, req []string) (string, any, error) {
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	players := make([]string, 0)
	for pseudo, p := range e.players {
		if player.room == p.room {
			players = append(players, pseudo)
		}
	}

	if len(players) == 0 {
		players = append(players, player.name)
	}

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
		Id:          "room." + player.room.Name,
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
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	direction := req[1]
	nextRoomName, exists := player.room.Exits[direction]
	if !exists {
		return "", "", errors.New(pr.ErrNoExit)
	}
	nextRoom := e.world.Rooms[nextRoomName]
	e.inform_room(player, player.room, "EVT ROOM PRESENCE LEAVE "+player.name)

	player.room = nextRoom
	e.inform_room(player, player.room, "EVT ROOM PRESENCE ENTER "+player.name)

	return fmt.Sprintf("OK room=%s", nextRoom.Name), "", nil
}

func (e *Engine) handleCmdChat(player *Player, req []string) (string, any, error) {
	if len(req) < 3 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	var chat string
	scope := strings.ToUpper(req[1])
	msg := strings.Join(req[2:], " ")

	if !slices.Contains([]string{pr.GlobalChat, pr.RoomChat, pr.GroupChat, pr.CombatChat}, scope) {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	for pseudo, p := range e.players {
		if pseudo != player.name {
			is_global := scope == pr.GlobalChat
			is_group := scope == pr.GroupChat && player.group != "" && p.group == player.group
			is_room := scope == pr.RoomChat && p.room == player.room
			is_combat := scope == pr.CombatChat && player.group != "" && p.group == player.group && p.inCombat

			if is_global || is_group || is_room || is_combat {
				chat = fmt.Sprintf("EVT %s CHAT %s %s", scope, player.name, msg)
				e.exchanger.ServerOutput <- pr.EngineResponse{Id: p.id, Msg: chat}
			}
		}
	}
	return "OK", "", nil
}

func (e *Engine) handleCmdGroup(player *Player, req []string) (string, any, error) {
	if len(req) < 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	var err error
	var res string
	scope := strings.ToUpper(req[1])

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

func (e *Engine) handleCmdCombat(player *Player, req []string) (string, any, error) {
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	if !player.inCombat {
		return "", "", errors.New(pr.ErrNotInCombat)
	}

	var err error
	var datas any
	var res string
	scope := strings.ToUpper(req[1])

	switch scope {
	case pr.StatsCombat:
		res, datas, err = e.get_combat_stats(player)
	default:
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	if err != nil {
		return "", "", err
	}

	return res, datas, nil
}

func (e *Engine) handleCmdStatus(player *Player, req []string) (string, any, error) {
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	res := pr.StatusCommandData{
		Hp:     player.stats.Hp,
		MaxHp:  player.stats.HpMax,
		Status: player.stats.Status,
	}
	return "OK", res, nil
}

func (e *Engine) handleCmdTake(player *Player, req []string) (string, any, error) {
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
			// Recompute quest progress now that the inventory changed, so
			// the player doesn't have to explicitly ask for it.
			e.refreshQuestProgress(player)
			return "OK taken=" + object, "", nil
		}
	}
	return "", "", errors.New(pr.ErrItemNotFound)
}

func (e *Engine) handleCmdDrop(player *Player, req []string) (string, any, error) {
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	object := req[1]
	for obj_index, obj := range player.inventory {
		if obj.Id == object {
			player.inventory = append(player.inventory[:obj_index], player.inventory[obj_index+1:]...)
			player.room.Items = append(player.room.Items, object)
			e.inform_room(player, player.room, "EVT ITEM DROPPED "+object)
			// Dropping the item can un-fulfil a quest's target item, so
			// recompute progress here too.
			e.refreshQuestProgress(player)
			return "OK dropped=" + object, "", nil
		}
	}
	return "", "", errors.New(pr.ErrItemNotInInventory)
}

func (e *Engine) handleCmdInventory(player *Player, req []string) (string, any, error) {
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

					// Check if player already has the quest
					for _, q := range player.quests {
						if q.Id == npc_datas.QuestId {
							return "", "", errors.New(pr.ErrNoQuestAvailable)
						}
					}

					quest := e.world.Quests[npc_datas.QuestId]

					// Give quest to player
					playerQuest := quest.Clone()
					playerQuest.Status = "active"
					playerQuest.Progress = "0/1" // Or default progress logic
					player.quests = append(player.quests, playerQuest)

					e.refreshQuestProgress(player)

					res := pr.QuestData{
						Id:          npc_datas.QuestId,
						Reward:      quest.Reward,
						Description: quest.Description,
						Status:      "active",
					}
					return "OK", res, nil
				}
			}
		}
	}
	return "", "", errors.New(pr.ErrNpcNotFound)
}

func (e *Engine) handleCmdQuests(player *Player, req []string) (string, any, error) {
	if len(req) != 1 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}
	var res []pr.TrackedQuestData
	for _, q := range player.quests {
		res = append(res, pr.TrackedQuestData{
			Id:          q.Id,
			Description: q.Description,
			Status:      q.Status,
			Progress:    q.Progress,
		})
	}
	if res == nil {
		res = []pr.TrackedQuestData{}
	}
	return "OK", res, nil
}

func (e *Engine) handleCmdCompleteQuest(player *Player, req []string) (string, any, error) {
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}
	questId := req[1]

	var playerQuest *Quest
	for _, q := range player.quests {
		if q.Id == questId {
			playerQuest = q
			break
		}
	}

	if playerQuest == nil || playerQuest.Status != "active" {
		return "", "", errors.New(pr.ErrNoQuestAvailable) // or a specific error
	}

	// Validate completion
	completed := false
	if playerQuest.TargetItem != "" {
		for _, item := range player.inventory {
			if item.Id == playerQuest.TargetItem {
				completed = true
				break
			}
		}
	} else if playerQuest.TargetNpc != "" {
		for _, npcName := range player.DefeatedNpcs {
			if npcName == playerQuest.TargetNpc {
				completed = true
				break
			}
		}
	} else {
		// If no target is specified, complete automatically
		completed = true
	}

	if !completed {
		return "", "", errors.New("ERR 406 QUEST_NOT_COMPLETED") // Custom error or existing one
	}

	playerQuest.Status = "completed"
	playerQuest.Progress = "1/1"

	// Give reward
	if playerQuest.Reward != "" {
		if item, exists := e.world.Items[playerQuest.Reward]; exists {
			player.inventory = append(player.inventory, item.Clone())
		}
	}

	if worldQuest, exists := e.world.Quests[playerQuest.Id]; exists {
		worldQuest.Status = "unavailable"
	}

	if playerQuest.TargetItem != "" {
		removeItemFromInventory(player, playerQuest.TargetItem)
		e.removeItemFromWorld(playerQuest.TargetItem)
		e.inform_all(player, "EVT ITEM REMOVED "+playerQuest.TargetItem)
	}

	e.inform_all(player, "EVT QUEST COMPLETED "+playerQuest.Id)

	res := pr.QuestData{
		Id:          playerQuest.Id,
		Description: playerQuest.Description,
		Status:      playerQuest.Status,
		Reward:      playerQuest.Reward,
	}

	return "OK", res, nil
}

func removeItemFromInventory(player *Player, itemId string) {
	for i, item := range player.inventory {
		if item.Id == itemId {
			player.inventory = append(player.inventory[:i], player.inventory[i+1:]...)
			return
		}
	}
}

func (e *Engine) removeItemFromWorld(itemId string) {
	for _, room := range e.world.Rooms {
		for i, it := range room.Items {
			if it == itemId {
				room.Items = append(room.Items[:i], room.Items[i+1:]...)
				break
			}
		}
	}
}

func (e *Engine) refreshQuestProgress(player *Player) {
	for _, q := range player.quests {
		if q.Status != "active" {
			continue
		}

		ready := false
		switch {
		case q.TargetItem != "":
			for _, item := range player.inventory {
				if item.Id == q.TargetItem {
					ready = true
					break
				}
			}
		case q.TargetNpc != "":
			for _, npcName := range player.DefeatedNpcs {
				if npcName == q.TargetNpc {
					ready = true
					break
				}
			}
		default:
			ready = true
		}

		if ready {
			q.Progress = "1/1"
		} else {
			q.Progress = "0/1"
		}
	}
}

func (e *Engine) handleCmdTalk(player *Player, req []string) (string, any, error) {
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	npc := req[1]
	for npc_name, npc_datas := range e.world.Npcs {
		if npc_name == npc {
			for _, room_npc := range player.room.Npcs {
				if room_npc == npc {
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

	targetName := req[1]

	combat_session, exists := e.activeCombats[player.stats.CombatId]

	var res string
	var fullResponse *FullTurnResponse

	if !exists {

		// Not an existing combat, try to initiate a new combat with NPC
		npc_copy, err := e.getValidTarget(player, targetName)
		if err != nil {
			return "", "", err
		}
		combat_session, res, fullResponse = e.initiateCombat(player, npc_copy)
	} else {
		// Player is already in combat
		npc_copy, err := e.getValidTarget(player, targetName)
		if err != nil {
			return "", "", err
		}

		if !combat_session.isFighterTurn(player) {
			return "", nil, errors.New(pr.ErrNotYourTurnToPlay)
		}
		res, fullResponse = combat_session.processCombatTurn(player, npc_copy)
	}

	if len(combat_session.Players) > 1 && fullResponse != nil {
		eventMsg, jsonErr := convertObjectToJson("EVT COMBAT UPDATE", fullResponse)
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

	if cs.State == StateCancelled {
		delete(e.activeCombats, cs.Id)
	} else if cs.TurnResponse != nil {
		eventMsg, _ := convertObjectToJson("EVT COMBAT UPDATE", cs.TurnResponse)
		e.inform_combat_players(cs, player, eventMsg)

		if cs.State != StateOngoing {
			e.end_combat(cs)
		}
	}

	msgEvent := fmt.Sprintf("EVT COMBAT ALLY_LEAVE_COMBAT %s", player.name)
	e.inform_combat_players(cs, player, msgEvent)
	return "OK", nil, nil
}

func (e *Engine) handleCmdUse(player *Player, req []string) (string, any, error) {
	if len(req) != 2 {
		return "", "", errors.New(pr.ErrInvalidCommand)
	}

	itemName := req[1]
	var itemToUse *Item
	itemIdx := -1
	for i, it := range player.inventory {
		if it.Id == itemName || it.Name == itemName {
			itemToUse = it
			itemIdx = i
			break
		}
	}

	if itemToUse == nil {
		return "", "", errors.New(pr.ErrItemNotInInventory)
	}

	if itemToUse.Type != "consumable" && itemToUse.Type != "weapon" {
		return "", "", errors.New(pr.ErrItemNotUsable)
	}

	var cs *CombatSession
	if player.inCombat {
		var exists bool
		cs, exists = e.activeCombats[player.stats.CombatId]
		if exists && !cs.isFighterTurn(player) {
			return "", "", errors.New(pr.ErrNotYourTurnToPlay)
		}
	}

	// Apply effect
	logAmount := 0
	switch itemToUse.Type {
	case "weapon":
		newWeapon := itemToUse.ConvertToWeapon()
		if newWeapon == nil {
			return "", "", errors.New(pr.ErrInternalServer)
		}
		player.equippedWeapon = newWeapon
	case "consumable":
		switch itemToUse.TargetStat {
		case "hp":
			player.stats.Hp += itemToUse.Amount
			if player.stats.Hp > player.stats.HpMax {
				player.stats.Hp = player.stats.HpMax
			}
			logAmount = itemToUse.Amount
		case "damage":
			player.stats.Damage += itemToUse.Amount
			logAmount = itemToUse.Amount
		case "shield":
			player.stats.Shield += itemToUse.Amount
			logAmount = itemToUse.Amount
		}
	}

	// Remove from inventory
	player.inventory = append(player.inventory[:itemIdx], player.inventory[itemIdx+1:]...)

	if player.inCombat && cs != nil {
		cs.TurnResponse = &FullTurnResponse{
			PlayerAction: ActionLog{
				ActorName: player.name, TargetName: player.name,
				Result: &CombatTurnResult{AttackerHp: player.stats.Hp, TargetHp: player.stats.Hp, Damage: -logAmount, Status: cs.State},
			},
			NpcReactions: []ActionLog{},
			CombatState:  cs.State,
		}

		cs.nextTurn()
		if cs.State == StateOngoing {
			cs.processNpcsTurn()
		}

		if len(cs.Players) > 1 {
			eventMsg, _ := convertObjectToJson("EVT COMBAT UPDATE", cs.TurnResponse)
			e.inform_combat_players(cs, player, eventMsg)
		}

		resResponse := cs.TurnResponse
		if cs.State != StateOngoing {
			e.end_combat(cs)
		}

		return "OK", resResponse, nil
	}
	return "OK", nil, nil
}
