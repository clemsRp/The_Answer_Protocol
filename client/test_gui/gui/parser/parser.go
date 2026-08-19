package parser

import (
	"test_gui/gui/variables"
)

// Parse Rooms
type Room struct {
	Tilesets [][]int
}

func ParseRooms(room_names []string) (map[string]*Room, error) {
	res := make(map[string]*Room)

	for _, room := range room_names {
		parsed_room, err := ParseRoom(room)
		if err != nil {
			return nil, err
		}
		res[room] = parsed_room
	}

	return res, nil
}

func ParseRoom(room_name string) (*Room, error) {
	// Load map
	room_map, err := LoadMap(room_name + "_map.json")
	if err != nil {
		return nil, err
	}

	// Get tilesets
	tilesets := make([][]int, room_map.Height)
	for y := range room_map.Height {
		temp_tile := make([]int, room_map.Width)
		for x := room_map.Width * y; x < room_map.Width*(y+1); x++ {
			temp_tile[x%room_map.Width] = room_map.Layers[0].Data[x]
		}
		tilesets[y] = temp_tile
	}

	return &Room{
		Tilesets: tilesets,
	}, nil
}

// Parse convert textures
func ParseConvertTextures(filePath string, textures *Textures) (*ConvertTextures, error) {
	// Load map
	convertData, err := LoadConvertTextures(filePath)
	if err != nil {
		return nil, err
	}

	// Get convert textures
	textureMap := make(ConvertTextures)
	for _, item := range convertData {
		for _, tile_ind := range item.TilesetNumbers {
			modified_item := get_modified_item(*textures, item, tile_ind)
			textureMap[tile_ind] = modified_item
		}
	}
	return &textureMap, nil
}

func get_modified_item(textures Textures, item TextureData, tile_ind int) TextureData {
	tile_num := tile_ind - 1

	img_width := textures[item.Path].Width / variables.FRAME_WIDTH

	modified_item := item
	modified_item.TilesetNumber = tile_ind
	modified_item.IndX = tile_num % int(img_width)
	modified_item.IndY = tile_num / int(img_width)

	return modified_item
}
