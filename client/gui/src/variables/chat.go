package variables

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Chat struct {
	Msg        string
	Pseudo     string
	EmoteIndex int
	Time       time.Time
}

type ChatPanel struct {
	Open            bool
	CurrentScope    string
	Msg             string
	LastMsgText     string
	LastMsgScope    string
	ScopeChats      map[string][]Chat
	LastNbChats     int
	Rect            rl.Rectangle
	ScrollRect      rl.Rectangle
	ScrollActive    bool
	Scroll          float32
	ScrollBarY      float32
	LastFrameScroll bool
}
