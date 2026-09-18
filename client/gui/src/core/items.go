package core

import (
	"tap/client/gui/src/ui"
	vars "tap/client/gui/src/variables"
	panel "tap/client/tui/panels"
	pr "tap/protocol"
)

func (app *App) GetNewRoomItems(roomItems []string) []*ui.Interaction {
	items := make([]*ui.Interaction, 0)

	for ind, it := range roomItems {
		item_start_x := 18
		item_start_y := 9
		frame_duration := 500

		item_emote := &ui.Emote{
			ID:           "item_" + it,
			Texture:      vars.CHEST_TEXTURE,
			X:            float32(item_start_x+2*ind) * app.Variables.Tileset_size,
			Y:            float32(item_start_y) * app.Variables.Tileset_size,
			Zoom:         app.Variables.Zoom,
			Rotation:     0,
			AnimDuration: 2 * frame_duration,
			Frames: []*ui.EmoteFrame{
				{
					Frame:    &ui.Frame{IndX: 0.9, IndY: 1, RatioX: 1.2, RatioY: 1},
					Duration: frame_duration,
				},
				{
					Frame:    &ui.Frame{IndX: 3.9, IndY: 1, RatioX: 1.2, RatioY: 1},
					Duration: frame_duration,
				},
			},
		}

		item_bubble := &ui.Button{
			ID:       "item_bubble",
			Texture:  vars.UI_SPRITE_TEXTURE,
			X:        (float32(item_start_x+2*ind) - 0.2) * app.Variables.Tileset_size,
			Y:        (float32(item_start_y) - 1.2) * app.Variables.Tileset_size,
			Rotation: -90,
			Zoom:     app.Variables.Zoom / 2,
			Normal:   ui.Frame{IndX: 28, IndY: 8, RatioX: 3, RatioY: 3},
			Pressed:  ui.Frame{IndX: 28, IndY: 8, RatioX: 3, RatioY: 3},
			Hover:    ui.Frame{IndX: 28, IndY: 8, RatioX: 3, RatioY: 3},
			OnClick:  func() {},
		}
		item_btn := &ui.Button{
			ID:       "item_btn",
			Texture:  vars.UI_SPRITE_TEXTURE,
			X:        (float32(item_start_x+2*ind) + 0.15) * app.Variables.Tileset_size,
			Y:        (float32(item_start_y) - 0.85) * app.Variables.Tileset_size,
			Rotation: -90,
			Zoom:     app.Variables.Zoom * 2 / 5,
			Normal:   ui.Frame{IndX: 52, IndY: 8, RatioX: 2, RatioY: 2},
			Hover:    ui.Frame{IndX: 52, IndY: 8, RatioX: 2, RatioY: 2},
			Pressed:  ui.Frame{IndX: 54, IndY: 8, RatioX: 2, RatioY: 2},
			OnClick: func() {
				app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: pr.CmdTake + " " + it}
			},
		}

		item := &ui.Interaction{
			ID:      "item",
			Emote:   item_emote,
			Buttons: []*ui.Button{item_bubble, item_btn},
		}

		items = append(items, item)
	}

	return items
}

func (app *App) GetNewInventory(inventory []string) ([]*ui.Button, []*ui.Emote) {
	invent_buttons := make([]*ui.Button, 0)
	invent_emotes := make([]*ui.Emote, 0)

	pseudo := app.Variables.Player.Pseudo
	start_x := float32(19)

	if len(pseudo) <= vars.PSEUDO_MAX_CHAR/3 {
		start_x = 11
	} else if len(pseudo) <= vars.PSEUDO_MAX_CHAR*2/3 {
		start_x = 15
	}
	start_x += 2
	start_x *= app.Variables.Zoom / 2 * vars.FRAME_WIDTH
	start_x -= 0.1 * app.Variables.Tileset_size
	start_y := app.Variables.Tileset_size

	for ind, it := range inventory {
		frame_duration := 500

		item_emote := &ui.Emote{
			ID:           "item_" + it,
			Texture:      vars.CHEST_TEXTURE,
			X:            start_x + app.Variables.Tileset_size*(float32(ind)+0.115),
			Y:            0.45*app.Variables.Tileset_size + start_y,
			Zoom:         app.Variables.Zoom / 1.5,
			Rotation:     0,
			AnimDuration: 2 * frame_duration,
			Frames: []*ui.EmoteFrame{
				{
					Frame:    &ui.Frame{IndX: 0.9, IndY: 1, RatioX: 1.2, RatioY: 1},
					Duration: frame_duration,
				},
				{
					Frame:    &ui.Frame{IndX: 3.9, IndY: 1, RatioX: 1.2, RatioY: 1},
					Duration: frame_duration,
				},
			},
		}

		item_btn := &ui.Button{
			ID:       "item_btn",
			Texture:  vars.UI_SPRITE_TEXTURE,
			X:        start_x + app.Variables.Tileset_size*(float32(ind)+0.645),
			Y:        0.2*app.Variables.Tileset_size + start_y,
			Rotation: 0,
			Zoom:     app.Variables.Zoom / 5,
			Normal:   ui.Frame{IndX: 52, IndY: 10, RatioX: 2, RatioY: 2},
			Hover:    ui.Frame{IndX: 52, IndY: 10, RatioX: 2, RatioY: 2},
			Pressed:  ui.Frame{IndX: 54, IndY: 10, RatioX: 2, RatioY: 2},
			OnClick: func() {
				app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: pr.CmdDrop + " " + it}
			},
		}

		invent_buttons = append(invent_buttons, item_btn)
		invent_emotes = append(invent_emotes, item_emote)
	}

	return invent_buttons, invent_emotes
}
