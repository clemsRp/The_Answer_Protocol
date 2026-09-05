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

// registerCustomValidations branche toutes les règles custom sur
// l'instance de validator. C'est le SEUL endroit à toucher pour ajouter
// une nouvelle règle "liste de valeurs autorisées" ou une nouvelle
// vérification d'existence croisée entre les maps de Map
// (rooms / items / npcs / quests).
//
// Pour ajouter une nouvelle structure à valider plus tard, il suffit
// généralement de :
//  1. ajouter des tags `validate:"..."` sur ses champs
//  2. éventuellement ajouter une entrée ici si elle a besoin d'une règle
//     custom (nouvelle liste de valeurs, ou nouvelle existence croisée)
func registerCustomValidations(v *validator.Validate) error {
	rules := map[string]validator.Func{
		// Règles "valeur dans une liste autorisée", basées sur les slices
		// définies ci-dessus. Ajouter une valeur valide = éditer la slice,
		// aucun changement de tag ailleurs dans le code.
		"valid_room_type":    inList(valid_maps),
		"valid_exit":         inList(exits),
		"valid_role":         inList(roles),
		"valid_npc_status":   inList(npc_status),
		"valid_quest_status": inList(quest_status),
		"valid_item_type":    inList(item_types),
		"valid_effect_type":  inList(consumable_type_effects),
		"valid_target_stat":  inList(consumable_target_stats),

		// Règles "existence croisée" : vérifient qu'un id référencé
		// (string) existe bien comme clé dans la sous-map correspondante
		// de la Map de premier niveau (récupérée via fl.Top()).
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

	// Règle de niveau "struct" : le tag `validate` classique ne peut pas
	// croiser un champ parent (Npc.Hostile) avec un champ d'une sous-struct
	// (Npc.Stats.Damage). On enregistre donc une validation dédiée sur Npc
	// qui impose un Damage strictement positif dès que le PNJ est hostile.
	v.RegisterStructValidation(validateHostileNpcDamage, Npc{})

	return nil
}

// inList crée une règle qui vérifie que la valeur du champ (string) fait
// partie de la liste donnée.
func inList(list []string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return is_inside(list, fl.Field().String())
	}
}

// existsIn crée une règle qui vérifie que la valeur du champ (un id,
// string) existe comme clé dans la sous-map de Map retournée par `get`.
// Générique : fonctionne pour n'importe quelle map[string]*T de la Map,
// donc ajouter une nouvelle vérification d'existence (ex: un nouveau type
// de ressource référencée par id) tient en une ligne.
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

// validateHostileNpcDamage garantit que tout PNJ hostile possède des Stats
// avec un Damage strictement positif. Un PNJ non hostile n'a pas besoin de
// dégâts puisqu'il n'entre jamais en combat.
func validateHostileNpcDamage(sl validator.StructLevel) {
	npc := sl.Current().Interface().(Npc)
	if !npc.Hostile {
		return
	}
	if npc.Stats == nil || npc.Stats.Damage <= 0 {
		sl.ReportError(npc.Stats, "Stats.Damage", "Damage", "npc_damage_required", "")
	}
}

// topMap récupère la Map de premier niveau passée à validate.Struct,
// qu'elle ait été passée par valeur ou par pointeur.
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
