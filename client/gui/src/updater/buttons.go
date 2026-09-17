package updater

import (
	"slices"
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
	leftPanelBtn := &ui.Button{
		ID:      "open_left_panel",
		Texture: vars.UI_SPRITE_TEXTURE,
		X:       up.app.Variables.Tileset_size,
		Y:       3 * up.app.Variables.Tileset_size,
		Zoom:    up.app.Variables.Zoom * 0.5,
		Normal:  ui.Frame{IndX: 40, IndY: 6, RatioX: 2, RatioY: 2},
		Pressed: ui.Frame{IndX: 42, IndY: 6, RatioX: 2, RatioY: 2},
		OnClick: func() {
			up.app.Variables.PanelsVariables.LeftPanel.Open = !up.app.Variables.PanelsVariables.LeftPanel.Open
		},
	}

	chatBtn := &ui.Button{
		ID:      "open_chat",
		Texture: vars.UI_SPRITE_TEXTURE,
		X:       2 * up.app.Variables.Tileset_size,
		Y:       3 * up.app.Variables.Tileset_size,
		Zoom:    up.app.Variables.Zoom * 0.5,
		Normal:  ui.Frame{IndX: 40, IndY: 8, RatioX: 2, RatioY: 2},
		Pressed: ui.Frame{IndX: 42, IndY: 8, RatioX: 2, RatioY: 2},
		OnClick: func() {
			up.app.Variables.PanelsVariables.Chat.Open = !up.app.Variables.PanelsVariables.Chat.Open
		},
	}

	quitBtn := &ui.Button{
		ID:      "quit",
		Texture: vars.UI_SPRITE_TEXTURE,
		X:       3 * up.app.Variables.Tileset_size,
		Y:       3 * up.app.Variables.Tileset_size,
		Zoom:    up.app.Variables.Zoom * 0.5,
		Normal:  ui.Frame{IndX: 48, IndY: 10, RatioX: 2, RatioY: 2},
		Pressed: ui.Frame{IndX: 50, IndY: 10, RatioX: 2, RatioY: 2},
		OnClick: func() {
			up.app.ActionsChan <- panel.Action{
				Type:    panel.ActionSendServer,
				Payload: pr.CmdQuit,
			}
			up.app.Stop()
		},
	}

	up.app.Manager.SetViewButtons("Game", []*ui.Button{chatBtn, leftPanelBtn, quitBtn})
}

func (up *Updater) buildChatButtons() {
	sendchatBtn := &ui.Button{
		ID:      "send_chat",
		Texture: vars.UI_SPRITE_TEXTURE,
		X:       29.6 * up.app.Variables.Tileset_size,
		Y:       15.5 * up.app.Variables.Tileset_size,
		Zoom:    up.app.Variables.Zoom * 0.5,
		Normal:  ui.Frame{IndX: 48, IndY: 0, RatioX: 2, RatioY: 2},
		Pressed: ui.Frame{IndX: 50, IndY: 0, RatioX: 2, RatioY: 2},
		OnClick: up.SendChat,
	}

	scopes := []string{"GLOBAL", "ROOM", "GROUP"}

	scopeBtnPrevX := (vars.CHAT_START_X - 0.67) * up.app.Variables.Tileset_size
	scopeBtnNextX := (vars.CHAT_START_X + vars.CHAT_WIDTH - 0.37) * up.app.Variables.Tileset_size
	scopeBtnY := (vars.CHAT_START_Y + vars.CHAT_HEIGHT/2) * up.app.Variables.Tileset_size

	previousScopeBtn := &ui.Button{
		ID:      "previous_scope",
		Texture: vars.UI_SPRITE_TEXTURE,
		X:       scopeBtnPrevX,
		Y:       scopeBtnY,
		Zoom:    up.app.Variables.Zoom,
		Normal:  ui.Frame{IndX: 17, IndY: 1, RatioX: 1, RatioY: 1},
		Pressed: ui.Frame{IndX: 18, IndY: 1, RatioX: 1, RatioY: 1},
		OnClick: func() {
			cur_scope := up.app.Variables.PanelsVariables.Chat.CurrentScope
			index := slices.Index(scopes, cur_scope)
			up.app.Variables.PanelsVariables.Chat.CurrentScope = scopes[(index+1)%len(scopes)]
		},
	}

	nextScopeBtn := &ui.Button{
		ID:      "next_scope",
		Texture: vars.UI_SPRITE_TEXTURE,
		X:       scopeBtnNextX,
		Y:       scopeBtnY,
		Zoom:    up.app.Variables.Zoom,
		Normal:  ui.Frame{IndX: 17, IndY: 0, RatioX: 1, RatioY: 1},
		Pressed: ui.Frame{IndX: 18, IndY: 0, RatioX: 1, RatioY: 1},
		OnClick: func() {
			cur_scope := up.app.Variables.PanelsVariables.Chat.CurrentScope
			index := slices.Index(scopes, cur_scope) - 1
			if index < 0 {
				index = len(scopes) - 1
			}
			up.app.Variables.PanelsVariables.Chat.CurrentScope = scopes[index]
		},
	}

	up.app.Manager.SetViewButtons("Chat", []*ui.Button{sendchatBtn, previousScopeBtn, nextScopeBtn})
}
