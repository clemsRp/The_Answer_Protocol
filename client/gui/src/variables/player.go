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
}
