package gui

import (
	vars "tap/client/gui/src/variables"
	panel "tap/client/tui/panels"

	pr "tap/protocol"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (app *App) UpdateConnectView() {
	if rl.IsKeyPressed(rl.KeyEnter) {
		app.actionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: pr.CmdConnect + " ali"}
	}
}

func (app *App) UpdateGameView() {
	app.UpdatePlayer()
}

func (app *App) UpdatePlayer() {
	current_room := app.rooms[app.variables.Current_room]
	if current_room == nil || len(current_room.Collisions) == 0 || len(current_room.Collisions[0]) == 0 {
		return
	}

	tile_size := app.variables.Tileset_size
	if tile_size == 0 {
		tile_size = 16
	}

	dir_x := 0
	dir_y := 0

	var dx, dy float32
	if rl.IsKeyDown(rl.KeyDown) {
		dir_y = 1
		dy += 4
	}
	if rl.IsKeyDown(rl.KeyUp) {
		dir_y = -1
		dy -= 4
	}
	if rl.IsKeyDown(rl.KeyRight) {
		dir_x = 1
		dx += 4
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		dir_x = -1
		dx -= 4
	}

	app.variables.Player.Direction = &vars.Direction{X: float32(dir_x), Y: float32(dir_y)}

	if dx == 0 && dy == 0 {
		return
	}

	newX := app.variables.Player.Position.X + dx
	if app.canMove(current_room, newX, app.variables.Player.Position.Y, tile_size) {
		app.variables.Player.Position.X = newX
	}

	newY := app.variables.Player.Position.Y + dy
	if app.canMove(current_room, app.variables.Player.Position.X, newY, tile_size) {
		app.variables.Player.Position.Y = newY
	}

	if app.variables.Player.Position.X < 0 {
		app.variables.Player.Position.X = 0
	}
	if app.variables.Player.Position.Y < 0 {
		app.variables.Player.Position.Y = 0
	}
}
