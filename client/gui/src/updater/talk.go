package updater

import (
	"slices"
	panel "tap/client/tui/panels"
	pr "tap/protocol"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (up *Updater) UpdateTalk() {
	t := up.app.Variables.PanelsVariables.Talk
	talk := t.Talking
	if talk == nil || talk.Result == "" {
		return
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) || rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace) {
		// End animation
		if !t.Finished {
			t.Finished = true
			return
		}

		// Handle next dialogue
		if t.Results != nil && slices.Contains(*t.Results, talk.Result) {
			t.NbLoop++

		} else {
			// Reset value for next loop
			t.NbLoop = 0
		}

		// End dialogue
		if t.NbLoop > 1 {
			t.NbLoop = 0
			t.Finished = false
			talk.Result = ""
			*t.Results = make([]string, 0)
			up.app.EndTalk()

		} else {
			// Next dialogue
			t.Finished = false
			talk.Result = ""
			up.actionsChan <- panel.Action{
				Type:    panel.ActionSendServer,
				Payload: pr.CmdTalk + " " + t.Talking.NpcID,
			}
		}
	}
}
