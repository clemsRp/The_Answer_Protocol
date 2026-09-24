package variables

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Talk struct {
	NpcID       string
	Result      string
	Start       time.Time
	ResultStart time.Time
}

type TalkPanel struct {
	Talking  *Talk
	LastTalk string
	Finished bool
	NbLoop   int
	Results  *[]string
	Rect     rl.Rectangle
}
