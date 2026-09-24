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
		up.SendNotif(oldX, oldY, oldDirX, oldDirY, false)

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

	up.UpdateMove()
	up.SendNotif(oldX, oldY, oldDirX, oldDirY, false)

	up.UpdatePlayersInspect()
}

func (up *Updater) UpdateMove() {
	const marge float32 = 1
	tileSize := up.app.Variables.Tileset_size
	pos := up.app.Variables.Player.Position
	dir := up.app.Variables.Player.Direction

	// Define room limits
	transitions := []struct {
		direction  string
		atBoundary bool
		isMoving   bool
		updatePos  func()
	}{
		{
			direction:  "east",
			atBoundary: pos.X >= (31-marge)*tileSize,
			isMoving:   dir.X == 1,
			updatePos:  func() { pos.X = marge * tileSize },
		},
		{
			direction:  "west",
			atBoundary: pos.X <= marge*tileSize,
			isMoving:   dir.X == -1,
			updatePos:  func() { pos.X = (31 - marge) * tileSize },
		},
		{
			direction:  "south",
			atBoundary: pos.Y >= (17-marge)*tileSize,
			isMoving:   dir.Y == 1,
			updatePos:  func() { pos.Y = marge * tileSize },
		},
		{
			direction:  "north",
			atBoundary: pos.Y <= marge*tileSize,
			isMoving:   dir.Y == -1,
			updatePos:  func() { pos.Y = (17 - marge) * tileSize },
		},
	}

	// Check player position to move
	for _, t := range transitions {
		if t.atBoundary && t.isMoving {
			up.actionsChan <- panel.Action{
				Type:    panel.ActionSendServer,
				Payload: pr.CmdMove + " " + t.direction,
			}
			t.updatePos()
			up.app.ResetRemotePlayers()
			break
		}
	}
}

func (up *Updater) SendNotif(oldX, oldY, oldDirX, oldDirY float32, first_notif bool) {
	newPosX := up.app.Variables.Player.Position.X
	newPosY := up.app.Variables.Player.Position.Y
	newDirX := up.app.Variables.Player.Direction.X
	newDirY := up.app.Variables.Player.Direction.Y
	emoteIndex := up.app.Variables.Player.EmoteIndex

	if oldX != newPosX || oldY != newPosY || oldDirX != newDirX || oldDirY != newDirY || first_notif {
		payload := fmt.Sprintf("%s %f %f %f %f %d", pr.CmdNotifyPlayerPosition, newPosX, newPosY, newDirX, newDirY, emoteIndex)

		up.actionsChan <- panel.Action{
			Type:    panel.ActionSendServer,
			Payload: payload,
		}
	}
}

func (up *Updater) UpdatePlayersInspect() {
	for _, p := range *up.app.Variables.RemotePlayers {
		mouse := rl.GetMousePosition()
		clicked := rl.IsMouseButtonPressed(rl.MouseButtonLeft)
		hover := rl.CheckCollisionPointRec(mouse, p.Rect())

		if hover && clicked {
			up.actionsChan <- panel.Action{
				Type:    panel.ActionSendServer,
				Payload: pr.CmdInspectPlayer + " " + p.Pseudo,
			}
		}
	}
}
