package core

import (
	"fmt"
	vars "tap/src/client/gui/src/variables"
	panel "tap/src/client/tui/panels"
	"tap/src/protocol"
	"time"
)

func (app *App) ShowConnectPage() {
	app.Variables.Current_view = "Connect"
}

func (app *App) ShowGamePage() {
	app.Variables.LastRoomSpanwTime = time.Now()
	app.EndTalk()
	app.Variables.Current_view = "Game"
	app.Variables.PanelsVariables.Chat.CurrentScope = "GLOBAL"
	app.Variables.PanelsVariables.Chat.Open = false
	app.Variables.PanelsVariables.Inspect.Open = false
	app.Variables.PanelsVariables.Datas.Open = false
	app.Variables.PanelsVariables.Group.Open = false

	app.Variables.Player.Position.X = app.Variables.StartingPosX
	app.Variables.Player.Position.Y = app.Variables.StartingPosY

	newPosX := app.Variables.Player.Position.X
	newPosY := app.Variables.Player.Position.Y

	newTileX := newPosX / app.Variables.Tileset_size
	newTileY := newPosY / app.Variables.Tileset_size

	newDirX := app.Variables.Player.Direction.X
	newDirY := app.Variables.Player.Direction.Y
	emoteIndex := app.Variables.Player.EmoteIndex

	payload := fmt.Sprintf("%s %f %f %f %f %f %f %d", protocol.CmdNotifyPlayerPosition, newPosX, newPosY, newTileX, newTileY, newDirX, newDirY, emoteIndex)

	app.QueueUpdate(func() {
		app.ActionsChan <- panel.Action{
			Type:    panel.ActionSendServer,
			Payload: payload,
		}
		app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: protocol.CmdGetPlayerPositions}
		app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: protocol.CmdGetPlayerEmotes}
		app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: protocol.CmdGetItemPositions}

		app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: protocol.CmdInspectSelf}
		app.Variables.PanelsVariables.Inspect.LastInspect = "SELF"
	})

	app.Variables.LastRoomChange = time.Now()
}

func (app *App) ShowCombatPage() {
	app.EndTalk()
	app.Variables.PanelsVariables.Chat.Open = false
	app.Variables.PanelsVariables.Group.Open = false
	app.Variables.PanelsVariables.Inspect.Open = false
	app.Variables.Current_view = "Combat"
	app.RebuildCombatInteractions()
}

func (app *App) ShowCombatResultPopup(result string, rewards []string) {
	app.Variables.PanelsVariables.CombatResult = &vars.CombatResultPanel{
		LastTime: time.Now(),
		Result:   result,
		Rewards:  rewards,
	}

	app.QueueUpdate(func() {
		app.ActionsChan <- panel.Action{
			Type:    panel.ActionSendServer,
			Payload: protocol.CmdInspectSelf,
		}
		app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: protocol.CmdInspectSelf}
		app.Variables.PanelsVariables.Inspect.LastInspect = "SELF"
	})
	app.Variables.PanelsVariables.Chat.CurrentScope = "GLOBAL"
	app.Variables.PanelsVariables.Chat.Open = false
	app.Variables.Current_view = "CombatResult"
}

func (app *App) ShowQuestCompletedPopup(questID, reward string) {
	for _, datas := range app.Variables.Npcs {
		if datas.QuestID == questID {
			datas.RequestedQuest = true
			datas.CompletedQuest = true
		}
	}
}

func (app *App) ShowPopupPage() {}
func (app *App) ClosePopup()    {}

func (app *App) Stop() {
	app.Running = false
}
