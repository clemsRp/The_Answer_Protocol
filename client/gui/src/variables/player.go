package variables

import "time"

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
	Speed         int
	Pseudo        string
	EmoteIndex    int
	LastTimeTyped time.Time
	MaxHp         int
	Hp            int
}

func GetPlayerVariables() *Player {
	return &Player{
		Direction: &Direction{
			X: 0,
			Y: 1,
		},
		EmoteIndex: 2,
		MaxHp:      100,
		Hp:         56,
	}
}
