package parser

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

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

	if _, ok := textures["Chest"]; !ok {
		fmt.Println("Chest texture introuvable dans client/gui/assets")
	}

	return &textures
}

func (t *Textures) UnloadTextures() {
	for _, texture := range *t {
		rl.UnloadTexture(texture)
	}
}
