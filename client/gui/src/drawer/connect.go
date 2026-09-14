package drawer

import (
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawConnectView() {
	if len(dr.app.Manager.Buttons("Connect")) == 0 {
		dr.buildConnectButtons()
	}

	// Draw game view
	rl.BeginTextureMode(dr.gameTexture)
	dr.DrawGameView()
	rl.EndTextureMode()

	// Add blur
	rl.BeginShaderMode(dr.blurShader)
	rl.DrawTextureRec(
		dr.gameTexture.Texture,
		rl.NewRectangle(
			0, 0,
			float32(dr.gameTexture.Texture.Width),
			-float32(dr.gameTexture.Texture.Height),
		),
		rl.NewVector2(0, 0),
		rl.White,
	)
	rl.EndShaderMode()

	connect_start_x := 8.0
	connect_start_y := 4.5
	connect_width := 16.0
	connect_height := 9.0

	// Logo
	logo := (*dr.app.Textures)[vars.LOGO_TEXTURE]
	logo_fraction := 50

	available_height := float32(connect_start_y) * dr.app.Variables.Tileset_size
	target_logo_height := available_height * float32(logo_fraction-2) / float32(logo_fraction)

	logo_zoom := target_logo_height / float32(logo.Height)
	logo_width_px := float32(logo.Width) * logo_zoom

	posX := (float32(dr.app.ScreenWidth) - logo_width_px) / 2
	posY := available_height * 1 / float32(logo_fraction)

	dr.DrawImage(
		vars.LOGO_TEXTURE,
		posX, posY,
		0, 0,
		float32(logo.Width)/float32(vars.FRAME_WIDTH),
		float32(logo.Height)/float32(vars.FRAME_HEIGHT),
		logo_zoom, 0,
	)

	// Connect panel
	rl.DrawRectangle(
		int32(float32(connect_start_x+1)*dr.app.Variables.Tileset_size),
		int32(float32(connect_start_y+1)*dr.app.Variables.Tileset_size),
		int32(float32(connect_width-2)*dr.app.Variables.Tileset_size),
		int32(float32(connect_height-1)*dr.app.Variables.Tileset_size),
		rl.NewColor(220, 185, 138, 255),
	)

	dr.DrawWoodFrameAt(
		vars.Position{X: 8, Y: 4.5},
		vars.Position{X: 24, Y: 13.5},
	)

	// Perso frame
	// TODO

	// Draw buttons
	dr.DrawButtons("Connect")
}
