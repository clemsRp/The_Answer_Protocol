package core

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"tap/client/state"
	panel "tap/client/tui/panels"
	"tap/protocol"
	pr "tap/protocol"
	"time"
	"unicode/utf8"

	vars "tap/client/gui/src/variables"
)

func (app *App) ShowConnectPage() {
	app.Variables.Current_view = "Connect"
}

func (app *App) ShowGamePage() {
	app.Variables.Current_view = "Game"
	app.Variables.Player.Position.X = float32(15.5 * app.Variables.Tileset_size)
	app.Variables.Player.Position.Y = float32(8.5 * app.Variables.Tileset_size)
	newPosX := app.Variables.Player.Position.X
	newPosY := app.Variables.Player.Position.Y
	newDirX := app.Variables.Player.Direction.X
	newDirY := app.Variables.Player.Direction.Y
	emoteIndex := app.Variables.Player.EmoteIndex
	payload := fmt.Sprintf("%s %f %f %f %f %d", protocol.CmdNotifyPlayerPosition, newPosX, newPosY, newDirX, newDirY, emoteIndex)

	app.ActionsChan <- panel.Action{
		Type:    panel.ActionSendServer,
		Payload: payload,
	}
	app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: protocol.CmdGetPlayerPositions}
	app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: protocol.CmdGetItemPositions}
}

func (app *App) ShowPopupPage() {}
func (app *App) ClosePopup()    {}

func (app *App) ShowCombatPage() {
	app.EndTalk()
	app.Variables.PanelsVariables.Chat.Open = false
	app.Variables.PanelsVariables.Group.Open = false
	app.Variables.PanelsVariables.Inspect.Open = false
	app.Variables.Current_view = "Combat"
	app.RebuildCombatInteractions()
}

func (app *App) ShowCombatResultPopup(result string, rewards []string) {
	app.ActionsChan <- panel.Action{
		Type:    panel.ActionSendServer,
		Payload: pr.CmdInspectSelf,
	}
	app.Variables.PanelsVariables.Inspect.LastInspect = "SELF"
	app.Variables.PanelsVariables.Chat.CurrentScope = "GLOBAL"
	app.Variables.PanelsVariables.Chat.Open = false
	app.ShowGamePage()
}

func (app *App) ShowQuestCompletedPopup(questID, reward string) {
	for _, datas := range app.Variables.Npcs {
		if datas.QuestID == questID {
			datas.RequestedQuest = true
			datas.CompletedQuest = true
		}
	}
}

func (app *App) UpdateNavigation(room *protocol.LookCommandData) {
	app.Variables.PanelsVariables.Room = room
	app.Variables.Current_room = strings.SplitN(room.Id, "room.", 2)[1]
}

func (app *App) UpdateItems(roomItems, inventory []string) {
	app.Variables.PanelsVariables.RoomItems = &roomItems
	app.Variables.PanelsVariables.InventoryItems = &inventory

	items := app.GetNewRoomItems(roomItems)
	app.Manager.SetViewInteractions("Items", items)

	invent_buttons, invent_emotes := app.GetNewInventory(inventory)
	app.Manager.SetViewEmotes("Inventory", invent_emotes)
	app.Manager.SetViewButtons("Inventory", invent_buttons)

	combat_use_buttons := app.GetNewCombatInventory(inventory)
	app.Manager.SetViewButtons("CombatInventory", combat_use_buttons)
}

func (app *App) UpdateDatas(text string) {}
func (app *App) UpdateInteraction(roomNpcs, players []string, npcData map[string]protocol.InspectNPCData, npcDialogues map[string]string, groupMembers []string, quests []protocol.TrackedQuestData, completed_quests []string) {
	npcs := app.GetNewRoomNpcs(roomNpcs)
	app.Manager.SetViewInteractions("Npcs", npcs)

	npc := app.Variables.PanelsVariables.Talk.LastTalk
	if talk_res, ok := npcDialogues[npc]; ok {
		app.SetTalkResult(npc, talk_res)
	}
}

func (app *App) UpdateGroup(groupState state.GroupState) {
	app.QueueUpdate(func() {
		app.Variables.PanelsVariables.GroupState = &groupState
		app.UpdateGroupPanel(groupState)

		prev := ""
		if sel := app.Manager.Selects("Group"); len(sel) > 0 {
			prev = sel[0].CurrentOption
		}
		selects := app.GetGroupSelects()
		current := ""
		if len(selects) > 0 {
			for _, opt := range selects[0].Options {
				if opt == prev {
					selects[0].CurrentOption = prev
					break
				}
			}
			current = selects[0].CurrentOption
		}

		app.Manager.SetViewOptions("Group", app.GetGroupOptions(current))
		app.Manager.SetViewSelects("Group", selects)
	})
}

func (app *App) UpdateGroupPanel(grS state.GroupState) {
	gr := app.Variables.PanelsVariables.Group

	// Sync everything
	gr.InGroup = grS.Group != ""
	gr.IsLeader = grS.Leader
	if gr.IsLeader {
		gr.Leader = app.GetPseudo()
	}
	gr.Grouped = grS.Grouped
	gr.UnGrouped = grS.UnGrouped
	gr.Invitations = grS.Invitations
	gr.Promote = grS.Promotion

	if !grS.SendPromotion {
		if gr.SendPromotion != "" {
			gr.Leader = gr.SendPromotion
		}
		gr.SendPromotion = ""
	}

	// Drop a pending promotion
	if gr.SendPromotion != "" && (!slices.Contains(gr.Grouped, gr.SendPromotion) || !gr.IsLeader) {
		gr.SendPromotion = ""
	}

	// Drop invitations
	gr.SendInvitations = slices.DeleteFunc(gr.SendInvitations, func(p string) bool {
		return slices.Contains(gr.Grouped, p) || !slices.Contains(gr.UnGrouped, p)
	})

	if !gr.InGroup {
		gr.Promote = false
		gr.SendPromotion = ""
	}
}

func (app *App) UpdateCombat(combatState state.CombatState) {
	app.Variables.PanelsVariables.CombatState = &combatState
	app.RebuildCombatInteractions()
}

func (app *App) UpdateQuests(quests []protocol.TrackedQuestData) {
	app.Variables.PanelsVariables.Quests = &quests
	for _, datas := range app.Variables.Npcs {
		app.syncNpcQuest(datas)
	}
}

func (app *App) syncNpcQuest(datas *vars.NpcDatas) {
	if datas.QuestID == "" || app.Variables.PanelsVariables.Quests == nil {
		return
	}
	for _, q := range *app.Variables.PanelsVariables.Quests {
		if q.Id == datas.QuestID {
			datas.RequestedQuest = true
			datas.CompletedQuest = q.Status == "completed"
			return
		}
	}
}

func (app *App) AppendChat(scope, user, msg string) {
	scope_up := strings.ToUpper(scope)
	chats := app.Variables.PanelsVariables.Chat.ScopeChats

	emote_index := app.Variables.Player.EmoteIndex
	if user != app.Variables.Player.Pseudo {
		emote_index = 0
		if remotePlayer, exists := (*app.Variables.RemotePlayers)[user]; exists {
			emote_index = remotePlayer.EmoteIndex
		}
	}

	new_chat := vars.Chat{
		Msg:        msg,
		Pseudo:     user,
		Time:       time.Now(),
		EmoteIndex: emote_index,
	}

	chats[scope_up] = append(chats[scope_up], new_chat)
}

func (app *App) UpdateInspector(text string) {
	app.QueueUpdate(func() {
		re := regexp.MustCompile(`\[.+?\]`)
		text = re.ReplaceAllString(text, "")

		switch app.Variables.PanelsVariables.Inspect.LastInspect {
		case "ROOM":
			app.UpdateNpcWithRoom(text)
		case "NPC":
			app.UpdateNpcWithNpc(text)
		case "SELF":
			app.UpdateSelfHp(text)

		default:

		}

		maxLen := 30
		lines := strings.Split(text, "\n")
		var result strings.Builder

		for lineIndex, line := range lines {
			if lineIndex > 0 {
				result.WriteString("\n")
			}

			words := strings.Fields(line)
			if len(words) == 0 {
				continue
			}

			currentLineLen := 0
			for i, word := range words {
				wordLen := utf8.RuneCountInString(word)

				if i == 0 {
					result.WriteString(word)
					currentLineLen = wordLen
					continue
				}

				if currentLineLen+1+wordLen > maxLen {
					result.WriteString("\n")
					result.WriteString(word)
					currentLineLen = wordLen

				} else {
					result.WriteString(" ")
					result.WriteString(word)
					currentLineLen += 1 + wordLen
				}
			}
		}

		app.Variables.PanelsVariables.Inspect.Datas = result.String()
	})
}

func (app *App) UpdateNpcWithRoom(text string) {
	lines := strings.Split(text, "\n")
	inNpcSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "NPCS:") {
			inNpcSection = true
			continue
		} else if strings.HasPrefix(line, "PLAYERS:") || strings.HasPrefix(line, "ITEMS:") {
			inNpcSection = false
			continue
		}

		if inNpcSection && strings.HasPrefix(line, "- ") {
			rawName := strings.TrimPrefix(line, "- ")

			isHostile := strings.Contains(rawName, "(hostile)")
			isQuest := strings.Contains(rawName, "(quest)")

			npcName := strings.ReplaceAll(rawName, "(hostile)", "")
			npcName = strings.ReplaceAll(npcName, "(quest)", "")
			npcName = strings.TrimSpace(npcName)

			npcKey := app.Variables.NpcConvertor[npcName]
			if npcKey == "" {
				npcKey = npcName
			}

			datas, ok := app.Variables.Npcs[npcKey]
			if !ok {
				datas = &vars.NpcDatas{}
				app.Variables.Npcs[npcKey] = datas
			}

			datas.Hostile = isHostile
			datas.HasQuest = isQuest
		}
	}
}

func (app *App) UpdateNpcWithNpc(text string) {
	prefix := "NAME: "
	startIndex := strings.Index(text, prefix)

	var npc string
	if startIndex != -1 {
		startIndex += len(prefix)
		endIndex := strings.Index(text[startIndex:], "\n")

		if endIndex != -1 {
			npc = text[startIndex : startIndex+endIndex]

		} else {
			npc = text[startIndex:]
		}

		npc = strings.TrimSpace(npc)
	}

	npc = app.Variables.NpcConvertor[npc]

	datas, ok := app.Variables.Npcs[npc]
	if !ok {
		datas = &vars.NpcDatas{}
		app.Variables.Npcs[npc] = datas
	}

	datas.Hostile = strings.Contains(text, "HOSTILE: YES")
	if m := regexp.MustCompile(`QUEST ID:\s*(\S+)`).FindStringSubmatch(text); m != nil {
		datas.QuestID = m[1]
	}
	datas.HasQuest = datas.QuestID != ""
	app.syncNpcQuest(datas)
}

func (app *App) UpdateSelfHp(text string) {
	re := regexp.MustCompile(`HP:\s*(\d+)\s*/\s*(\d+)`)
	match := re.FindStringSubmatch(text)

	if len(match) == 3 {
		current, _ := strconv.ParseFloat(match[1], 32)
		max, _ := strconv.ParseFloat(match[2], 32)

		app.Variables.Player.Hp = int(current)
		app.Variables.Player.MaxHp = int(max)
	}
}

func (app *App) AppendCombatChat(user, msg string) {
	app.AppendChat("COMBAT", user, msg)
}

func (app *App) AppendServerResponse(res protocol.ServerResponse) {}
func (app *App) AppendCliMessage(text string)                     {}
func (app *App) AppendCliResponse(res protocol.ServerResponse)    {}

func (app *App) UpdateRemotePlayerPosition(pseudo string, x, y, dirX, dirY float32, emoteIndex int) {
	if pseudo == app.Variables.Player.Pseudo {
		return
	}

	if remotePlayer, exists := (*app.Variables.RemotePlayers)[pseudo]; exists {
		remotePlayer.Position.X = x
		remotePlayer.Position.Y = y
		remotePlayer.Direction.X = dirX
		remotePlayer.Direction.Y = dirY
		remotePlayer.EmoteIndex = emoteIndex
		remotePlayer.Zoom = app.Variables.Zoom
	} else {
		(*app.Variables.RemotePlayers)[pseudo] = &vars.Player{
			Pseudo:     pseudo,
			Position:   &vars.Position{X: x, Y: y},
			Direction:  &vars.Direction{X: dirX, Y: dirY},
			EmoteIndex: emoteIndex,
			Zoom:       app.Variables.Zoom,
		}
	}
}

func (app *App) AddRemotePlayer(pseudo string) {
	if _, exists := (*app.Variables.RemotePlayers)[pseudo]; !exists {
		(*app.Variables.RemotePlayers)[pseudo] = &vars.Player{
			Pseudo:    pseudo,
			Position:  &vars.Position{X: app.Variables.StartingPosX, Y: app.Variables.StartingPosY},
			Direction: &vars.Direction{X: 0, Y: 1},
			Zoom:      app.Variables.Zoom,
		}
	}
	app.UpdateRemotePlayerPosition(pseudo, app.Variables.StartingPosX, app.Variables.StartingPosY, 0, 1, 0)
}

func (app *App) RemoveRemotePlayer(pseudo string) {
	if app.Variables.RemotePlayers != nil {
		new_remote_players := make(map[string]*vars.Player)

		for player_pseudo, player := range *app.Variables.RemotePlayers {
			if player_pseudo != pseudo {
				new_remote_players[player_pseudo] = player
			}
		}

		app.Variables.RemotePlayers = &new_remote_players
	}
}

func (app *App) ResetRemotePlayers() {
	remote_players := make(map[string]*vars.Player)
	app.Variables.RemotePlayers = &remote_players
}

func (app *App) UpdateItemPosition(name string, x, y float32) {
	if item, exists := (*app.Variables.ItemPositions)[name]; exists {
		(*item).X = x
		(*item).Y = y
	} else {
		(*app.Variables.ItemPositions)[name] = &vars.Position{X: x, Y: y}
	}

	for _, item := range app.Manager.Interactions("Game") {
		if item.Emote.ID == "item_"+name {
			// Update emote position
			item.Emote.X = x
			item.Emote.Y = y

			// Update buttons positions
			item.Buttons[0].X = x - 0.2*app.Variables.Tileset_size
			item.Buttons[0].Y = y - 1.2*app.Variables.Tileset_size
			item.Buttons[1].X = x + 0.15*app.Variables.Tileset_size
			item.Buttons[1].Y = y - 0.85*app.Variables.Tileset_size
			break
		}
	}
}

func (app *App) GetPseudo() string {
	return app.Variables.Player.Pseudo
}

func (app *App) SetPseudo(pseudo string) {
	app.Variables.Player.Pseudo = pseudo
}

func (app *App) Stop() {
	app.Running = false
}
