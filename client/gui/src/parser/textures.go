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

		if !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))

			if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
				baseName := filepath.Base(path)
				textureName := strings.TrimSuffix(baseName, filepath.Ext(baseName))

				textures[textureName] = rl.LoadTexture(path)
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
func DrawImage(texture rl.Texture2D, posX, posY, indX, indY, ratioX, ratioY, zoom, rotation float32) { //[cite: 18]
	sourceW := ratioX * float32(vars.FRAME_WIDTH)
	sourceH := ratioY * float32(vars.FRAME_HEIGHT)

	sourceRec := rl.NewRectangle(
		indX*float32(vars.FRAME_WIDTH), indY*float32(vars.FRAME_HEIGHT),
		sourceW, sourceH,
	)

	destW := float32(vars.FRAME_WIDTH) * zoom
	destH := float32(vars.FRAME_HEIGHT) * zoom

	destRec := rl.NewRectangle(posX+(destW/2), posY+(destH/2), destW, destH)

	origin := rl.NewVector2(destW/2, destH/2)

	rl.DrawTexturePro(texture, sourceRec, destRec, origin, rotation, rl.White)
}
