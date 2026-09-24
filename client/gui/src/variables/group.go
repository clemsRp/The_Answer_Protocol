package variables

import rl "github.com/gen2brain/raylib-go/raylib"

type GroupPanel struct {
	Open            bool
	InGroup         bool
	IsLeader        bool
	Leader          string
	Grouped         []string
	UnGrouped       []string
	Invitations     []string
	SendInvitations []string
	SendPromotion   string
	Promote         bool

	LastNbOptions int

	Rect            rl.Rectangle
	ScrollRect      rl.Rectangle
	ScrollActive    bool
	Scroll          float32
	ScrollBarY      float32
	LastFrameScroll bool
}
