package parser

import (
	"encoding/json"
	"errors"
)

var (
	room_names = make([]string, 0)
	item_names = make([]string, 0)
	npc_names  = make([]string, 0)
)

func valid_map(file []byte) error {
	// Convert content to map[string]any
	var world []map[string]any
	err := json.Unmarshal(file, &world)
	if err != nil {
		return errors.New("Invalid file: JSON file must be parsable")
	}

	world_map := world[0]

	// Valid raw map
	if len(world) != 1 {
		return errors.New("Invalid map: incorrect number of maps")
	}

	// Handle duplicates keys
	err = IsValidKeys(file)
	if err != nil {
		return err
	}

	// Handle invalid fields / types
	err = IsValidFields(world_map)
	if err != nil {
		return err
	}

	// Handle invalid values
	err = IsValidValues(world_map)
	if err != nil {
		return err
	}

	return nil
}
