package drawer

import (
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawCombatActionsZone() {
	dr.DrawCombatFrame(
		vars.COMBAT_START_X, vars.COMBAT_ACTIONS_START_Y,
		vars.COMBAT_LEFT_END_X, vars.COMBAT_END_Y,
	)

	if dr.app.IsPlayerTurn() {
		// TODO: Draw Actions
		return
	}

	tile := dr.app.Variables.Tileset_size
	font := int32(dr.app.Variables.FontSize)
	text := "Waiting"

	centerX := float32(vars.COMBAT_START_X+vars.COMBAT_LEFT_END_X) / 2 * tile
	centerY := float32(vars.COMBAT_ACTIONS_START_Y+vars.COMBAT_END_Y) / 2 * tile

	rl.DrawText(
		text,
		int32(centerX)-rl.MeasureText(text, font)/2,
		int32(centerY)-font/2,
		font, dr.app.Colors["pseudo_text"],
	)
}
