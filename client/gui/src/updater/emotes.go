package updater

import (
	"strconv"
	"tap/client/gui/src/ui"
	vars "tap/client/gui/src/variables"
)

func (up *Updater) buildConnectEmotes() {
	connect_start_x := 10
	connect_start_y := 6

	emote_x := (float32(connect_start_x) + 0.5) * up.app.Variables.Tileset_size
	emote_y := (float32(connect_start_y) + 0.5) * up.app.Variables.Tileset_size

	lengths := []int{1, 2, 5, 4, 2, 2, 2, 2, 2, 2, 2, 2, 1, 2, 1}
	emotes := make([]*ui.Emote, 0)

	for ind_y, ind_x_max := range lengths {
		emote_frames := make([]*ui.EmoteFrame, 0)
		frame_duration := 300

		anim_duration := 0
		for ind_x := range ind_x_max {
			frame := ui.Frame{IndX: float32(2 * ind_x), IndY: float32(2 * ind_y), RatioX: 2, RatioY: 2}
			emote_frame := ui.EmoteFrame{
				Frame:    &frame,
				Duration: frame_duration,
			}
			anim_duration += frame_duration

			emote_frames = append(emote_frames, &emote_frame)
		}

		emote := ui.Emote{
			ID:           "emote" + strconv.Itoa(ind_y),
			Texture:      vars.EMOTES_TEXTURE,
			X:            emote_x,
			Y:            emote_y,
			Zoom:         up.app.Variables.Zoom,
			Rotation:     0,
			AnimDuration: anim_duration,
			Frames:       emote_frames,
		}

		emotes = append(emotes, &emote)
	}

	up.app.Manager.SetViewEmotes("Connect", emotes)
}

func (up *Updater) buildGameEmotes() {
	emote_x := 1.5 * up.app.Variables.Tileset_size
	emote_y := 1.5 * up.app.Variables.Tileset_size

	lengths := []int{1, 2, 5, 4, 2, 2, 2, 2, 2, 2, 2, 2, 1, 2, 1}
	emotes := make([]*ui.Emote, 0)

	for ind_y, ind_x_max := range lengths {
		emote_frames := make([]*ui.EmoteFrame, 0)
		frame_duration := 300

		anim_duration := 0
		for ind_x := range ind_x_max {
			frame := ui.Frame{IndX: float32(2 * ind_x), IndY: float32(2 * ind_y), RatioX: 2, RatioY: 2}
			emote_frame := ui.EmoteFrame{
				Frame:    &frame,
				Duration: frame_duration,
			}
			anim_duration += frame_duration

			emote_frames = append(emote_frames, &emote_frame)
		}

		emote := ui.Emote{
			ID:           "emote" + strconv.Itoa(ind_y),
			Texture:      vars.EMOTES_TEXTURE,
			X:            emote_x,
			Y:            emote_y,
			Zoom:         up.app.Variables.Zoom / 2,
			Rotation:     0,
			AnimDuration: anim_duration,
			Frames:       emote_frames,
		}

		emotes = append(emotes, &emote)
	}

	up.app.Manager.SetViewEmotes("Game", emotes)
}
