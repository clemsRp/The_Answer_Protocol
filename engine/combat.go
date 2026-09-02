package engine

import pr "tap/protocol"

func (e *Engine) get_combat_stats() (string, any, error) {
	res := pr.CombatStatsCommandData{
		Leader:      "clement",
		CurrentTurn: "bob",
		Team: map[string]pr.CombatPersonData{
			"clement": {
				Name:      "clement",
				Hp:        75,
				Inventory: []string{"Heuuuuuuuuuuuuuu"},
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
				Hp:        67,
				Inventory: []string{"sabre"},
			},
			/* "granny_jeanine": {
				Name:      "granny_jeanine",
				Hp:        18,
				Inventory: []string{"bonjour", "hello"},
			}, */
		},
	}

	return "OK", res, nil
}
