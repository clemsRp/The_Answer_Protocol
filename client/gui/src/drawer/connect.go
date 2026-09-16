package drawer

import (
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawConnectView() {
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
		2, 0, true,
	)

	// Draw emotes
	dr.DrawWoodFrameAt(
		vars.Position{X: 10, Y: 6},
		vars.Position{X: 13, Y: 9},
		1, 2, false,
	)
	dr.DrawConnectEmotes()
	dr.DrawConnectButtons()

	// Draw input
	posX = 13.5 * dr.app.Variables.Tileset_size
	posY = 6 * dr.app.Variables.Tileset_size
	dialogue := (*dr.app.Textures)[vars.MSG_BUBBLE_TEXTURE]
	dr.DrawImage(
		vars.MSG_BUBBLE_TEXTURE,
		posX, posY,
		0, 0,
		float32(dialogue.Width)/float32(vars.FRAME_WIDTH),
		float32(dialogue.Height)/float32(vars.FRAME_HEIGHT),
		logo_zoom*1.5, 0,
	)

	rl.DrawText(
		"ENTER PSEUDO",
		int32(posX+dr.app.Variables.Tileset_size*0.5),
		int32(posY+dr.app.Variables.Tileset_size*0.3),
		int32(dr.app.Variables.FontSize), dr.app.Colors["panel_text"],
	)

	input_start := int32(posX + dr.app.Variables.Tileset_size)
	input_width := 22.5*dr.app.Variables.Tileset_size - float32(input_start)
	input_height := 2 * dr.app.Variables.FontSize
	border := input_height * 0.1

	// Bordered input
	baseX := int32(posX + dr.app.Variables.Tileset_size)
	baseY := int32(posY + 1.5*dr.app.Variables.Tileset_size)

	dr.DrawBorderedInput(
		baseX,
		baseY,
		int32(input_width),
		int32(input_height),
		int32(border),
		rl.Black, rl.White,
	)

	// Display pseudo
	rl.DrawText(
		dr.app.Variables.Player.Pseudo,
		int32(posX+1.25*dr.app.Variables.Tileset_size),
		int32(posY+1.75*dr.app.Variables.Tileset_size),
		int32(dr.app.Variables.FontSize), rl.Black,
	)

	// Draw cursor
	dr.DrawCursor(
		int32(posX+dr.app.Variables.Tileset_size)+int32(border),
		int32(posY+1.75*dr.app.Variables.Tileset_size),
		int32(border), int32(dr.app.Variables.FontSize),
		dr.app.Variables.Player.Pseudo+"  ", rl.Black,
	)
}
