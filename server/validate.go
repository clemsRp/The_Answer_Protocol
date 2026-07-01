package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
)

type Field struct {
	f_name string
	f_type string
}

var typeRegistry = map[string]reflect.Type{
	"string": reflect.TypeOf(""),
	"int":    reflect.TypeOf(0),
	"bool":   reflect.TypeOf(true),
	"map":    reflect.TypeOf(map[string]any{}),
	"list":   reflect.TypeOf([]any{}),
}

var (
	MapFieldTypes = []Field{
		{"rooms", "map"},
		{"items", "map"},
		{"npcs", "map"},
		{"quests", "map"},
	}
	RoomFieldTypes = []Field{
		{"name", "string"},
		{"description", "string"},
		{"exits", "map"},
		{"items", "list"},
		{"npcs", "list"},
	}
	ItemFieldTypes = []Field{
		{"name", "string"},
		{"description", "string"},
		{"obtainable", "bool"},
	}
	NpcFieldTypes = []Field{
		{"name", "string"},
		{"description", "string"},
		{"dialogue", "list"},
		{"role", "string"},
		{"quest_id", "string"},
		{"stats", "map"},
	}
	QuestFieldTypes = []Field{
		{"description", "string"},
		{"reward", "string"},
		{"status", "string"},
	}

	directions = []string{
		"north",
		"south",
		"east",
		"west",
	}

	room_names = make([]string, 0)
	item_names = make([]string, 0)
	npc_names  = make([]string, 0)
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
			return fmt.Errorf("Invalid map: duplicate key found inside '%s'", field)
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

func check_fields_and_types(datas map[string]any, fields []Field) error {
	for _, field := range fields {
		if _, ok := datas[field.f_name]; !ok {
			return errors.New(fmt.Sprintf("Invalid map: missing key '%s'", field.f_name))

		} else if field.f_type != "any" && reflect.TypeOf(datas[field.f_name]) != typeRegistry[field.f_type] {
			return errors.New(
				fmt.Sprintf(
					"Invalid type '%s': expected '%s' got '%s'",
					field.f_name, typeRegistry[field.f_type], reflect.TypeOf(datas[field.f_name]),
				),
			)
		}
	}

	if len(datas) != len(fields) {
		return errors.New("Too many or missing fields")
	}

	return nil
}

func IsValidFields(world_map map[string]any) error {
	// Map fields
	err := check_fields_and_types(world_map, MapFieldTypes)
	if err != nil {
		return err
	}

	major_fields := map[string][]Field{
		"rooms":  RoomFieldTypes,
		"items":  ItemFieldTypes,
		"npcs":   NpcFieldTypes,
		"quests": QuestFieldTypes,
	}

	field_names := map[string][]string{
		"rooms": room_names,
		"items": item_names,
		"npcs":  npc_names,
	}

	// Check major fields
	for major_field, field_types := range major_fields {

		// Cast any into map of any
		field, ok := world_map[major_field].(map[string]any)
		if !ok {
			return errors.New("An error occured while casting major field")
		}

		for field_name, field_datas := range field {
			// Check fields and types for each sub field
			if fieldMap, ok := field_datas.(map[string]any); ok {

				err = check_fields_and_types(fieldMap, field_types)
				if err != nil {
					return err
				}
			}

			// Append name to name list
			field_names[major_field] = append(field_names[major_field], field_name)
		}
	}

	return nil
}

func IsValidValues(world_map map[string]any) error {
	return nil
}

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
