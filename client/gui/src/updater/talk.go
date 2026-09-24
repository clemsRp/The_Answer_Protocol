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

	// Init Buttons/Emotes
	if len(up.app.Manager.Buttons("Talk")) == 0 {
		up.buildTalkButtons()
	}

	up.app.Manager.Update("Talk")

	mouse := rl.GetMousePosition()
	clicked := rl.IsMouseButtonPressed(rl.MouseButtonLeft)
	hover := rl.CheckCollisionPointRec(mouse, up.app.Variables.PanelsVariables.Talk.Rect)

	next := hover && clicked

	if next || rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace) {
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
