package engine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func Get_map(map_path string) (*Map, error) {
	// Get map file content
	file, err := os.ReadFile(map_path)
	if err != nil {
		return nil, errors.New("Invalid file path: Permission denied or File doesn't exist")
	}
	var worlds []Map

	reader := bytes.NewReader(file)
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	err = dec.Decode(&worlds)
	if err != nil {
		return nil, fmt.Errorf("Invalid file: JSON file must be parsable:\n%s\n", err)
	}
	
	
	// Handle invalid values
	err = IsValidMap(&worlds[0])
	if err != nil {
		return nil, err
	}

	for id, item := range worlds[0].Items {
		item.Id = id
	}

	return &worlds[0], nil
}
