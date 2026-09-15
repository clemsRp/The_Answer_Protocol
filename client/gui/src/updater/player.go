package updater

import (
	"fmt"
	"math"
	vars "tap/client/gui/src/variables"
	panel "tap/client/tui/panels"
	pr "tap/protocol"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (up *Updater) UpdatePlayer() {
	current_room := up.app.Rooms[up.app.Variables.Current_room]
	if current_room == nil || len(current_room.Collisions) == 0 || len(current_room.Collisions[0]) == 0 {
		return
	}

	tile_size := up.app.Variables.Tileset_size
	if tile_size == 0 {
		tile_size = 16
	}

	oldX := up.app.Variables.Player.Position.X
	oldY := up.app.Variables.Player.Position.Y
	oldDirX := up.app.Variables.Player.Direction.X
	oldDirY := up.app.Variables.Player.Direction.Y
	speed := float32(up.app.Variables.Player.Speed)

	dir_x := 0
	dir_y := 0

	var dx, dy float32
	if rl.IsKeyDown(rl.KeyDown) {
		dir_y = 1
		dy += speed
	} else if rl.IsKeyDown(rl.KeyUp) {
		dir_y = -1
		dy -= speed
	}

	if rl.IsKeyDown(rl.KeyRight) {
		dir_x = 1
		dx += speed
	} else if rl.IsKeyDown(rl.KeyLeft) {
		dir_x = -1
		dx -= speed
	}

	up.app.Variables.Player.Direction = &vars.Direction{X: float32(dir_x), Y: float32(dir_y)}

	if dx == 0 && dy == 0 {
		newPosX := up.app.Variables.Player.Position.X
		newPosY := up.app.Variables.Player.Position.Y
		newDirX := up.app.Variables.Player.Direction.X
		newDirY := up.app.Variables.Player.Direction.Y

		if oldX != newPosX || oldY != newPosY || oldDirX != newDirX || oldDirY != newDirY {
			payload := fmt.Sprintf("%s %f %f %f %f", pr.CmdNotifyPosition, newPosX, newPosY, newDirX, newDirY)

			up.actionsChan <- panel.Action{
				Type:    panel.ActionSendServer,
				Payload: payload,
			}
		}
	} else if dir_x != 0 && dir_y != 0 {
		dx *= float32(math.Sqrt(0.5))
		dy *= float32(math.Sqrt(0.5))
	}

	newX := up.app.Variables.Player.Position.X + dx
	if up.canMove(current_room, newX, up.app.Variables.Player.Position.Y, int(tile_size)) {
		up.app.Variables.Player.Position.X = newX
	}

	newY := up.app.Variables.Player.Position.Y + dy
	if up.canMove(current_room, up.app.Variables.Player.Position.X, newY, int(tile_size)) {
		up.app.Variables.Player.Position.Y = newY
	}

	if up.app.Variables.Player.Position.X < 0 {
		up.app.Variables.Player.Position.X = 0
	}
	if up.app.Variables.Player.Position.Y < 0 {
		up.app.Variables.Player.Position.Y = 0
	}

	newPosX := up.app.Variables.Player.Position.X
	newPosY := up.app.Variables.Player.Position.Y
	newDirX := up.app.Variables.Player.Direction.X
	newDirY := up.app.Variables.Player.Direction.Y

	if oldX != newPosX || oldY != newPosY || oldDirX != newDirX || oldDirY != newDirY {
		payload := fmt.Sprintf("%s %f %f %f %f", pr.CmdNotifyPosition, newPosX, newPosY, newDirX, newDirY)

		up.actionsChan <- panel.Action{
			Type:    panel.ActionSendServer,
			Payload: payload,
		}
	}
}
