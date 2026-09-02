package parser

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Textures map[string]rl.Texture2D

func LoadTextures() *Textures {
	textures := make(Textures)

	err := filepath.WalkDir("client/gui/assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Ignore directories, only get files
		if !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))

			// Get supported files
			if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
				// Get texture
				textures[path] = rl.LoadTexture(path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error charging the textures :", err)
	}

	return &textures
}

func (t *Textures) UnloadTextures() {
	for _, texture := range *t {
		rl.UnloadTexture(texture)
	}
}

func DrawImage(texture rl.Texture2D, posX, posY, indX, indY, ratioX, ratioY, zoom float32) {
	sourceRec := rl.NewRectangle(
		indX*vars.FRAME_WIDTH, indY*vars.FRAME_HEIGHT,
		ratioX*vars.FRAME_WIDTH, ratioY*vars.FRAME_HEIGHT,
	)
	destRec := rl.NewRectangle(posX, posY, sourceRec.Width*zoom, sourceRec.Height*zoom)
	origin := rl.NewVector2(0, 0)

	rl.DrawTexturePro(texture, sourceRec, destRec, origin, 0, rl.White)
}
