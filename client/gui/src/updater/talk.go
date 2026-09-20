package updater

import rl "github.com/gen2brain/raylib-go/raylib"

func (up *Updater) UpdateTalk() {
	talk := up.app.Variables.PanelsVariables.Talk.Talking
	if talk == nil || talk.Result == "" {
		return
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) || rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace) {
		up.app.EndTalk()
	}
}
