package main

import (
	"os"
	"encoding/json"
	"errors"
)

type Room struct {
    Items      []string          `json:"items"`
    Npcs       []string          `json:"npcs"`
    Neighbours map[string]string `json:"neighbours"`
}

type Item struct {
}

type Npc struct {
}

type Map struct {
    Rooms map[string]Room `json:"rooms"`
    Items map[string]Item `json:"items"`
    Npcs  map[string]Npc  `json:"npcs"`
}

func get_map(map_path string) (Map, error) {
	// Get map file content
	file, err := os.ReadFile(map_path)
	if err != nil {
		return Map{}, errors.Error("Invalid file path: Permission denied or File doesn't exist")
	}

	// Convert the content in json object
	var map []Map
	err = json.Unmarshal(file, &map)
	if err != nil {
		return Map{}, errors.Error("Invalid map, doesn't respect map format")
	}

	return map[0], nil
}
