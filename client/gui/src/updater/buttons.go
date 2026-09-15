package updater

import (
	"tap/client/gui/src/ui"
	vars "tap/client/gui/src/variables"
	panel "tap/client/tui/panels"
	pr "tap/protocol"
)

func (up *Updater) buildConnectButtons() {
	connect_start_x := 8.0
	connect_start_y := 4.5
	connect_width := 16.0
	connect_height := 9.0

	play_x := float32(connect_start_x+connect_width/2) * up.app.Variables.Tileset_size
	play_y := float32(connect_start_y+connect_height) * up.app.Variables.Tileset_size

	center_x := 3 * up.app.Variables.Tileset_size
	center_y := 3 * up.app.Variables.Tileset_size

	// Play buttons
	playBtn := &ui.Button{
		ID:      "play",
		Texture: vars.PLAY_TEXTURE,
		X:       play_x - center_x,
		Y:       play_y - center_y,
		Zoom:    up.app.Variables.Zoom,
		Normal:  ui.Frame{IndX: 0, IndY: 2, RatioX: 6, RatioY: 2},
		Pressed: ui.Frame{IndX: 6, IndY: 2, RatioX: 6, RatioY: 2},
		OnClick: func() {
			up.app.ActionsChan <- panel.Action{
				Type:    panel.ActionSendServer,
				Payload: pr.CmdConnect + " " + up.app.Variables.Player.Pseudo,
			}
		},
	}

	// Emote buttons
	emote_x := (float32(connect_start_x) + 0.5) * up.app.Variables.Tileset_size
	emote_y := (float32(connect_start_y) + 0.5) * up.app.Variables.Tileset_size

	prevBtn := &ui.Button{
		ID:      "previous_emote",
		Texture: vars.UI_SPRITE_TEXTURE,
		X:       emote_x + 0.5*up.app.Variables.Tileset_size,
		Y:       emote_y + 2*up.app.Variables.Tileset_size,
		Zoom:    up.app.Variables.Zoom,
		Normal:  ui.Frame{IndX: 17, IndY: 1, RatioX: 1, RatioY: 1},
		Pressed: ui.Frame{IndX: 18, IndY: 1, RatioX: 1, RatioY: 1},
		OnClick: func() {
			up.app.Variables.Player.EmoteIndex--
			if up.app.Variables.Player.EmoteIndex < 0 {
				up.app.Variables.Player.EmoteIndex = len(up.app.Manager.Emotes("Connect")) - 1
			}
		},
	}
	nextBtn := &ui.Button{
		ID:      "next_emote",
		Texture: vars.UI_SPRITE_TEXTURE,
		X:       emote_x + 4.5*up.app.Variables.Tileset_size,
		Y:       emote_y + 2*up.app.Variables.Tileset_size,
		Zoom:    up.app.Variables.Zoom,
		Normal:  ui.Frame{IndX: 17, IndY: 0, RatioX: 1, RatioY: 1},
		Pressed: ui.Frame{IndX: 18, IndY: 0, RatioX: 1, RatioY: 1},
		OnClick: func() {
			up.app.Variables.Player.EmoteIndex++
			up.app.Variables.Player.EmoteIndex %= len(up.app.Manager.Emotes("Connect"))
		},
	}

	up.app.Manager.SetViewButtons("Connect", []*ui.Button{playBtn, prevBtn, nextBtn})
}

func (up *Updater) buildGameButtons() {
	up.app.Manager.SetViewButtons("Game", []*ui.Button{})
}
