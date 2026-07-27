package parser

import (
	"errors"
	"fmt"
	"reflect"
)

type Field struct {
	f_name string
	f_type string
}

var typeRegistry = map[string]reflect.Type{
	"string": reflect.TypeOf(""),
	"int":    reflect.TypeOf(0.0),
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
	StatsFieldTypes = []Field{
        {"hp", "int"},
        {"max_hp", "int"},
        {"status", "string"},
    }
)

func check_fields_and_types(datas map[string]any, fields []Field) error {
	for _, field := range fields {
		// Handle missing key
		if _, ok := datas[field.f_name]; !ok {
			return errors.New(fmt.Sprintf("Invalid map: missing key '%s'", field.f_name))

		// Handle any type
		} else if field.f_type == "any" {
			continue

		// Handle empty strings
		} else if field.f_type == "string" && field.f_name != "quest_id" && datas[field.f_name] == "" {
			return errors.New(fmt.Sprintf("Invalid map: '%s' field shouldn't be an empty string", field.f_name))

		// Handle other types
		} else if reflect.TypeOf(datas[field.f_name]) != typeRegistry[field.f_type] {
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

func check_exits(exitsData map[string]any, field_name string) error {
    for exitDir, exitDest := range exitsData {
        if reflect.TypeOf(exitDest) != typeRegistry["string"] {
            return errors.New(fmt.Sprintf("Invalid type for exit '%s' in room '%s': expected 'string'", exitDir, field_name))
        }
    }

	return nil
}

func check_stats(statsData map[string]any) error {
	err := check_fields_and_types(statsData, StatsFieldTypes)
    if err != nil {
        return err
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

	field_names := map[string]*[]string{
		"rooms":  &room_names,
		"items":  &item_names,
		"npcs":   &npc_names,
		"quests": &quest_names,
	}

	// Check major fields
	for major_field, field_types := range major_fields {

		// Cast any into map of any
		field, ok := world_map[major_field].(map[string]any)
		if !ok {
			return errors.New("An error occured while casting major field")
		}

		for field_name, field_datas := range field {
			// Handle empty keys
			if field_name == "" {
				return errors.New("Invalid map: key shouldn't be an empty string")
			}

			// Check fields and types for each sub field
			if fieldMap, ok := field_datas.(map[string]any); ok {

				// Check field types
				err = check_fields_and_types(fieldMap, field_types)
				if err != nil {
					return err
				}

				// Check exits
				if major_field == "rooms" {
					if exitsData, ok := fieldMap["exits"].(map[string]any); ok {
						err = check_exits(exitsData, field_name)
						if err != nil {
							return err
						}
					}
				}
				
				// Check stats
				if major_field == "npcs" {
					if statsData, ok := fieldMap["stats"].(map[string]any); ok {
						err = check_stats(statsData)
						if err != nil {
							return err
						}
					}
				}
			}

			// Append name to name list
			*field_names[major_field] = append(*field_names[major_field], field_name)
		}
	}

	return nil
}
