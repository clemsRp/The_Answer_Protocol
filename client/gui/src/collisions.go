package gui

import (
	"tap/client/gui/src/parser"

	vars "tap/client/gui/src/variables"
)

type Pos vars.Position
type Dir vars.Direction

type Collision struct {
	Start Pos
	End   Pos
}

var (
	collisions_convertor = map[int][]Collision{
		0: {
			{Start: Pos{X: 8, Y: 0}, End: Pos{X: 16, Y: 16}},
		},
		1: {
			{Start: Pos{X: 0, Y: 0}, End: Pos{X: 8, Y: 16}},
		},
		2: {
			{Start: Pos{X: 0, Y: 0}, End: Pos{X: 16, Y: 8}},
		},
		3: {
			{Start: Pos{X: 0, Y: 8}, End: Pos{X: 16, Y: 16}},
		},
		4: {
			{Start: Pos{X: 0, Y: 0}, End: Pos{X: 16, Y: 16}},
		},
		5: {
			{Start: Pos{X: 0, Y: 6}, End: Pos{X: 6, Y: 10}},
			{Start: Pos{X: 6, Y: 10}, End: Pos{X: 10, Y: 16}},
		},
		6: {
			{Start: Pos{X: 6, Y: 6}, End: Pos{X: 16, Y: 10}},
			{Start: Pos{X: 6, Y: 10}, End: Pos{X: 10, Y: 16}},
		},
		7: {
			{Start: Pos{X: 6, Y: 0}, End: Pos{X: 10, Y: 10}},
			{Start: Pos{X: 0, Y: 6}, End: Pos{X: 6, Y: 10}},
		},
		8: {
			{Start: Pos{X: 6, Y: 0}, End: Pos{X: 10, Y: 10}},
			{Start: Pos{X: 10, Y: 6}, End: Pos{X: 16, Y: 10}},
		},
		9: {
			{Start: Pos{X: 0, Y: 6}, End: Pos{X: 16, Y: 10}},
		},
		10: {
			{Start: Pos{X: 6, Y: 0}, End: Pos{X: 10, Y: 16}},
		},
		11: {
			{Start: Pos{X: 5, Y: 6}, End: Pos{X: 16, Y: 10}},
		},
		12: {
			{Start: Pos{X: 0, Y: 6}, End: Pos{X: 11, Y: 10}},
		},
		13: {
			{Start: Pos{X: 6, Y: 5}, End: Pos{X: 10, Y: 16}},
		},
		14: {
			{Start: Pos{X: 6, Y: 0}, End: Pos{X: 10, Y: 11}},
		},
		15: {
			{Start: Pos{X: 0, Y: 3}, End: Pos{X: 16, Y: 10}},
		},
		16: {
			{Start: Pos{X: 0, Y: 3}, End: Pos{X: 11, Y: 10}},
		},
		17: {
			{Start: Pos{X: 5, Y: 3}, End: Pos{X: 16, Y: 10}},
		},
		20: {
			{Start: Pos{X: 0, Y: 3}, End: Pos{X: 10, Y: 10}},
			{Start: Pos{X: 6, Y: 10}, End: Pos{X: 10, Y: 16}},
		},
		21: {
			{Start: Pos{X: 6, Y: 3}, End: Pos{X: 16, Y: 10}},
			{Start: Pos{X: 6, Y: 10}, End: Pos{X: 10, Y: 16}},
		},
		22: {
			{Start: Pos{X: 6, Y: 0}, End: Pos{X: 10, Y: 10}},
			{Start: Pos{X: 0, Y: 3}, End: Pos{X: 6, Y: 10}},
		},
		23: {
			{Start: Pos{X: 6, Y: 0}, End: Pos{X: 10, Y: 10}},
			{Start: Pos{X: 10, Y: 3}, End: Pos{X: 16, Y: 10}},
		},
		24: {
			{Start: Pos{X: 0, Y: 0}, End: Pos{X: 8, Y: 16}},
			{Start: Pos{X: 8, Y: 8}, End: Pos{X: 16, Y: 16}},
		},
		25: {
			{Start: Pos{X: 0, Y: 8}, End: Pos{X: 16, Y: 16}},
			{Start: Pos{X: 8, Y: 0}, End: Pos{X: 16, Y: 8}},
		},
		26: {
			{Start: Pos{X: 0, Y: 0}, End: Pos{X: 8, Y: 16}},
			{Start: Pos{X: 8, Y: 0}, End: Pos{X: 16, Y: 8}},
		},
		27: {
			{Start: Pos{X: 0, Y: 0}, End: Pos{X: 16, Y: 8}},
			{Start: Pos{X: 8, Y: 8}, End: Pos{X: 16, Y: 16}},
		},
	}
)

func (app *App) canMove(room *parser.Map, x, y float32, tile_size int) bool {
	x_off := vars.FRAME_WIDTH * app.variables.Zoom
	y_off := vars.FRAME_HEIGHT * app.variables.Zoom
	marge := 2
	bottom_part := vars.FRAME_HEIGHT * 2 / 3

	// Define hitbox corners
	up_left := Pos{X: x + float32(marge), Y: y + float32(marge+bottom_part)}
	up_right := Pos{X: x + x_off - float32(marge), Y: y + float32(marge+bottom_part)}
	down_left := Pos{X: x + float32(marge), Y: y + y_off - float32(marge)}
	down_right := Pos{X: x + x_off - float32(marge), Y: y + y_off - float32(marge)}

	// Get all hitbox cells
	cells := []Pos{up_left, up_right, down_right, down_left}
	dirs := map[Pos]Dir{
		up_left:    {X: 1, Y: 0},
		up_right:   {X: 0, Y: 1},
		down_right: {X: -1, Y: 0},
		down_left:  {X: 0, Y: -1},
	}

	for index := range 4 {
		start := cells[index]
		end := cells[(index+1)%4]

		start_x := start.X
		start_y := start.Y

		dir := dirs[start]

		for start_x != end.X || start_y != end.Y {
			start_x += dir.X
			start_y += dir.Y
			cells = append(cells, Pos{X: start_x, Y: start_y})
		}
	}

	// Check collisions
	for _, cell := range cells {
		if app.isColliding(room, cell.X, cell.Y, tile_size) {
			return false
		}
	}

	return true
}

func (app *App) isColliding(room *parser.Map, x, y float32, tile_size int) bool {
	indX := int(x) / tile_size
	indY := int(y) / tile_size

	if indY < 0 || indY >= len(room.Collisions) || indX < 0 || indX >= len(room.Collisions[0]) {
		return true
	}

	if room.Collisions[indY][indX] == 0 {
		return false
	}
	gid := room.Collisions[indY][indX] - room.CollisionFirstGID

	rects, ok := collisions_convertor[gid]
	if !ok {
		rects = []Collision{
			{Start: Pos{X: 0, Y: 0}, End: Pos{X: float32(vars.FRAME_WIDTH), Y: float32(vars.FRAME_HEIGHT)}},
		}
	}

	localX := float32(int(x)%tile_size) / app.variables.Zoom
	localY := float32(int(y)%tile_size) / app.variables.Zoom

	for _, r := range rects {
		if localX >= r.Start.X && localX < r.End.X && localY >= r.Start.Y && localY < r.End.Y {
			return true
		}
	}

	return false
}
