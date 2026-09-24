package variables

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Direction struct {
	X float32
	Y float32
}

type Position struct {
	X float32
	Y float32
}

type Player struct {
	Direction     *Direction
	Position      *Position
	Zoom          float32
	Speed         int
	Pseudo        string
	EmoteIndex    int
	LastTimeTyped time.Time
	MaxHp         int
	Hp            int
}

func (p *Player) Rect() rl.Rectangle {
	w := float32(FRAME_WIDTH) * p.Zoom
	h := float32(FRAME_HEIGHT) * p.Zoom
	return rl.NewRectangle(p.Position.X, p.Position.Y, w, h)
}

func GetPlayerVariables() *Player {
	return &Player{
		Direction: &Direction{
			X: 0,
			Y: 1,
		},
		EmoteIndex: 2,
		MaxHp:      100,
		Hp:         100,
	}
}
