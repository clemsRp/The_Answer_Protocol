package parser

import (
	"errors"
	"fmt"
	pr "tap/protocol"
)

var (
	exits = []string{
		pr.North,
		pr.South,
		pr.East,
		pr.West,
	}

	directions = map[string]string{
		pr.North: pr.South,
		pr.South: pr.North,
		pr.West:  pr.East,
		pr.East:  pr.West,
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
)

func is_inside(elements []string, value string) bool {
	for _, element := range elements {
		if element == value {
			return true
		}
	}

	return false
}

func (room Room) IsValidRoom(m Map, room_id string) error {
	// Check exits
	for exit_dir, exit_room := range room.Exits {
		if !is_inside(exits, exit_dir) {
			return fmt.Errorf("Invalid map: '%s' exits doesn't exist", exit_dir)

		} else if !is_inside(room_names, exit_room) {
			return fmt.Errorf("Invalid map: '%s' room doesn't exist", exit_room)
		}
	}

	// Check items
	for _, item := range room.Items {
		if !is_inside(item_names, item) {
			return fmt.Errorf("Invalid map: '%s' item doesn't exist", item)
		}
	}

	// Check npcs
	for _, npc := range room.Npcs {
		if !is_inside(npc_names, npc) {
			return fmt.Errorf("Invalid map: '%s' npc doesn't exist", npc)
		}
	}

	// Check directions
	for exit_dir, exit_room := range room.Exits {
		// Check opposite direction is present
		if _, ok := directions[exit_dir]; !ok || m.Rooms[exit_room].Exits[directions[exit_dir]] != room_id {
			return errors.New("Invalid map: exits aren't consistents")
		}
	}

	return nil
}

func (npc Npc) IsValidNpc() error {
	// Check role
	if !is_inside(roles, npc.Role) {
		return fmt.Errorf("Invalid map: '%s' role isn't valid", npc.Role)

		// Check quest id
	} else if npc.QuestId != "" && !is_inside(quest_names, npc.QuestId) {
		return fmt.Errorf("Invalid map: '%s' quest id doesn't exist", npc.QuestId)
	}

	// Check stats
	if npc.Stats.Hp < 0 {
		return fmt.Errorf(
			"Invalid map: hp must be a positive integer != '%d'",
			npc.Stats.Hp,
		)

	} else if npc.Stats.Hp > npc.Stats.MaxHp {
		return fmt.Errorf(
			"Invalid map: max_hp must be a positive integer != '%d' and greater or equal than hp (%d)",
			npc.Stats.MaxHp,
			npc.Stats.Hp,
		)

	} else if !is_inside(npc_status, npc.Stats.Status) {
		return fmt.Errorf("Invalid map: '%s' status doesn't exist", npc.Stats.Status)
	}

	return nil
}

func (quest Quest) IsValidQuest() error {
	if !is_inside(quest_status, quest.Status) {
		return fmt.Errorf("Invalid map: '%s' status doesn't exist", quest.Status)
	}

	return nil
}

func (m Map) IsValidMap() error {
	var err error

	// Check each rooms
	for room_id, room := range m.Rooms {
		err = room.IsValidRoom(m, room_id)
		if err != nil {
			return err
		}
	}

	// Check each npcs
	for _, npc := range m.Npcs {
		err = npc.IsValidNpc()
		if err != nil {
			return err
		}
	}

	// Check each quests
	for _, quest := range m.Quests {
		err = quest.IsValidQuest()
		if err != nil {
			return err
		}
	}

	return nil
}
