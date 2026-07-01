package parser 

import (
	"errors"
	"bytes"
	"encoding/json"
	"io"
	"fmt"
)

func IsValidKeys(file []byte) error {
	// Unmarshal into generic structure
	var world []map[string]map[string]any
	if err := json.Unmarshal(file, &world); err != nil || len(world) == 0 {
		return errors.New("Invalid file: JSON file must be parsable")
	}

	// Initialize the decoder
	dec := json.NewDecoder(bytes.NewReader(file))

	countInJson := make(map[string]int)
	var currentMajorField string
	inMajorFieldsObject := false

	for {
		// Read next token
		t, err := dec.Token()
		if err == io.EOF {
			break
		}

		if err != nil {
			return errors.New("Invalid file: JSON file must be parsable")
		}

		// Handle delimiters
		if delim, ok := t.(json.Delim); ok {
			if delim == '{' {
				// Enter major field object
				if currentMajorField != "" && !inMajorFieldsObject {
					inMajorFieldsObject = true
				}

			} else if delim == '}' {

				// Exit major field object
				if inMajorFieldsObject {
					inMajorFieldsObject = false
					currentMajorField = ""
				}
			}

			// Handle major fields and keys
		} else if str, ok := t.(string); ok {
			if !inMajorFieldsObject && (str == "rooms" || str == "items" || str == "npcs" || str == "quests") {
				// Set current major field
				currentMajorField = str

			} else if inMajorFieldsObject {
				// Increment key counter
				countInJson[currentMajorField]++

				// Skip the value to avoid sub keys
				if err := skipValue(dec); err != nil {
					return err
				}
			}
		}
	}

	// Compare counts to detect duplicate keys
	world_map := world[0]
	for _, field := range []string{"rooms", "items", "npcs", "quests"} {
		if len(world_map[field]) != countInJson[field] {
			return errors.New(
				fmt.Sprintf(
					"Invalid map: duplicate key found inside '%s'",
					field,
				),
			)
		}
	}

	return nil
}

func skipValue(dec *json.Decoder) error {
	// Read next token
	t, err := dec.Token()
	if err != nil {
		return err
	}

	// Check if token is a single value
	_, ok := t.(json.Delim)
	if !ok {
		return nil
	}

	// Skip object or list until matching closure
	depth := 1
	for depth > 0 {
		t, err = dec.Token()
		if err != nil {
			return err
		}

		if d, ok := t.(json.Delim); ok {
			if d == '{' || d == '[' {
				depth++

			} else if d == '}' || d == ']' {
				depth--
			}
		}
	}
	return nil
}
