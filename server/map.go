package main

import (
	"encoding/json"
	"errors"
	"os"
)

type Stats struct {
	Hp     int    `json:"hp"`
	Max_hp int    `json:"max_hp"`
	Status string `json:"status"`
}

type Quest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Target      string `json:"target"`
	Reward      string `json:"reward"`
}

type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Obtainable  bool   `json:"obtainable"`
}

type NPC struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Dialogue    []string `json:"dialogue"`
	Role        string   `json:"role"`
	Stats       Stats    `json:"stats"`
}

type Room struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Exits       map[string]string `json:"exits"`
	Items       []string          `json:"items"`
	Npcs        []string          `json:"npcs"`
}

type Map struct {
	Rooms  map[string]*Room  `json:"rooms"`
	Items  map[string]*Item  `json:"items"`
	Npcs   map[string]*NPC   `json:"npcs"`
	Quests map[string]*Quest `json:"quests"`
}

// Validate the map
func (m Map) IsValid() bool {
    for _, room := range m.Rooms {
		// Check rooms
        if room == nil {
            return false
        }

        // Check exits
        for _, targetRoomKey := range room.Exits {
            if _, exists := m.Rooms[targetRoomKey]; !exists {
                return false
            }
        }

        // Check items
        for _, itemKey := range room.Items {
            if _, exists := m.Items[itemKey]; !exists {
                return false
            }
        }

        // Check npcs
        for _, npcKey := range room.Npcs {
            if _, exists := m.Npcs[npcKey]; !exists {
                return false
            }
        }
    }

    return true
}

func get_map(map_path string) (Map, error) {
	// Get map file content
	file, err := os.Open(map_path)
	if err != nil {
		return Map{}, errors.New("Invalid file path: Permission denied or File doesn't exist")
	}

	// Convert the content in json object
	var world []Map

	// Handle additional fields
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&world)
	if err != nil {
		return Map{}, errors.New("Invalid world, doesn't respect world format")
	}

	// Check there is only 1 map
	if len(world) > 1 {
		return Map{}, errors.New("Invalid map, too many maps")
	}

	// Check world exits, references ...
	if !world[0].IsValid() {
		return Map{}, errors.New("Invalid map, check that datas are consistents")
	}

	return world[0], nil
}
