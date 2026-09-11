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
	rl.DrawRectangle(0, 0, 750, 250, rl.Red)
}

func (app *App) DrawGameView() {
	app.DrawMap()
	app.DrawPlayer()

	for _, layer := range app.rooms[app.variables.Current_room].Layers {
		if layer.Name == "InFrontOfPlayer" {
			app.DrawLayer(layer, app.rooms[app.variables.Current_room].Tilesets)
		}
	}
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

		posX := float32(gridX * app.variables.Tileset_size)
		posY := float32(gridY * app.variables.Tileset_size)

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
	texture, ind_x, ind_y := app.GetPlayerTexture()
	parser.DrawImage(
		texture,
		app.variables.Player.Position.X, app.variables.Player.Position.Y,
		float32(ind_x), float32(ind_y), 1, 1, app.variables.Zoom, 0,
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

	return (*app.textures)[texture_name], ind_x, ind_y
}
