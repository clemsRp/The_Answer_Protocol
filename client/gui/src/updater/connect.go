package updater

import (
	vars "tap/client/gui/src/variables"
	panel "tap/client/tui/panels"
	pr "tap/protocol"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (up *Updater) UpdateConnectView() {
	up.app.Manager.Update("Connect")

	up.UpdatePseudo()

	if rl.IsKeyPressed(rl.KeyEnter) {
		up.actionsChan <- panel.Action{
			Type:    panel.ActionSendServer,
			Payload: pr.CmdConnect + " " + up.app.Variables.Player.Pseudo,
		}
		up.app.Variables.Player.Pseudo = ""
	}

	up.UpdateGameView()
}

func (up *Updater) UpdatePseudo() {
	pseudo := up.app.Variables.Player.Pseudo
	if rl.IsKeyPressed(rl.KeyBackspace) || rl.IsKeyPressedRepeat(rl.KeyBackspace) {
		up.app.Variables.Player.LastTimeTyped = time.Now()
		if len(pseudo) > 0 {
			runes := []rune(pseudo)
			pseudo = string(runes[:len(runes)-1])
		}
	}

	char := rl.GetCharPressed()

	for char > 0 {
		up.app.Variables.Player.LastTimeTyped = time.Now()
		if char >= 32 && len(pseudo) < vars.PSEUDO_MAX_CHAR {
			pseudo += string(char)
		}

		char = rl.GetCharPressed()
	}

	up.app.Variables.Player.Pseudo = pseudo
}
