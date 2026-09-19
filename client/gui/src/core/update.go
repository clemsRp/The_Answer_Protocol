package core

import (
	"fmt"
	"strings"
	"tap/client/state"
	panel "tap/client/tui/panels"
	"tap/protocol"
	"time"

	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (app *App) ShowConnectPage() {
	app.Variables.Current_view = "Connect"
}

func (app *App) ShowGamePage() {
	app.Variables.Current_view = "Game"
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

func (app *App) ShowCombatPage()                                       {}
func (app *App) ShowPopupPage()                                        {}
func (app *App) ClosePopup()                                           {}
func (app *App) ShowCombatResultPopup(result string, rewards []string) {}
func (app *App) ShowQuestCompletedPopup(questID, reward string)        {}

func (app *App) UpdateNavigation(room *protocol.LookCommandData) {
	app.Variables.PanelsVariables.Room = room
	app.Variables.Current_room = strings.SplitN(room.Id, "room.", 2)[1]
}

func (app *App) UpdateItems(roomItems, inventory []string) {
	app.Variables.PanelsVariables.RoomItems = &roomItems
	app.Variables.PanelsVariables.InventoryItems = &inventory

	items := app.GetNewRoomItems(roomItems)
	app.Manager.SetViewInteractions("Game", items)

	invent_buttons, invent_emotes := app.GetNewInventory(inventory)
	app.Manager.SetViewEmotes("Inventory", invent_emotes)
	app.Manager.SetViewButtons("Inventory", invent_buttons)
}

func (app *App) UpdateDatas(text string) {}
func (app *App) UpdateInteraction(npcs, players []string, npcData map[string]protocol.InspectNPCData, npcDialogues map[string]string, groupMembers []string, quests []protocol.TrackedQuestData, completed_quests []string) {
}

func (app *App) UpdateGroup(groupState state.GroupState) {
	app.Variables.PanelsVariables.GroupState = &groupState
}

func (app *App) UpdateCombat(combatState state.CombatState) {
	app.Variables.PanelsVariables.CombatState = &combatState
}

func (app *App) UpdateQuests(quests []protocol.TrackedQuestData) {
	app.Variables.PanelsVariables.Quests = &quests
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

func (app *App) UpdateInspector(text string)                      {}
func (app *App) AppendCombatChat(user, msg string)                {}
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
	} else {
		(*app.Variables.RemotePlayers)[pseudo] = &vars.Player{
			Pseudo:     pseudo,
			Position:   &vars.Position{X: x, Y: y},
			Direction:  &vars.Direction{X: dirX, Y: dirY},
			EmoteIndex: emoteIndex,
		}
	}
}

func (app *App) AddRemotePlayer(pseudo string) {
	if _, exists := (*app.Variables.RemotePlayers)[pseudo]; !exists {
		(*app.Variables.RemotePlayers)[pseudo] = &vars.Player{
			Pseudo:    pseudo,
			Position:  &vars.Position{X: app.Variables.StartingPosX, Y: app.Variables.StartingPosY},
			Direction: &vars.Direction{X: 0, Y: 1},
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
			for _, button := range item.Buttons {
				button.X = x
				button.Y = y
			}
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
	app.closeOnce.Do(func() {
		if app.Textures != nil {
			app.Textures.UnloadTextures()
		}
		rl.CloseWindow()
	})
}
