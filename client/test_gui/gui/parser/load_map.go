package parser

import (
	"encoding/json"
	"os"
)

type Map struct {
	CompressionLevel int       `json:"compressionlevel"`
	Height           int       `json:"height"`
	Width            int       `json:"width"`
	Infinite         bool      `json:"infinite"`
	TileWidth        int       `json:"tilewidth"`
	TileHeight       int       `json:"tileheight"`
	Orientation      string    `json:"orientation"`
	RenderOrder      string    `json:"renderorder"`
	TiledVersion     string    `json:"tiledversion"`
	Version          string    `json:"version"`
	Type             string    `json:"type"`
	NextLayerID      int       `json:"nextlayerid"`
	NextObjectID     int       `json:"nextobjectid"`
	Layers           []Layer   `json:"layers"`
	Tilesets         []Tileset `json:"tilesets"`
}

type Layer struct {
	Data    []int   `json:"data"`
	Height  int     `json:"height"`
	Width   int     `json:"width"`
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Opacity float32 `json:"opacity"`
	Type    string  `json:"type"`
	Visible bool    `json:"visible"`
	X       int     `json:"x"`
	Y       int     `json:"y"`
}

type Tileset struct {
	FirstGID int    `json:"firstgid"`
	Source   string `json:"source"`
}

func LoadMap(filePath string) (*Map, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var gameMap Map
	err = json.Unmarshal(fileData, &gameMap)
	if err != nil {
		return nil, err
	}

	return &gameMap, nil
}
