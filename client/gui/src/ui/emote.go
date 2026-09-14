package ui

import "time"

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

func (m *Manager) Emotes(view string) []*Emote {
	return m.views_emotes[view]
}

func (m *Manager) UpdateEmotes(view string) {
}
