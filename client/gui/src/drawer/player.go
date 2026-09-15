package drawer

import (
	"time"

	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawPlayers() {
	if dr.app.Variables.Player != nil && dr.app.Variables.Player.Position != nil {
		dr.DrawPlayer(dr.app.Variables.Player)
	}

	if dr.app.Variables.RemotePlayers != nil {
		for _, p := range dr.app.Variables.RemotePlayers {
			if p.Position != nil && p.Position.X >= 0 {
				dr.DrawPlayer(p)
			}
		}
	}
}

func (dr *Drawer) DrawPlayer(p *vars.Player) {
	texture_name, ind_x, ind_y := dr.GetPlayerTexture(p)
	dr.DrawImage(
		texture_name,
		p.Position.X, p.Position.Y,
		float32(ind_x), float32(ind_y), 1, 1, dr.app.Variables.Zoom, 0,
	)
}

func (dr *Drawer) DrawPseudos() {
	if dr.app.Variables.Player != nil && dr.app.Variables.Player.Position != nil {
		dr.drawSinglePseudo(dr.app.Variables.Player)
	}

	if dr.app.Variables.RemotePlayers != nil {
		for _, p := range dr.app.Variables.RemotePlayers {
			if p.Position != nil && p.Position.X >= 0 {
				dr.drawSinglePseudo(p)
			}
		}
	}
}

func (dr *Drawer) drawSinglePseudo(p *vars.Player) {
	if p.Pseudo == "" {
		return
	}
	font_size := dr.app.Variables.FontSize
	text_size := rl.MeasureText(p.Pseudo, int32(font_size))
	center_text := (vars.FRAME_WIDTH*int32(dr.app.Variables.Zoom) - text_size) / 2

	rl.DrawText(
		p.Pseudo,
		int32(p.Position.X)+center_text,
		int32(p.Position.Y)-int32(font_size),
		int32(font_size), rl.White,
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
