package engine

import (
	"errors"
	"fmt"
)

const (
	RoomEntrance     = "entrance"
	RoomHealthAisle  = "health_aisle"
	RoomFreshSection = "fresh_section"
)

const (
	South = "south"
	North = "north"
	East  = "east"
	West  = "west"
)

var (
	valid_maps = []string{
		RoomEntrance,
		RoomHealthAisle,
		RoomFreshSection,
	}

	exits = []string{
		North,
		South,
		East,
		West,
	}

	directions = map[string]string{
		North: South,
		South: North,
		West:  East,
		East:  West,
	}

	roles = []string{
		"quest",
		"dialogue",
	}

	npc_status = []string{
		"healthy",
		"dead",
	}

	quest_status = []string{
		"available",
		"progress",
		"unavailable",
	}
	item_types = []string{
		"ressource",
		"consumable",
		"weapon",
	}
	consumable_type_effects = []string{
		"heal",
		"buff",
		"cure",
	}
	consumable_target_stats = []string{
		"hp",
		"mana",
		"max_hp",
		"status",
		"initiative",
	}
)

func is_inside(elements []string, value string) bool {
	for _, element := range elements {
		if element == value {
			return true
		}
	}

	return false
}
func IsValidRoom(room *Room, m *Map, room_id string) error {
	// Check exits
	for exit_dir, exit_room := range room.Exits {
		if !is_inside(exits, exit_dir) {
			return fmt.Errorf("Invalid map: '%s' exits doesn't exist", exit_dir)
		}
		if _, exists := m.Rooms[exit_room]; !exists {
			return fmt.Errorf("Invalid map: '%s' room doesn't exist", exit_room)
		}
	}

	// Check items
	for _, item := range room.Items {
		if _, exists := m.Items[item]; !exists {
			return fmt.Errorf("Invalid map: '%s' item doesn't exist", item)
		}
	}

	// Check npcs
	for _, npc := range room.Npcs {
		if _, exists := m.Npcs[npc]; !exists {
			return fmt.Errorf("Invalid map: '%s' npc doesn't exist", npc)
		}
	}

	for exit_dir, exit_room := range room.Exits {
		if _, ok := directions[exit_dir]; !ok || m.Rooms[exit_room].Exits[directions[exit_dir]] != room_id {
			return errors.New("Invalid map: exits aren't consistents")
		}
	}

	return nil
}

func IsValidNpc(npc *Npc, m *Map) error {
	// Check role
	if !is_inside(roles, npc.Role) {
		return fmt.Errorf("Invalid map: '%s' role isn't valid", npc.Role)
	}

	if npc.QuestId != "" {
		if _, exists := m.Quests[npc.QuestId]; !exists {
			return fmt.Errorf("Invalid map: '%s' quest id doesn't exist", npc.QuestId)
		}
	}
	if npc.Hostile {
		// Check stats
		if npc.Stats.Hp < 0 {
			return fmt.Errorf(
				"Invalid map: hp must be a positive integer != '%d'",
				npc.Stats.Hp,
			)
		} else if npc.Stats.Hp > npc.Stats.HpMax {
			return fmt.Errorf(
				"Invalid map: max_hp must be a positive integer != '%d' and greater or equal than hp (%d)",
				npc.Stats.HpMax,
				npc.Stats.Hp,
			)
		} else if !is_inside(npc_status, npc.Stats.Status) {
			return fmt.Errorf("Invalid map: '%s' status doesn't exist", npc.Stats.Status)
		} else if npc.Stats.Initiative <= 0 {
			return fmt.Errorf("Invalid map: '%s' initiative can't be under or equal to 0", npc.Name)

		}
	}

	return nil
}

func IsValidQuest(quest *Quest) error {
	if !is_inside(quest_status, quest.Status) {
		return fmt.Errorf("Invalid map: '%s' status doesn't exist", quest.Status)
	}

	return nil
}

func IsValidMap(m *Map) error {
	var err error

	// Check each rooms
	for room_id, room := range m.Rooms {
		err = IsValidRoom(room, m, room_id)
		if err != nil {
			return err
		}
	}

	// Check each npcs
	for _, npc := range m.Npcs {
		err = IsValidNpc(npc, m)
		if err != nil {
			return err
		}
	}

	// Check each quests
	for _, quest := range m.Quests {
		err = IsValidQuest(quest)
		if err != nil {
			return err
		}
	}

	for _, item := range m.Items {
		err = isValidItem(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func isValidItem(item *Item) error {
	if !is_inside(item_types, item.Type) {
		return fmt.Errorf("Invalid item: type '%s'  doesn't exist", item.Type)
	}
	if item.Name == "" {
		return fmt.Errorf("Invalid name of item. item: %s", item.Name)
	}
	if item.Worth < 0 {
		return fmt.Errorf("Invalid item's worth is < 0 item: %s", item.Name)
	}
	switch item.Type {
	case "ressource":
	case "consumable":
		if item.EffectType == "" || !is_inside(consumable_type_effects, item.EffectType) {
			return fmt.Errorf("Invalid item effect type %s", item.EffectType)
		}
		if item.TargetStat == "" || !is_inside(consumable_target_stats, item.TargetStat) {
			return fmt.Errorf("Invalid target_stat for consumable item.  item: %s", item.Name)
		}
		if item.Charges <= 0 {
			return fmt.Errorf("Invalid number of charge on consumable item (must be > 0) item: %s", item.Name)
		}
		if item.Duration <= 0 {
			return fmt.Errorf("Invalid number of duration on consumable item (must be > 0) item: %s", item.Name)
		}
		if item.Amount <= 0 {
			return fmt.Errorf("Invalid amount on consumable item (must be > 0) item: %s", item.Name)
		}
	case "weapon":
		if item.Damage <= 0 {
			return fmt.Errorf("Invalid weapon damage on weapon item (must be > 0): %s", item.Name)
		}
	}
	return nil
}
