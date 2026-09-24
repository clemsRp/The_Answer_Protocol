package drawer

import (
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) drawCombatFightersZone() {
	// Draw frame
	dr.drawCombatFrame(
		vars.COMBAT_START_X, vars.COMBAT_START_Y,
		vars.COMBAT_LEFT_END_X, vars.COMBAT_FIGHTERS_END_Y,
	)

	tile := dr.app.Variables.Tileset_size
	zoom := dr.app.Variables.Zoom
	frame_size := 4

	opponents := dr.app.Variables.PanelsVariables.CombatState.Opponents
	team := dr.app.Variables.PanelsVariables.CombatState.Team

	// Calculate placement variables
	panel_start_x := vars.COMBAT_START_X * tile
	panel_start_y := vars.COMBAT_START_Y * tile
	panel_width := (vars.COMBAT_LEFT_END_X - vars.COMBAT_START_X) * tile
	panel_height := (vars.COMBAT_END_Y - vars.COMBAT_START_Y) * tile

	opp_case_width := panel_width / float32(len(opponents))
	team_case_width := panel_width / float32(len(team))

	// Draw opponents
	center := float32(opp_case_width-0.5*float32(frame_size)*vars.FRAME_WIDTH*zoom) / 2
	for index := range len(opponents) {
		opp_offset := float32(index) * opp_case_width
		dr.DrawRealWoodFrame(
			float32(panel_start_x+opp_offset+center),
			float32(panel_start_y)+tile,
			frame_size, frame_size,
			0.5, 2, false,
		)
	}

	// Draw team
	center = float32(team_case_width-0.5*float32(frame_size)*vars.FRAME_WIDTH*zoom) / 2
	for index := range len(team) {
		team_offset := float32(index) * team_case_width
		dr.DrawRealWoodFrame(
			float32(panel_start_x+team_offset+center),
			float32(panel_start_y+panel_height)-(float32(frame_size)+4.5)*tile,
			frame_size, frame_size,
			0.5, 2, false,
		)
	}

	// Draw Combat Fighter
	for _, item := range dr.app.Manager.Interactions("Combat") {
		if item.Emote != nil {
			frame := item.Emote.CurrentFrame(dr.app.Variables.StartTime)
			dr.DrawImage(
				item.Emote.Texture,
				item.Emote.X-0.1*tile, item.Emote.Y-0.1*tile,
				frame.IndX, frame.IndY,
				frame.RatioX, frame.RatioY,
				item.Emote.Zoom/1.5,
				item.Emote.Rotation,
			)
		}

		pseudo := dr.LimitString(item.ID, 12)
		pseudo_len := rl.MeasureText(
			pseudo, int32(dr.app.Variables.FontSize),
		)
		offset := float32(frame_size)*vars.FRAME_WIDTH*zoom*1.5/4
		center = offset - float32(pseudo_len)/2

		color := dr.app.Colors["pseudo_text"]
		if item.Emote != nil {
			rl.DrawText(
				pseudo,
				int32(item.Emote.X-tile)+int32(center),
				int32(item.Emote.Y-0.8*tile),
				int32(dr.app.Variables.FontSize),
				color,
			)
		}
	}
}

