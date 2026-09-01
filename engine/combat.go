package engine

import pr "tap/protocol"

func (e *Engine) get_combat_stats() (string, any, error) {
	res := pr.CombatStatsCommandData{
		Leader:      "clement",
		CurrentTurn: "bob",
		Team: map[string]pr.CombatPersonData{
			"clement": {
				Name:      "clement",
				Hp:        100,
				Inventory: []string{"sabre", "ton crane humide"},
			},
			"bob": {
				Name:      "bob",
				Hp:        100,
				Inventory: []string{"faaaaaaaaaa", "uuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuu"},
			},
		},
		Opponents: map[string]pr.CombatPersonData{
			"grandpa_gaston": {
				Name:      "grandpa_gaston",
				Hp:        100,
				Inventory: []string{"sabre", "ton crane humide"},
			},
			"granny_jeanine": {
				Name:      "granny_jeanine",
				Hp:        100,
				Inventory: []string{"sabre", "ton crane humide"},
			},
		},
	}

	return "OK", res, nil
}
