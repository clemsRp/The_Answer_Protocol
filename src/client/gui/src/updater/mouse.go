package updater

import rl "github.com/gen2brain/raylib-go/raylib"

func (up *Updater) UpdateMouse() {
	mouse := rl.GetMousePosition()
	clicked := rl.IsMouseButtonPressed(rl.MouseButtonLeft)
	// down := rl.IsMouseButtonDown(rl.MouseButtonLeft)

	up.app.Variables.Mouse.Clicked = clicked
	// up.app.Variables.Mouse.Down = down
	up.app.Variables.Mouse.X = mouse.X
	up.app.Variables.Mouse.Y = mouse.Y
}
