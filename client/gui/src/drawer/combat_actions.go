package drawer

import (
	"fmt"
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawCombatActionsZone() {
	splitX := float32(vars.COMBAT_START_X) + float32(vars.COMBAT_LEFT_END_X-vars.COMBAT_START_X)/2.0

	dr.DrawCombatFrame(
		vars.COMBAT_START_X, vars.COMBAT_ACTIONS_START_Y,
		vars.COMBAT_LEFT_END_X, vars.COMBAT_END_Y,
	)

	font := dr.app.Variables.FontSize

	if dr.app.Variables.PanelsVariables.CombatState.CurrentTurn == dr.app.Variables.Player.Pseudo {
		dr.DrawCombatActionsButtons()

		for _, b := range dr.app.Manager.Buttons("CombatActions") {
			label := ""
			if b.ID == "combat_attack" {
				label = "Attack"
			} else if b.ID == "combat_flee" {
				label = "Flee"
			}

			if label != "" {
				btnW := float32(vars.FRAME_WIDTH) * b.Zoom * b.Normal.RatioX
				btnH := float32(vars.FRAME_HEIGHT) * b.Zoom * b.Normal.RatioY
				centerX := b.X + btnW/2
				centerY := b.Y + btnH/2

				textWidth := rl.MeasureText(label, int32(font))
				rl.DrawText(
					label,
					int32(centerX)-textWidth/2,
					int32(centerY)-int32(font)/2,
					int32(font),
					dr.app.Colors["panel_text"],
				)
			}
		}

	} else {
		tile := dr.app.Variables.Tileset_size
		font := int32(dr.app.Variables.FontSize)
		text := "Waiting"

		centerX := float32(vars.COMBAT_START_X+vars.COMBAT_LEFT_END_X) / 3 * tile
		centerY := float32(vars.COMBAT_ACTIONS_START_Y+vars.COMBAT_END_Y) / 2 * tile

		rl.DrawText(
			text,
			int32(centerX)-rl.MeasureText(text, font)/2,
			int32(centerY)-font/2,
			font, dr.app.Colors["pseudo_text"],
		)
	}

	tile := dr.app.Variables.Tileset_size
	rightStartX := (splitX + 0.4) * tile
	startY := (vars.COMBAT_ACTIONS_START_Y + 0.5) * tile
	lineHeight := 1.15 * font

	cs := dr.app.Variables.PanelsVariables.CombatState
	persoDmg := 0
	groupDmg := 0
	teamCount := 1

	if cs != nil {
		persoDmg = cs.PersonalDamage
		groupDmg = cs.TotalGroupDamage
		if len(cs.Team) > 0 {
			teamCount = len(cs.Team)
		}
	}
	if groupDmg < persoDmg {
		groupDmg = persoDmg
	}
	avgDmg := float64(groupDmg) / float64(teamCount)

	rl.DrawRectangle(
		int32(rightStartX),
		int32(startY+0.5*tile),
		int32(tile/16),
		int32((vars.COMBAT_END_Y-vars.COMBAT_ACTIONS_START_Y-2)*tile),
		dr.app.Colors["group_text"],
	)

	gap := tile

	rl.DrawText(
		"Donnees Globales",
		int32(rightStartX+gap),
		int32(startY+1.5*gap),
		int32(1.1*font),
		dr.app.Colors["pseudo_text"],
	)

	rl.DrawText(
		fmt.Sprintf("Degats perso: %d", persoDmg),
		int32(rightStartX+gap),
		int32(startY+lineHeight+1.5*gap),
		int32(font),
		dr.app.Colors["panel_text"],
	)

	rl.DrawText(
		fmt.Sprintf("Degats groupe: %d", groupDmg),
		int32(rightStartX+gap),
		int32(startY+2*lineHeight+1.5*gap),
		int32(font),
		dr.app.Colors["panel_text"],
	)

	rl.DrawText(
		fmt.Sprintf("Moyenne / pers: %.1f", avgDmg),
		int32(rightStartX+gap),
		int32(startY+3*lineHeight+1.5*gap),
		int32(font),
		dr.app.Colors["group_text"],
	)
}
