package main

import (
	"encoding/json"
	"errors"
	"os"
)

type Stats struct {
    Hp     int `json:"hp"`
	Max_hp int `json:"max_hp"`
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
    Rooms  map[string]Room  `json:"rooms"`
    Items  map[string]Item  `json:"items"`
    Npcs   map[string]NPC   `json:"npcs"`
    Quests map[string]Quest `json:"quests"`
}

func get_map(map_path string) (Map, error) {
	// Get map file content
	file, err := os.ReadFile(map_path)
	if err != nil {
		return Map{}, errors.New("Invalid file path: Permission denied or File doesn't exist")
	}

	// Convert the content in json object
	var world []Map
	err = json.Unmarshal(file, &world)
	if err != nil {
		return Map{}, errors.New("Invalid world, doesn't respect world format")
	}

	return world[0], nil
}
