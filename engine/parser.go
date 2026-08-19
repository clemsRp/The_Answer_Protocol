package engine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	if err := registerCustomValidations(validate); err != nil {
		panic(fmt.Sprintf("engine: failed to register validators: %v", err))
	}
}

func Get_map(map_path string) (*Map, error) {
	// Get map file content
	file, err := os.ReadFile(map_path)
	if err != nil {
		return nil, errors.New("Invalid file path: Permission denied or File doesn't exist")
	}

	var worlds []Map

	if err := checkDuplicates(file); err != nil {
		return nil, fmt.Errorf("Invalid file: %w", err)
	}

	reader := bytes.NewReader(file)
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&worlds); err != nil {
		return nil, fmt.Errorf("Invalid file: JSON file must be parsable:\n%w", err)
	}

	if len(worlds) == 0 {
		return nil, errors.New("Invalid file: no map found in file")
	}

	world := &worlds[0]

	if err := validate.Struct(world); err != nil {
		return nil, formatValidationError(err)
	}
	if err := validateMapConsistency(world); err != nil {
		return nil, err
	}

	// Add ids to map structures
	for id, item := range world.Items {
		item.Id = id
	}

	for id, room := range world.Rooms {
		room.Id = id
	}

	for id, npc := range world.Npcs {
		npc.Id = id
	}

	return world, nil
}

func validateMapConsistency(m *Map) error {
	for room_id, room := range m.Rooms {
		for exit_dir, exit_room_id := range room.Exits {
			opposite_dir, ok := directions[exit_dir]
			if !ok {
				continue
			}
			target_room, exists := m.Rooms[exit_room_id]
			if !exists {
				continue
			}
			if target_room.Exits[opposite_dir] != room_id {
				return fmt.Errorf(
					"Invalid map: exits aren't consistent between '%s' and '%s'",
					room_id, exit_room_id,
				)
			}
		}
	}
	return nil
}

func formatValidationError(err error) error {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		fe := validationErrs[0]
		return fmt.Errorf(
			"Invalid map: '%s' failed validation '%s' (valeur reçue: '%v')",
			fe.Namespace(),
			fe.Tag(),
			fe.Value(),
		)
	}
	return fmt.Errorf("Invalid map: %w", err)
}

func checkDuplicates(data []byte) error {
	d := json.NewDecoder(bytes.NewReader(data))
	for {
		err := decodeAndCheck(d)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func decodeAndCheck(d *json.Decoder) error {
	t, err := d.Token()
	if err != nil {
		return err
	}

	delim, ok := t.(json.Delim)
	if !ok {
		return nil
	}

	if delim == '{' {
		keys := make(map[string]bool)
		for d.More() {
			tKey, err := d.Token()
			if err != nil {
				return err
			}
			key := tKey.(string)

			if keys[key] {
				return fmt.Errorf("duplicate key: '%s'", key)
			}
			keys[key] = true

			if err := decodeAndCheck(d); err != nil {
				return err
			}
		}
		_, err = d.Token()
		return err
	}

	if delim == '[' {
		for d.More() {
			if err := decodeAndCheck(d); err != nil {
				return err
			}
		}
		_, err = d.Token()
		return err
	}

	return nil
}
