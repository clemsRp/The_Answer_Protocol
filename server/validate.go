package main

import (
	"encoding/json"
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

	directions := []string{
		"north",
		"south",
		"east",
		"west",
	}

	room_names = make([]string)
	item_names = make([]string)
	npc_names  = make([]string)
)

func check_fields_and_types(datas map[string]any, fields []Field) bool {
	for _, field := range fields {
		if _, ok := datas[field.f_name]; !ok {
			return false

		} else if field.f_type != "any" && reflect.TypeOf(datas[field.f_name]) != typeRegistry[field.f_type] {
			fmt.Println(field.f_name)
			fmt.Println(datas[field.f_name])
			fmt.Println(reflect.TypeOf(datas[field.f_name]))
			fmt.Println(typeRegistry[field.f_type])
			return false
		}
	}

	return len(datas) == len(fields)
}

func IsValidFields(world_map map[string]any) bool {
	// Map fields
	map_fields := check_fields_and_types(world_map, MapFieldTypes)
	if !map_fields {
		return false
	}

	major_fields := map[string][]Field{
		"rooms": RoomFieldTypes,
		"items": ItemFieldTypes,
		"npcs": NpcFieldTypes,
		"quests": QuestFieldTypes,
	}

	// Check major fields
	for major_field, field_types := range major_fields {

		// Cast any into map of any
		field, ok := world_map[major_field].(map[string]any)
		if !ok {
			return false
		}

		// Check fields and types for each sub field
		for _, field_datas := range field {
			if fieldMap, ok := field_datas.(map[string]any); ok {
				if !check_fields_and_types(fieldMap, field_types) {
					return false
				}
			}
		}
	}

	return true
}

func IsValidValues(world_map map[string]any) bool {
	return true
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

	} else if !IsValidFields(world_map) {
		return errors.New("Invalid map: missing fields or unknown fields detected")

	} else if !IsValidValues(world_map) {
		return errors.New("Invalid map: invalid or duplicate values")
	}

	return nil
}
