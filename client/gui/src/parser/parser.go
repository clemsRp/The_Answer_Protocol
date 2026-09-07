package parser

import "path/filepath"

type Room struct {
	Tilesets [][]int
}

func ParseRooms(map_paths []string) (map[string]*Map, error) {
	res := make(map[string]*Map)

	for _, map_path := range map_paths {
		room_map, err := LoadMap(map_path + ".json")
		if err != nil {
			return nil, err
		}
		room_key := filepath.Base(map_path)

		res[room_key] = room_map
	}

	return res, nil
}
