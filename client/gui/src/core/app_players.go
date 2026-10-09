package core

import (
	"slices"
	vars "tap/client/gui/src/variables"
	"tap/client/state"
	"time"
)

func (app *App) GetPseudo() string {
	return app.Variables.Player.Pseudo
}

func (app *App) SetPseudo(pseudo string) {
	app.Variables.Player.Pseudo = pseudo
}

func (app *App) UpdateWho(nb_players int) {
	app.Variables.PanelsVariables.Datas.NbServerPlayers = nb_players
}

func (app *App) UpdateRemotePlayerPosition(pseudo string, x, y, tile_x, tile_y, dirX, dirY float32, emoteIndex int) {
	if pseudo == app.Variables.Player.Pseudo {
		return
	}

	potential_tile_x := x / app.Variables.Tileset_size
	potential_tile_y := y / app.Variables.Tileset_size

	if potential_tile_x != tile_x {
		x = tile_x * app.Variables.Tileset_size
	}
	if potential_tile_y != tile_y {
		y = tile_y * app.Variables.Tileset_size
	}

	if remotePlayer, exists := (*app.Variables.RemotePlayers)[pseudo]; exists {
		remotePlayer.Position.X = x
		remotePlayer.Position.Y = y
		remotePlayer.Direction.X = dirX
		remotePlayer.Direction.Y = dirY
		remotePlayer.EmoteIndex = emoteIndex
		remotePlayer.Zoom = app.Variables.Zoom
	} else {
		(*app.Variables.RemotePlayers)[pseudo] = &vars.Player{
			Pseudo:     pseudo,
			Position:   &vars.Position{X: x, Y: y},
			Direction:  &vars.Direction{X: dirX, Y: dirY},
			EmoteIndex: emoteIndex,
			Zoom:       app.Variables.Zoom,
		}
	}
}

func (app *App) UpdateRemotePlayerEmotes(pseudo string, emote int) {
	if pseudo == app.Variables.Player.Pseudo {
		return
	}

	if remotePlayer, exists := (*app.Variables.RemotePlayerEmotes)[pseudo]; exists {
		remotePlayer.Emote = emote
		remotePlayer.ChoiceTime = time.Now()

	} else {
		(*app.Variables.RemotePlayerEmotes)[pseudo] = &vars.PlayerEmote{
			Emote:      emote,
			ChoiceTime: time.Now(),
		}
	}
}

func (app *App) AddRemotePlayer(pseudo string) {
	x := app.Variables.StartingPosX
	y := app.Variables.StartingPosY

	tileX := x / app.Variables.Tileset_size
	tileY := y / app.Variables.Tileset_size

	if _, exists := (*app.Variables.RemotePlayers)[pseudo]; !exists {
		(*app.Variables.RemotePlayers)[pseudo] = &vars.Player{
			Pseudo:    pseudo,
			Position:  &vars.Position{X: x, Y: y},
			Direction: &vars.Direction{X: 0, Y: 1},
			Zoom:      app.Variables.Zoom,
		}
	}
	app.UpdateRemotePlayerPosition(pseudo, x, y, tileX, tileY, 0, 1, 0)
}

func (app *App) RemoveRemotePlayer(pseudo string) {
	if app.Variables.RemotePlayers != nil {
		new_remote_players := make(map[string]*vars.Player)

		for player_pseudo, player := range *app.Variables.RemotePlayers {
			if player_pseudo != pseudo {
				new_remote_players[player_pseudo] = player
			}
		}

		app.Variables.RemotePlayers = &new_remote_players
	}
}

func (app *App) ResetRemotePlayers() {
	remote_players := make(map[string]*vars.Player)
	app.Variables.RemotePlayers = &remote_players
}

func (app *App) UpdateGroup(groupState state.GroupState) {
	app.QueueUpdate(func() {
		app.Variables.PanelsVariables.GroupState = &groupState
		app.UpdateGroupPanel(groupState)

		prev := ""
		if sel := app.Manager.Selects("Group"); len(sel) > 0 {
			prev = sel[0].CurrentOption
		}
		selects := app.GetGroupSelects()
		current := ""
		if len(selects) > 0 {
			for _, opt := range selects[0].Options {
				if opt == prev {
					selects[0].CurrentOption = prev
					break
				}
			}
			current = selects[0].CurrentOption
		}

		app.Manager.SetViewOptions("Group", app.GetGroupOptions(current))
		app.Manager.SetViewSelects("Group", selects)
	})
}

func (app *App) UpdateGroupPanel(grS state.GroupState) {
	gr := app.Variables.PanelsVariables.Group

	// Sync everything
	gr.InGroup = grS.Group != ""
	gr.IsLeader = grS.Leader
	if gr.IsLeader {
		gr.Leader = app.GetPseudo()
	}
	gr.Grouped = grS.Grouped
	gr.UnGrouped = grS.UnGrouped
	gr.Invitations = grS.Invitations
	gr.Promote = grS.Promotion

	if !grS.SendPromotion {
		if gr.SendPromotion != "" {
			gr.Leader = gr.SendPromotion
		}
		gr.SendPromotion = ""
	}

	// Drop a pending promotion
	if gr.SendPromotion != "" && (!slices.Contains(gr.Grouped, gr.SendPromotion) || !gr.IsLeader) {
		gr.SendPromotion = ""
	}

	// Drop invitations
	gr.SendInvitations = slices.DeleteFunc(gr.SendInvitations, func(p string) bool {
		return slices.Contains(gr.Grouped, p) || !slices.Contains(gr.UnGrouped, p)
	})

	if !gr.InGroup {
		gr.Promote = false
		gr.SendPromotion = ""
	}
}
