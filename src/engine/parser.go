package engine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	pr "tap/src/protocol"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	if err := registerCustomValidations(validate); err != nil {
		panic(fmt.Sprintf("engine: failed to register validators: %v", err))
	}
}

func Get_map(file []byte) (*Map, error) {
	// Get map file content

	var worlds []Map

	if err := checkDuplicates(file); err != nil {
		return nil, fmt.Errorf("invalid file: %w", err)
	}

	reader := bytes.NewReader(file)
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&worlds); err != nil {
		return nil, fmt.Errorf("invalid file: JSON file must be parsable:\n%w", err)
	}

	if len(worlds) == 0 {
		return nil, errors.New("invalid file: no map found in file")
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

	for id, quest := range world.Quests {
		quest.Id = id
	}

	// Check default item exist
	_, exists := world.Items[pr.DefaultItem]
	if !exists {
		return nil, errors.New(pr.ErrDefaultItemNotFound)
	}

	return world, nil
}

func validateMapConsistency(m *Map) error {
	if _, exists := m.Rooms[RoomPlaceDuVillage]; !exists {
		return fmt.Errorf("invalid map: start room '%s' does not exist", RoomPlaceDuVillage)
	}

	for questId, q := range m.Quests {
		if q.TargetItem != "" {
			if _, exists := m.Items[q.TargetItem]; !exists {
				return fmt.Errorf("invalid map: quest '%s' targets non-existent item '%s'", questId, q.TargetItem)
			}
		}
		if q.TargetNpc != "" {
			if _, exists := m.Npcs[q.TargetNpc]; !exists {
				return fmt.Errorf("invalid map: quest '%s' targets non-existent npc '%s'", questId, q.TargetNpc)
			}
		}
		if q.Reward != "" {
			if _, exists := m.Items[q.Reward]; !exists {
				return fmt.Errorf("invalid map: quest '%s' rewards non-existent item '%s'", questId, q.Reward)
			}
		}
	}

	for npcId, n := range m.Npcs {
		for _, loot := range n.ItemsReward {
			if _, exists := m.Items[loot]; !exists {
				return fmt.Errorf("invalid map: npc '%s' drops non-existent item '%s'", npcId, loot)
			}
		}
	}

	for room_id, room := range m.Rooms {
		for exit_dir, exit_room_id := range room.Exits {
			opposite_dir, ok := Directions[exit_dir]
			if !ok {
				continue
			}
			target_room, exists := m.Rooms[exit_room_id]
			if !exists {
				continue
			}
			if target_room.Exits[opposite_dir] != room_id {
				return fmt.Errorf(
					"invalid map: exits aren't consistent between '%s' and '%s'",
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
			"invalid map: '%s' failed validation '%s' (valeur reçue: '%v')",
			fe.Namespace(),
			fe.Tag(),
			fe.Value(),
		)
	}
	return fmt.Errorf("invalid map: %s", err.Error())
}

func checkDuplicates(data []byte) error {
	d := json.NewDecoder(bytes.NewReader(data))
	for {
		err := decodeAndCheck(d)
		if err != nil {
			if err == io.EOF {
				break
			}
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
