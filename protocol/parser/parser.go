package parser

import (
	"encoding/json"
	"errors"
	"os"
)

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

	// Convert map from file to pr.Map structure
	var worlds []Map
	err = json.Unmarshal(file, &worlds)
	if err != nil {
		return Map{}, errors.New("Invalid file: JSON file must be parsable")
	}

	// Handle invalid values
	err = worlds[0].IsValidMap()
	if err != nil {
		return Map{}, err
	}

	return worlds[0], nil
}
