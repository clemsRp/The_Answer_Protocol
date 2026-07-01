package parser

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
	Description string `json:"description"`
	Reward      string `json:"reward"`
	Status      string `json:"status"`
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
	QuestId     string   `json:"quest_id"`
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

func Get_map(map_path string) (Map, error) {
	// Get map file content
	file, err := os.ReadFile(map_path)
	if err != nil {
		return Map{}, errors.New("Invalid file path: Permission denied or File doesn't exist")
	}

	// Validate map
	err = valid_map(file)
	if err != nil {
		return Map{}, err
	}

	// Convert map from file to Map structure
	var worlds []Map
	err = json.Unmarshal(file, &worlds)
	if err != nil {
		return Map{}, errors.New("Invalid file: JSON file must be parsable")
	}

	return worlds[0], nil
}
