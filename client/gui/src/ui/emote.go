package ui

import (
	vars "tap/client/gui/src/variables"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type EmoteFrame struct {
	Frame    *Frame
	Duration int
}

type Emote struct {
	ID      string
	Texture string

	X, Y     float32
	Zoom     float32
	Rotation float32

	AnimDuration int
	Frames       []*EmoteFrame
	Inspect      func(emote string)
}

func (e *Emote) Rect() rl.Rectangle {
	if len(e.Frames) == 0 {
		return rl.NewRectangle(0, 0, 0, 0)
	}

	frame := e.Frames[0].Frame
	w := float32(vars.FRAME_WIDTH) * e.Zoom * frame.RatioX
	h := float32(vars.FRAME_HEIGHT) * e.Zoom * frame.RatioY

	return rl.NewRectangle(e.X, e.Y, w, h)
}

func (e *Emote) CurrentFrame(start_time time.Time) Frame {
	elapsedMillis := time.Since(start_time).Milliseconds()
	cur_time := elapsedMillis % int64(e.AnimDuration)

	total_anim_time := 0
	var error_frame *EmoteFrame
	for _, emote_frame := range e.Frames {
		error_frame = emote_frame
		total_anim_time += int(emote_frame.Duration)
		if int64(total_anim_time) >= cur_time {
			return *emote_frame.Frame
		}
	}

	return *error_frame.Frame
}

func (m *Manager) SetViewEmotes(view string, emotes []*Emote) {
	m.views_emotes[view] = emotes
}

func (m *Manager) Emotes(view string) []*Emote {
	return m.views_emotes[view]
}

func (m *Manager) UpdateEmotes(view string) {
	emotes := m.views_emotes[view]
	if len(emotes) == 0 {
		return
	}

	mouse := rl.GetMousePosition()
	down := rl.IsMouseButtonDown(rl.MouseButtonLeft)

	for _, emote := range emotes {
		// Check emote hover
		hovered := rl.CheckCollisionPointRec(mouse, emote.Rect())

		if hovered && down {
			emote.Inspect(emote.ID)
		}
	}
}
