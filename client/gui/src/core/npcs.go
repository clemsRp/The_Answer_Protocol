package core

import (
	"tap/client/gui/src/ui"
	vars "tap/client/gui/src/variables"
	panel "tap/client/tui/panels"
	pr "tap/protocol"
)

type NpcDatas struct {
	texture        string
	ratioX, ratioY float32
	indXs, indYs   []float32
	pos            vars.Position
}

var (
	// TODO Define all npcs
	npc_convertor = map[string]NpcDatas{}
)

func (app *App) GetNewRoomNpcs(roomNpcs []string) []*ui.Interaction {
	// Initialize the interactions array
	npcs := make([]*ui.Interaction, 0)

	for _, np := range roomNpcs {
		// Set default placement and animation constants
		frame_duration := 100

		// Retrieve texture and frames for the current npc
		text, pos, frames := get_npc_datas(np, frame_duration)

		pos_x := pos.X * app.Variables.Tileset_size
		pos_y := pos.Y * app.Variables.Tileset_size

		// Create the visual emote (sprite) for the npc
		npc_emote := &ui.Emote{
			ID:           "npc_" + np,
			Texture:      text,
			X:            pos_x,
			Y:            pos_y,
			Zoom:         app.Variables.Zoom,
			Rotation:     0,
			AnimDuration: len(frames) * frame_duration,
			Frames:       frames,
		}

		// Create the background bubble for the interaction
		npc_bubble := &ui.Button{
			ID:       "npc_bubble",
			Texture:  vars.UI_SPRITE_TEXTURE,
			X:        pos_x + 0.8*app.Variables.Tileset_size,
			Y:        pos_y + 0.2*app.Variables.Tileset_size,
			Rotation: 0,
			Zoom:     app.Variables.Zoom / 2,
			Normal:   ui.Frame{IndX: 28, IndY: 8, RatioX: 3, RatioY: 3},
			Pressed:  ui.Frame{IndX: 28, IndY: 8, RatioX: 3, RatioY: 3},
			Hover:    ui.Frame{IndX: 28, IndY: 8, RatioX: 3, RatioY: 3},
			OnClick:  func() {},
		}

		// Create the interactive button to pick up the npc
		npc_btn := &ui.Button{
			ID:       "npc_btn",
			Texture:  vars.UI_SPRITE_TEXTURE,
			X:        pos_x + 1.15*app.Variables.Tileset_size,
			Y:        pos_y + 0.55*app.Variables.Tileset_size,
			Rotation: 0,
			Zoom:     app.Variables.Zoom * 2 / 5,
			Normal:   ui.Frame{IndX: 40, IndY: 8, RatioX: 2, RatioY: 2},
			Hover:    ui.Frame{IndX: 40, IndY: 8, RatioX: 2, RatioY: 2},
			Pressed:  ui.Frame{IndX: 42, IndY: 8, RatioX: 2, RatioY: 2},
			OnClick: func() {
				app.Variables.PanelsVariables.Chat.Open = false
				app.Variables.PanelsVariables.Group.Open = false
				app.Variables.PanelsVariables.Talk.LastTalk = np
				app.StartTalk(np)
				app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: pr.CmdTalk + " " + np}
				app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: pr.CmdInspectNpc + " " + np}
			},
		}

		// Group the emote and buttons into a single interaction entity
		npc := &ui.Interaction{
			ID:      "npc",
			Emote:   npc_emote,
			Buttons: []*ui.Button{npc_bubble, npc_btn},
		}

		npcs = append(npcs, npc)
	}

	return npcs
}

func get_npc_datas(npc string, frame_duration int) (string, vars.Position, []*ui.EmoteFrame) {
	frames := make([]*ui.EmoteFrame, 0)

	// Define needed datas
	var texture string
	var ratioX, ratioY float32
	var indXs, indYs []float32
	var pos vars.Position

	// Get datas depending on npc
	if d, ok := npc_convertor[npc]; ok {
		texture = d.texture
		indXs = d.indXs
		indYs = d.indYs
		ratioX = d.ratioX
		ratioY = d.ratioY
		pos = d.pos

	} else {
		texture = vars.NPC_TEXTURE
		indXs = []float32{0, 1, 2, 3, 4, 5, 6, 7}
		indYs = []float32{0, 0, 0, 0, 0, 0, 0, 0}
		ratioX = 1
		ratioY = 2
		pos = vars.Position{X: 6, Y: 8}
	}

	// Create emote frames
	for index := range indXs {
		frame := &ui.EmoteFrame{
			Frame:    &ui.Frame{IndX: indXs[index], IndY: indYs[index], RatioX: ratioX, RatioY: ratioY},
			Duration: frame_duration,
		}

		frames = append(frames, frame)
	}

	return texture, pos, frames
}
