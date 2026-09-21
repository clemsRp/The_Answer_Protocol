package core

import (
	vars "tap/client/gui/src/variables"
	"time"
)

func (app *App) StartTalk(npcID string) {
	app.Variables.PanelsVariables.Talk.Talking = &vars.Talk{NpcID: npcID, Start: time.Now()}
}

func (app *App) SetTalkResult(npcID, result string) {
	if app.Variables.PanelsVariables.Talk.Talking == nil || app.Variables.PanelsVariables.Talk.Talking.NpcID != npcID {
		return
	}
	app.Variables.PanelsVariables.Talk.Talking.ResultStart = time.Now()
	app.Variables.PanelsVariables.Talk.Talking.Result = result
	*app.Variables.PanelsVariables.Talk.Results = append(*app.Variables.PanelsVariables.Talk.Results, result)
}

func (app *App) EndTalk() {
	app.Variables.PanelsVariables.Talk.Talking = nil
}
