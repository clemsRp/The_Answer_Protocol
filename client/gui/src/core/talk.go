package core

import (
	"slices"
	vars "tap/client/gui/src/variables"
	"time"
)

func (app *App) StartTalk(npcID string) bool {
	t := app.Variables.PanelsVariables.Talk
	if t.FirstPhrases == nil {
		t.FirstPhrases = make(map[string]string)
	}

	t.Talking = &vars.Talk{NpcID: npcID, Start: time.Now()}
	t.Finished = false
	*t.Results = (*t.Results)[:0]

	if first, ok := t.FirstPhrases[npcID]; ok {
		app.SetTalkResult(npcID, first)
		return true
	}
	return false
}

func (app *App) OnTalkResponse(npcName, dialogue string) {
	app.SetTalkResult(npcName, dialogue)
}

func (app *App) SetTalkResult(npcID, result string) {
	t := app.Variables.PanelsVariables.Talk
	if t.Talking == nil || t.Talking.NpcID != npcID {
		return
	}

	if slices.Contains(*t.Results, result) {
		app.EndTalk()
		return
	}

	if len(*t.Results) == 0 {
		if t.FirstPhrases == nil {
			t.FirstPhrases = make(map[string]string)
		}
		if _, ok := t.FirstPhrases[npcID]; !ok {
			t.FirstPhrases[npcID] = result
		}
	}

	*t.Results = append(*t.Results, result)
	t.Talking.Result = result
	t.Talking.ResultStart = time.Now()
	t.Finished = false
}

func (app *App) EndTalk() {
	t := app.Variables.PanelsVariables.Talk
	t.Talking = nil
	t.Finished = false
	*t.Results = (*t.Results)[:0]
}