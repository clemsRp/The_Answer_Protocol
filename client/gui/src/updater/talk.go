package updater

import (
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
	hover := rl.CheckCollisionPointRec(mouse, t.Rect)

	next := hover && clicked

	if next || rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace) {
		// Skip animation
		if !t.Finished {
			t.Finished = true
			return
		}

		// Ask next phrase
		t.Finished = false
		npcID := talk.NpcID
		talk.Result = ""
		up.actionsChan <- panel.Action{
			Type:    panel.ActionSendServer,
			Payload: pr.CmdTalk + " " + npcID,
		}
	}
}