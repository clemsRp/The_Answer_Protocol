package updater

import (
	panel "tap/client/tui/panels"
	pr "tap/protocol"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (up *Updater) UpdateConnectView() {
	if rl.IsKeyPressed(rl.KeyEnter) {
		up.actionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: pr.CmdConnect + " clement"}
	}

	up.UpdateGameView()
}
