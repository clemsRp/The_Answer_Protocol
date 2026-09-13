package gui

import (
	"path/filepath"
	"strings"
	"tap/client/gui/src/parser"
	"time"

	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)



var (
	Down  = vars.Direction{X: 0, Y: 1}
	Up    = vars.Direction{X: 0, Y: -1}
	Left  = vars.Direction{X: -1, Y: 0}
	Right = vars.Direction{X: 1, Y: 0}

	DownLeft  = vars.Direction{X: -1, Y: 1}
	DownRight = vars.Direction{X: 1, Y: 1}
	UpLeft    = vars.Direction{X: -1, Y: -1}
	UpRight   = vars.Direction{X: 1, Y: -1}
)

func (app *App) DrawConnectView() {
	// Draw game view
	rl.BeginTextureMode(app.gameTexture)
	app.DrawGameView()
	rl.EndTextureMode()

	// Add blur
	rl.BeginShaderMode(app.blurShader)
	rl.DrawTextureRec(
		app.gameTexture.Texture,
		rl.NewRectangle(
			0, 0,
			float32(app.gameTexture.Texture.Width),
			-float32(app.gameTexture.Texture.Height),
		),
		rl.NewVector2(0, 0),
		rl.White,
	)
	rl.EndShaderMode()

	connect_start_x := 8.0
	connect_start_y := 4.5
	connect_width := 16.0
	connect_height := 9.0

	// Draw connect view
	// Logo
	logo := (*app.textures)[vars.LOGO_TEXTURE]
	logo_fraction := 50

	available_height := float32(connect_start_y) * app.variables.Tileset_size
	target_logo_height := available_height * float32(logo_fraction-2) / float32(logo_fraction)

	logo_zoom := target_logo_height / float32(logo.Height)
	logo_width_px := float32(logo.Width) * logo_zoom

	posX := (float32(app.screenWidth) - logo_width_px) / 2
	posY := available_height * 1 / float32(logo_fraction)

	parser.DrawImage(
		logo,
		posX, posY,
		0, 0,
		float32(logo.Width)/float32(vars.FRAME_WIDTH),
		float32(logo.Height)/float32(vars.FRAME_HEIGHT),
		logo_zoom, 0,
	)

	// Connect panel
	rl.DrawRectangle(
		int32(float32(connect_start_x+1)*app.variables.Tileset_size),
		int32(float32(connect_start_y+1)*app.variables.Tileset_size),
		int32(float32(connect_width-2)*app.variables.Tileset_size),
		int32(float32(connect_height-1)*app.variables.Tileset_size),
		rl.NewColor(220, 185, 138, 255),
	)

	app.DrawWoodFrameAt(
		vars.Position{X: 8, Y: 4.5},
		vars.Position{X: 24, Y: 13.5},
	)

	// Play button
	parser.DrawImage(
		(*app.textures)[vars.PLAY_TEXTURE],
		0, 0, 0, 2, 6, 2, app.variables.Zoom, 0,
	)
}

func (app *App) DrawGameView() {
	app.DrawMap()
	app.DrawPlayer()

	for _, layer := range app.rooms[app.variables.Current_room].Layers {
		if layer.Name == "InFrontOfPlayer" {
			app.DrawLayer(layer, app.rooms[app.variables.Current_room].Tilesets)
		}
	}

	app.DrawPseudo()
}

func (app *App) DrawMap() {
	cur_room := app.rooms[app.variables.Current_room]
	if cur_room == nil {
		return
	}

	for _, layer := range cur_room.Layers {
		if !layer.Visible || layer.Name == "InFrontOfPlayer" {
			continue
		}
		app.DrawLayer(layer, cur_room.Tilesets)
	}

	app.DrawWoodFrameAt(
		vars.Position{X: 0, Y: 0},
		vars.Position{X: 32, Y: 18},
	)
}

func (app *App) DrawLayer(layer parser.Layer, tilesets []parser.Tileset) {
	const (
		FLIPPED_HORIZONTALLY_FLAG = 0x80000000
		FLIPPED_VERTICALLY_FLAG   = 0x40000000
		FLIPPED_DIAGONALLY_FLAG   = 0x20000000
	)

	for i, rawTile := range layer.Data {
		if rawTile == 0 {
			continue
		}

		flipH := (rawTile & FLIPPED_HORIZONTALLY_FLAG) != 0
		flipV := (rawTile & FLIPPED_VERTICALLY_FLAG) != 0
		flipD := (rawTile & FLIPPED_DIAGONALLY_FLAG) != 0

		tile := rawTile & 0x0FFFFFFF
		if tile == 0 {
			continue
		}

		gridX := i % layer.Width
		gridY := i / layer.Width

		posX := float32(gridX * int(app.variables.Tileset_size))
		posY := float32(gridY * int(app.variables.Tileset_size))

		var activeTileset parser.Tileset
		for j := len(tilesets) - 1; j >= 0; j-- {
			if tile >= tilesets[j].FirstGID {
				activeTileset = tilesets[j]
				break
			}
		}

		if activeTileset.FirstGID == 0 {
			continue
		}

		sourcePath := activeTileset.Source
		baseName := filepath.Base(sourcePath)
		textureName := strings.TrimSuffix(baseName, filepath.Ext(baseName))

		texture, ok := (*app.textures)[textureName]
		if !ok {
			continue
		}

		localID := tile - activeTileset.FirstGID
		localID = app.GetLocalID(localID, textureName)

		cols := int(texture.Width) / vars.FRAME_WIDTH
		if cols == 0 {
			cols = 1
		}

		indX := float32(localID % cols)
		indY := float32(localID / cols)

		var rotation float32 = 0
		var rX, rY float32 = 1, 1

		if flipD {
			rotation = 90
			rX = 1
			rY = -1
			if flipH {
				rX = -rX
			}
			if flipV {
				rY = -rY
			}
		} else {
			if flipH {
				rX = -1
			}
			if flipV {
				rY = -1
			}
		}

		parser.DrawImage(texture, posX, posY, indX, indY, rX, rY, app.variables.Zoom, rotation)
	}
}

func (app *App) DrawPlayer() {
	// Draw player texture
	texture, ind_x, ind_y := app.GetPlayerTexture()
	parser.DrawImage(
		texture,
		app.variables.Player.Position.X, app.variables.Player.Position.Y,
		float32(ind_x), float32(ind_y), 1, 1, app.variables.Zoom, 0,
	)
}

func (app *App) DrawPseudo() {
	pseudo := app.GetPseudo()
	font_size := 20

	text_size := rl.MeasureText(pseudo, int32(font_size))
	center_text := (vars.FRAME_WIDTH*int32(app.variables.Zoom) - text_size) / 2

	rl.DrawText(
		pseudo,
		int32(app.variables.Player.Position.X)+center_text,
		int32(app.variables.Player.Position.Y)-int32(font_size),
		int32(font_size), rl.White,
	)
}

func (app *App) GetPlayerTexture() (rl.Texture2D, int, int) {
	// Get directions
	dir_x := app.variables.Player.Direction.X
	dir_y := app.variables.Player.Direction.Y

	// Get texture_name
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

	ind_y = (int(time.Since(app.variables.StartTime)) % (4 * vars.PLAYER_ANIM_DURATION)) / vars.PLAYER_ANIM_DURATION

	if dir_x == 0 && dir_y == 0 {
		texture_name = vars.PLAYER_REST_TEXTURE
		ind_x = (int(time.Since(app.variables.StartTime)) % (2 * vars.PLAYER_ANIM_DURATION * 3)) / (vars.PLAYER_ANIM_DURATION * 3)
		ind_x = 3*ind_x + 1
		ind_y = 1
	}

	return (*app.textures)[texture_name], ind_x, ind_y
}

func (app *App) DrawWoodFrameAt(visualStart vars.Position, visualEnd vars.Position) {
	start := vars.Position{
		X: visualStart.X / 2,
		Y: visualStart.Y / 2,
	}
	end := vars.Position{
		X: visualEnd.X - start.X + 1,
		Y: visualEnd.Y - start.Y + 1,
	}
	app.DrawWoodFrame(start, end)
}

func (app *App) DrawWoodFrame(start vars.Position, end vars.Position) {
	frame_start_x := 12
	frame_start_y := 0

	frame_width := int((end.X-start.X)/2 + 1)
	frame_height := int((end.Y-start.Y)/2 + 1)

	for x := range frame_width {
		for y := range frame_height {
			// Skip center
			if x != 0 && y != 0 && x != frame_width-1 && y != frame_height-1 {
				continue
			}

			// Get indexs
			var ind_x int
			var ind_y int
			switch x {
			case 0:
				ind_x = 0
			case frame_width - 1:
				ind_x = 2
			default:
				ind_x = 1
			}
			switch y {
			case 0:
				ind_y = 0
			case frame_height - 1:
				ind_y = 2
			default:
				ind_y = 1
			}

			// Draw frame
			parser.DrawImage(
				(*app.textures)[vars.WOOD_FRAME_TEXTURE],
				2*app.variables.Tileset_size*(float32(x)+start.X)-app.variables.Tileset_size,
				2*app.variables.Tileset_size*(float32(y)+start.Y)-app.variables.Tileset_size,
				float32(frame_start_x+ind_x), float32(frame_start_y+ind_y),
				1, 1, 2*app.variables.Zoom, 0,
			)
		}
	}
}
