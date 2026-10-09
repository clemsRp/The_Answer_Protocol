package parser

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Textures map[string]rl.Texture2D

func LoadTextures(assetsFS embed.FS) *Textures {
	textures := make(Textures)

	err := fs.WalkDir(assetsFS, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))

			if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
				baseName := filepath.Base(path)
				textureName := strings.TrimSuffix(baseName, filepath.Ext(baseName))

				data, err := assetsFS.ReadFile(path)
				if err != nil {
					return err
				}

				img := rl.LoadImageFromMemory(ext, data, int32(len(data)))
				textures[textureName] = rl.LoadTextureFromImage(img)
				rl.UnloadImage(img)
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
