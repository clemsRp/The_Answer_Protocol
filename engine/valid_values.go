package engine

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

const (
	RoomFermeCharbonneau = "ferme_charbonneau"
	RoomFermeUzoloise    = "ferme_uzoloise"
	RoomFontaneilles     = "fontaneilles"
	RoomVergerDeLaSavane = "verger_de_la_savane"
	RoomFromagerie       = "fromagerie"
	RoomPlaceDuVillage   = "place_du_village"
	RoomForetMalicieuse  = "foret_malicieuse"
	RoomCamping          = "camping"
	RoomRiviere          = "riviere"
	RoomBarrage          = "barrage"
	RoomLacDuBarrage     = "lac_du_barrage"
	RoomJardin           = "jardin"

	South = "south"
	North = "north"
	East  = "east"
	West  = "west"
)

var (
	ValidMaps = []string{
		RoomFermeCharbonneau,
		RoomFermeUzoloise,
		RoomFontaneilles,
		RoomVergerDeLaSavane,
		RoomFromagerie,
		RoomPlaceDuVillage,
		RoomForetMalicieuse,
		RoomCamping,
		RoomRiviere,
		RoomBarrage,
		RoomLacDuBarrage,
		RoomJardin,
	}

	Exits = []string{
		North,
		South,
		East,
		West,
	}

	Directions = map[string]string{
		North: South,
		South: North,
		West:  East,
		East:  West,
	}

	Roles = []string{
		"quest",
		"dialogue",
		"enemy",
	}

	NpcStatus = []string{
		"healthy",
		"dead",
	}

	QuestStatus = []string{
		"available",
		"progress",
		"unavailable",
	}
	ItemTypes = []string{
		"ressource",
		"consumable",
		"weapon",
		"currency",
	}
	ConsumableTypeEffects = []string{
		"heal",
		"buff",
		"cure",
	}
	ConsumableTargetStats = []string{
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
		"valid_room_type":    inList(ValidMaps),
		"valid_exit":         inList(Exits),
		"valid_role":         inList(Roles),
		"valid_npc_status":   inList(NpcStatus),
		"valid_quest_status": inList(QuestStatus),
		"valid_item_type":    inList(ItemTypes),
		"valid_effect_type":  inList(ConsumableTypeEffects),
		"valid_target_stat":  inList(ConsumableTargetStats),

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
