package core

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"tap/client/state"
	"tap/protocol"
	pr "tap/protocol"
	"time"
	"unicode/utf8"

	vars "tap/client/gui/src/variables"
)

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

func (app *App) UpdateItemPosition(name string, x, y float32) {
	if item, exists := (*app.Variables.ItemPositions)[name]; exists {
		(*item).X = x
		(*item).Y = y
	} else {
		(*app.Variables.ItemPositions)[name] = &vars.Position{X: x, Y: y}
	}

	for _, item := range app.Manager.Interactions("Game") {
		if item.Emote.ID == "item_"+name {
			item.Emote.X = x
			item.Emote.Y = y

			item.Buttons[0].X = x - 0.2*app.Variables.Tileset_size
			item.Buttons[0].Y = y - 1.2*app.Variables.Tileset_size
			item.Buttons[1].X = x + 0.15*app.Variables.Tileset_size
			item.Buttons[1].Y = y - 0.85*app.Variables.Tileset_size
			break
		}
	}
}

func (app *App) UpdateInteraction(roomNpcs, players []string, npcData map[string]protocol.InspectNPCData, npcDialogues map[string]string, groupMembers []string, quests []protocol.TrackedQuestData, completed_quests []string) {
	npcs := app.GetNewRoomNpcs(roomNpcs)
	app.Manager.SetViewInteractions("Npcs", npcs)

	npc := app.Variables.PanelsVariables.Talk.LastTalk
	if talk_res, ok := npcDialogues[npc]; ok {
		app.SetTalkResult(npc, talk_res)
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

func (app *App) UpdateStatus(datas any) {
	bytes, err := json.Marshal(datas)

	if err == nil {
		var statusData pr.StatusCommandData
		json.Unmarshal(bytes, &statusData)
		fmt.Println(statusData.Hp, statusData.MaxHp)
		app.Variables.Player.Hp = int(statusData.Hp)
		app.Variables.Player.MaxHp = int(statusData.MaxHp)
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

func (app *App) AppendCombatChat(user, msg string) {
	app.AppendChat("COMBAT", user, msg)
}

func (app *App) UpdateDatas(text string)                          {}
func (app *App) AppendServerResponse(res protocol.ServerResponse) {}
func (app *App) AppendCliMessage(text string)                     {}
func (app *App) AppendCliResponse(res protocol.ServerResponse)    {}
