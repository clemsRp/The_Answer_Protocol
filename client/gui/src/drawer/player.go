package drawer

import (
	"sort"
	"time"

	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawPlayers() {
	// Get all players
	all_players := make([]*vars.Player, 0)
	all_players = append(all_players, dr.app.Variables.Player)
	for _, remote_player := range *dr.app.Variables.RemotePlayers {
		all_players = append(all_players, remote_player)
	}

	// Sort players by Y position
	sort.Slice(all_players, func(i, j int) bool {
		return all_players[i].Position.Y < all_players[j].Position.Y
	})

	// Draw players
	for _, player := range all_players {
		dr.DrawPlayer(player)
	}
}

func (dr *Drawer) DrawPlayer(p *vars.Player) {
	texture_name, ind_x, ind_y := dr.GetPlayerTexture(p)
	dr.DrawImage(
		texture_name,
		p.Position.X, p.Position.Y,
		float32(ind_x), float32(ind_y), 1, 1, dr.app.Variables.Zoom, 0,
	)

	if p.Pseudo == dr.app.Variables.Player.Pseudo {
		dr.DrawPlayerEmotesPanel()

	} else {
		datas, ok := (*dr.app.Variables.RemotePlayerEmotes)[p.Pseudo]
		if ok && datas.Emote != -1 && time.Since(datas.ChoiceTime) <= 5*time.Second {

			dr.DrawImage(
				vars.PLAYER_EMOTES_BUBBLE_TEXTURE,
				float32(p.Position.X)+dr.app.Variables.Tileset_size,
				float32(p.Position.Y)-0.1*dr.app.Variables.Tileset_size,
				0, 0, 4, 4, dr.app.Variables.Zoom/4, 90,
			)

			dr.DrawImage(
				vars.PLAYER_EMOTES_TEXTURE,
				float32(p.Position.X)+1.07*dr.app.Variables.Tileset_size,
				float32(p.Position.Y)-0.025*dr.app.Variables.Tileset_size,
				float32(datas.Emote+1), 0, 1, 1,
				dr.app.Variables.Zoom*0.8, 0,
			)
		}
	}
}

func (dr *Drawer) DrawPseudos() {
	if dr.app.Variables.Player != nil && dr.app.Variables.Player.Position != nil {
		dr.drawSinglePseudo(dr.app.Variables.Player, rl.Black)
	}

	if (*dr.app.Variables.RemotePlayers) != nil {
		for _, p := range *dr.app.Variables.RemotePlayers {
			if p.Position != nil && p.Position.X >= 0 {
				dr.drawSinglePseudo(p, rl.White)
			}
		}
	}
}

func (dr *Drawer) drawSinglePseudo(p *vars.Player, color rl.Color) {
	if p.Pseudo == "" {
		return
	}
	pseudo := dr.LimitString(p.Pseudo, 12)

	font_size := dr.app.Variables.FontSize
	text_size := rl.MeasureText(pseudo, int32(font_size))
	center_text := (vars.FRAME_WIDTH*int32(dr.app.Variables.Zoom) - text_size) / 2

	rl.DrawText(
		pseudo,
		int32(p.Position.X)+center_text,
		int32(p.Position.Y)-int32(font_size),
		int32(font_size), color,
	)
}

func (dr *Drawer) GetPlayerTexture(p *vars.Player) (string, int, int) {
	dir_x := p.Direction.X
	dir_y := p.Direction.Y

	texture_name := vars.PLAYER_TEXTURE

	var ind_x int
	var ind_y int

	switch (vars.Direction{X: dir_x, Y: dir_y}) {
	case Up:
		ind_x = 0
	case UpLeft:
		ind_x = 1
	case Left:
		ind_x = 2
	case DownLeft:
		ind_x = 3
	case Down:
		ind_x = 4
	case DownRight:
		ind_x = 5
	case Right:
		ind_x = 6
	case UpRight:
		ind_x = 7
	default:
		ind_x = 4
	}

	ind_y = (int(time.Since(dr.app.Variables.StartTime)) % (4 * vars.PLAYER_ANIM_DURATION)) / vars.PLAYER_ANIM_DURATION

	if dir_x == 0 && dir_y == 0 {
		texture_name = vars.PLAYER_REST_TEXTURE
		ind_x = (int(time.Since(dr.app.Variables.StartTime)) % (2 * vars.PLAYER_ANIM_DURATION * 3)) / (vars.PLAYER_ANIM_DURATION * 3)
		ind_x = 3*ind_x + 1
		ind_y = 1
	}

	return texture_name, ind_x, ind_y
}

func (dr *Drawer) DrawPlayerEmotesPanel() {
	emotes := dr.app.Variables.PanelsVariables.Emotes

	if emotes.LastEmoteIndex >= 0 {
		dr.DrawImage(
			vars.PLAYER_EMOTES_BUBBLE_TEXTURE,
			float32(dr.app.Variables.Player.Position.X)+dr.app.Variables.Tileset_size,
			float32(dr.app.Variables.Player.Position.Y)-0.1*dr.app.Variables.Tileset_size,
			0, 0, 4, 4, dr.app.Variables.Zoom/4, 90,
		)

		dr.DrawImage(
			vars.PLAYER_EMOTES_TEXTURE,
			float32(dr.app.Variables.Player.Position.X)+1.07*dr.app.Variables.Tileset_size,
			float32(dr.app.Variables.Player.Position.Y)-0.025*dr.app.Variables.Tileset_size,
			float32(emotes.LastEmoteIndex+1), 0, 1, 1,
			dr.app.Variables.Zoom*0.8, 0,
		)
		return
	}

	dr.DrawImage(
		vars.PLAYER_EMOTES_BUBBLE_TEXTURE,
		float32(dr.app.Variables.Player.Position.X)+0.8*dr.app.Variables.Tileset_size,
		float32(dr.app.Variables.Player.Position.Y)-0.1*dr.app.Variables.Tileset_size,
		9, 6, 2, 2, dr.app.Variables.Zoom/4, 45,
	)

	if emotes.Open {
		dr.DrawEmotePanel()
		dr.DrawPlayerEmotes()
	}
}

func (dr *Drawer) DrawEmotePanel() {
	startX := float32(dr.app.Variables.Player.Position.X) + dr.app.Variables.Tileset_size
	startY := float32(dr.app.Variables.Player.Position.Y) - 0.1*dr.app.Variables.Tileset_size
	scale := dr.app.Variables.Zoom / 4

	segmentWidth := 32.0 * scale

	nbMiddleParts := 6

	dr.DrawImage(
		vars.PLAYER_EMOTES_BUBBLE_TEXTURE,
		startX,
		startY,
		0, 3, 4, 1, scale, 90,
	)

	for i := 0; i < nbMiddleParts; i++ {
		offsetX := float32(i+1) * segmentWidth

		dr.DrawImage(
			vars.PLAYER_EMOTES_BUBBLE_TEXTURE,
			startX+offsetX-0.15*dr.app.Variables.Tileset_size,
			startY-0.12*dr.app.Variables.Tileset_size,
			0, 1, 4, 2, scale, 90,
		)
	}

	dr.DrawImage(
		vars.PLAYER_EMOTES_BUBBLE_TEXTURE,
		startX+float32(nbMiddleParts+1)*segmentWidth-0.3*dr.app.Variables.Tileset_size,
		startY,
		0, 0, 4, 1, scale, 90,
	)
}

func (dr *Drawer) DrawPlayerEmotes() {
	startX := float32(dr.app.Variables.Player.Position.X) + dr.app.Variables.Tileset_size
	startY := float32(dr.app.Variables.Player.Position.Y) - 0.1*dr.app.Variables.Tileset_size

	for i := 0; i < 6; i++ {
		dr.DrawImage(
			vars.PLAYER_EMOTES_TEXTURE,
			startX+(float32(i)/2+0.6)*dr.app.Variables.Tileset_size,
			startY-0.12*dr.app.Variables.Tileset_size,
			float32(i+1), 1, 1, 1,
			dr.app.Variables.Zoom/2, 0,
		)
	}
}
