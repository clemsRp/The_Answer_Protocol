package drawer

import (
	"time"

	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawPlayer() {
	texture_name, ind_x, ind_y := dr.GetPlayerTexture()
	dr.DrawImage(
		texture_name,
		dr.app.Variables.Player.Position.X, dr.app.Variables.Player.Position.Y,
		float32(ind_x), float32(ind_y), 1, 1, dr.app.Variables.Zoom, 0,
	)
}

func (dr *Drawer) DrawPseudo() {
	pseudo := dr.app.GetPseudo()
	font_size := 20

	text_size := rl.MeasureText(pseudo, int32(font_size))
	center_text := (vars.FRAME_WIDTH*int32(dr.app.Variables.Zoom) - text_size) / 2

	rl.DrawText(
		pseudo,
		int32(dr.app.Variables.Player.Position.X)+center_text,
		int32(dr.app.Variables.Player.Position.Y)-int32(font_size),
		int32(font_size), rl.White,
	)
}

func (dr *Drawer) GetPlayerTexture() (string, int, int) {
	dir_x := dr.app.Variables.Player.Direction.X
	dir_y := dr.app.Variables.Player.Direction.Y

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
