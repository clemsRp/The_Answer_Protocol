package drawer

import (
	"path/filepath"
	"strings"
	"tap/client/gui/src/parser"

	vars "tap/client/gui/src/variables"
)

func (dr *Drawer) DrawMap() {
	cur_room := dr.app.Rooms[dr.app.Variables.Current_room]
	if cur_room == nil {
		return
	}

	for _, layer := range cur_room.Layers {
		if !layer.Visible || layer.Name == "InFrontOfPlayer" {
			continue
		}
		dr.DrawLayer(layer, cur_room.Tilesets)
	}

	dr.DrawWoodFrameAt(
		vars.Position{X: 0, Y: 0},
		vars.Position{X: 32, Y: 18},
	)
}

func (dr *Drawer) DrawLayer(layer parser.Layer, tilesets []parser.Tileset) {
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

		posX := float32(gridX * int(dr.app.Variables.Tileset_size))
		posY := float32(gridY * int(dr.app.Variables.Tileset_size))

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

		texture, ok := (*dr.app.Textures)[textureName]
		if !ok {
			continue
		}

		localID := tile - activeTileset.FirstGID
		localID = dr.GetLocalID(localID, textureName)

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

		dr.DrawImage(textureName, posX, posY, indX, indY, rX, rY, dr.app.Variables.Zoom, rotation)
	}
}
