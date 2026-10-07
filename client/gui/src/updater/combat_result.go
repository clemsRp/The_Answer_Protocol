package updater

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (up *Updater) UpdateCombatResultView() {
	clicked := rl.IsMouseButtonPressed(rl.MouseButtonLeft)

	next := rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace) || clicked
	right_time := time.Since(up.app.Variables.LastCombatResultTime) >= 5*time.Second

	if next && right_time {
		up.app.ShowGamePage()
	}
}
