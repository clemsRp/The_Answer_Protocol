package updater

import (
	"math"
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

func (up *Updater) canMove(room *parser.Map, x, y float32, tile_size int) bool {
	x_off := vars.FRAME_WIDTH * up.app.Variables.Zoom
	y_off := vars.FRAME_HEIGHT * up.app.Variables.Zoom
	marge := 2
	bottom_part := vars.FRAME_HEIGHT * 2 / 3

	up_left := Pos{X: x + float32(marge), Y: y + float32(marge+bottom_part)}
	up_right := Pos{X: x + x_off - float32(marge), Y: y + float32(marge+bottom_part)}
	down_left := Pos{X: x + float32(marge), Y: y + y_off - float32(marge)}
	down_right := Pos{X: x + x_off - float32(marge), Y: y + y_off - float32(marge)}

	edges := [][2]Pos{
		{up_left, up_right},
		{up_right, down_right},
		{down_right, down_left},
		{down_left, up_left},
	}

	cells := []Pos{up_left, up_right, down_right, down_left}

	for _, edge := range edges {
		start, end := edge[0], edge[1]
		steps := int(math.Round(math.Max(
			math.Abs(float64(end.X-start.X)),
			math.Abs(float64(end.Y-start.Y)),
		)))
		if steps == 0 {
			continue
		}
		stepX := (end.X - start.X) / float32(steps)
		stepY := (end.Y - start.Y) / float32(steps)
		for i := 1; i < steps; i++ {
			cells = append(cells, Pos{
				X: start.X + stepX*float32(i),
				Y: start.Y + stepY*float32(i),
			})
		}
	}

	for _, cell := range cells {
		if up.isColliding(room, cell.X, cell.Y, tile_size) {
			return false
		}
	}
	return true
}

func (up *Updater) isColliding(room *parser.Map, x, y float32, tile_size int) bool {
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

	localX := float32(int(x)%tile_size) / up.app.Variables.Zoom
	localY := float32(int(y)%tile_size) / up.app.Variables.Zoom

	for _, r := range rects {
		if localX >= r.Start.X && localX < r.End.X && localY >= r.Start.Y && localY < r.End.Y {
			return true
		}
	}

	return false
}
