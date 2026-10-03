package core

import (
	"fmt"
	panel "tap/client/tui/panels"
	"tap/protocol"
)

func (app *App) ShowConnectPage() {
	app.Variables.Current_view = "Connect"
}

func (app *App) ShowGamePage() {
	app.Variables.Current_view = "Game"
	app.Variables.Player.Position.X = app.Variables.StartingPosX
	app.Variables.Player.Position.Y = app.Variables.StartingPosY
	newPosX := app.Variables.Player.Position.X
	newPosY := app.Variables.Player.Position.Y
	newDirX := app.Variables.Player.Direction.X
	newDirY := app.Variables.Player.Direction.Y
	emoteIndex := app.Variables.Player.EmoteIndex
	payload := fmt.Sprintf("%s %f %f %f %f %d", protocol.CmdNotifyPlayerPosition, newPosX, newPosY, newDirX, newDirY, emoteIndex)

	app.QueueUpdate(func() {
		app.ActionsChan <- panel.Action{
			Type:    panel.ActionSendServer,
			Payload: payload,
		}
		app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: protocol.CmdGetPlayerPositions}
		app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: protocol.CmdGetPlayerEmotes}
		app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: protocol.CmdGetItemPositions}
	})
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
	app.QueueUpdate(func() {
		app.ActionsChan <- panel.Action{
			Type:    panel.ActionSendServer,
			Payload: protocol.CmdInspectSelf,
		}
		app.Variables.PanelsVariables.Inspect.Open = true
	})
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

func (app *App) ShowPopupPage() {}
func (app *App) ClosePopup()    {}

func (app *App) Stop() {
	app.Running = false
}
