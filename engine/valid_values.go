package engine

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

const (
	RoomEntrance         = "entrance"
	RoomProduceSection   = "produce_section"
	RoomMeatCounter      = "meat_counter"
	RoomFishCounter      = "fish_counter"
	RoomCleaningProducts = "cleaning_products"
	RoomPastries         = "pastries"
	RoomClothingAisle    = "clothing_aisle"
	RoomCheckoutLanes    = "checkout_lanes"
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
		RoomProduceSection,
		RoomMeatCounter,
		RoomFishCounter,
		RoomCleaningProducts,
		RoomPastries,
		RoomClothingAisle,
		RoomCheckoutLanes,
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
		"enemy",
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
		"currency",
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
		"damage",
		"shield",
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

func registerCustomValidations(v *validator.Validate) error {
	rules := map[string]validator.Func{
		"valid_room_type":    inList(valid_maps),
		"valid_exit":         inList(exits),
		"valid_role":         inList(roles),
		"valid_npc_status":   inList(npc_status),
		"valid_quest_status": inList(quest_status),
		"valid_item_type":    inList(item_types),
		"valid_effect_type":  inList(consumable_type_effects),
		"valid_target_stat":  inList(consumable_target_stats),

		"room_exists":  existsIn(func(m *Map) map[string]*Room { return m.Rooms }),
		"item_exists":  existsIn(func(m *Map) map[string]*Item { return m.Items }),
		"npc_exists":   existsIn(func(m *Map) map[string]*Npc { return m.Npcs }),
		"quest_exists": existsIn(func(m *Map) map[string]*Quest { return m.Quests }),
	}

	for tag, fn := range rules {
		if err := v.RegisterValidation(tag, fn); err != nil {
			return fmt.Errorf("registering validator '%s': %w", tag, err)
		}
	}

	v.RegisterStructValidation(validateHostileNpcDamage, Npc{})

	return nil
}

func inList(list []string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return is_inside(list, fl.Field().String())
	}
}

func existsIn[T any](get func(*Map) map[string]*T) validator.Func {
	return func(fl validator.FieldLevel) bool {
		m := topMap(fl)
		if m == nil {
			return false
		}
		_, exists := get(m)[fl.Field().String()]
		return exists
	}
}

func validateHostileNpcDamage(sl validator.StructLevel) {
	npc := sl.Current().Interface().(Npc)
	if !npc.Hostile {
		return
	}
	if npc.Stats == nil || npc.Stats.Damage <= 0 {
		sl.ReportError(npc.Stats, "Stats.Damage", "Damage", "npc_damage_required", "")
	}
}

func topMap(fl validator.FieldLevel) *Map {
	top := fl.Top()
	if top.Kind() == reflect.Ptr {
		if top.IsNil() {
			return nil
		}
		top = top.Elem()
	}
	m, ok := top.Interface().(Map)
	if !ok {
		return nil
	}
	return &m
}
