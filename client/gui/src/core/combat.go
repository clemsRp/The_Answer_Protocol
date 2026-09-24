package core

import (
	"slices"
	"tap/client/gui/src/ui"
	vars "tap/client/gui/src/variables"
)

func (app *App) IsPlayerTurn() bool {
	cs := app.Variables.PanelsVariables.CombatState
	return cs != nil && cs.CurrentTurn != "" && cs.CurrentTurn == app.Variables.Player.Pseudo
}

func (app *App) RebuildCombatInteractions() {
	cs := app.Variables.PanelsVariables.CombatState
	if cs == nil {
		return
	}

	tile := app.Variables.Tileset_size
	zoom := app.Variables.Zoom
	frame_size := float32(4)

	panel_start_x := float32(vars.COMBAT_START_X) * tile
	panel_start_y := float32(vars.COMBAT_START_Y) * tile
	panel_width := float32(vars.COMBAT_LEFT_END_X-vars.COMBAT_START_X) * tile
	panel_height := float32(vars.COMBAT_END_Y-vars.COMBAT_START_Y) * tile

	opponents := cs.Opponents
	team := cs.Team

	interactions := make([]*ui.Interaction, 0)

	getEmoteForFighter := func(pseudo string, x, y float32) *ui.Emote {
		emoteIndex := 0
		if pseudo == app.Variables.Player.Pseudo {
			emoteIndex = app.Variables.Player.EmoteIndex
		} else if app.Variables.RemotePlayers != nil {
			if remote, ok := (*app.Variables.RemotePlayers)[pseudo]; ok {
				emoteIndex = remote.EmoteIndex
			}
		}

		emotesConnect := app.Manager.Emotes("Connect")
		if len(emotesConnect) == 0 {
			return nil
		}
		if emoteIndex < 0 || emoteIndex >= len(emotesConnect) {
			emoteIndex = 0
		}

		baseEmote := emotesConnect[emoteIndex]
		return &ui.Emote{
			ID:           "combat_" + pseudo,
			Texture:      baseEmote.Texture,
			X:            x,
			Y:            y,
			Zoom:         zoom * 0.75,
			AnimDuration: baseEmote.AnimDuration,
			Frames:       baseEmote.Frames,
		}
	}

	// Opponents
	if len(opponents) > 0 {
		opp_case_width := panel_width / float32(len(opponents))
		center := (opp_case_width - 0.5*frame_size*float32(vars.FRAME_WIDTH)*zoom) / 2

		keys := make([]string, 0, len(opponents))
		for k := range opponents {
			keys = append(keys, k)
		}
		slices.Sort(keys)

		for idx, name := range keys {
			opp_offset := float32(idx) * opp_case_width
			frameX := panel_start_x + opp_offset + center
			frameY := panel_start_y + tile

			emoteX := frameX + 0.35*tile
			emoteY := frameY + 0.35*tile

			fighterName := name
			interactions = append(interactions, &ui.Interaction{
				ID:    fighterName,
				Emote: getEmoteForFighter(fighterName, emoteX, emoteY),
				Inspect: func(target string) {
					if app.Variables.PanelsVariables.CombatState != nil {
						app.Variables.PanelsVariables.CombatState.SelectedPerson = target
					}
				},
			})
		}
	}

	// Team
	if len(team) > 0 {
		team_case_width := panel_width / float32(len(team))
		center := (team_case_width - 0.5*frame_size*float32(vars.FRAME_WIDTH)*zoom) / 2

		keys := make([]string, 0, len(team))
		for k := range team {
			keys = append(keys, k)
		}
		slices.Sort(keys)

		for idx, name := range keys {
			team_offset := float32(idx) * team_case_width
			frameX := panel_start_x + team_offset + center
			frameY := (panel_start_y + panel_height) - (frame_size+4.5)*tile

			emoteX := frameX + 0.35*tile
			emoteY := frameY + 0.35*tile

			fighterName := name
			interactions = append(interactions, &ui.Interaction{
				ID:    fighterName,
				Emote: getEmoteForFighter(fighterName, emoteX, emoteY),
				Inspect: func(target string) {
					if app.Variables.PanelsVariables.CombatState != nil {
						app.Variables.PanelsVariables.CombatState.SelectedPerson = target
					}
				},
			})
		}
	}

	app.Manager.SetViewInteractions("Combat", interactions)
}
