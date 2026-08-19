package parser

import (
	"encoding/json"
	"errors"
	"os"
)

type TextureData struct {
	TilesetNumbers []int  `json:"tileset_numbers"`
	TilesetNumber  int    `json:"-"`
	Path           string `json:"path"`
	IndX           int    `json:"-"`
	IndY           int    `json:"-"`
}

type ConvertTextures map[int]TextureData

func LoadConvertTextures(filePath string) ([]TextureData, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var convertData []TextureData
	err = json.Unmarshal(fileData, &convertData)
	if err != nil {
		return nil, err
	}

	return convertData, nil
}

func (c *ConvertTextures) GetTilesetDatas(tileset_number int) (string, int, int, error) {
	tileset, ok := (*c)[tileset_number]
	if !ok {
		return "", 0, 0, errors.New("Invalid tileset")
	}

	return tileset.Path, tileset.IndX, tileset.IndY, nil
}
