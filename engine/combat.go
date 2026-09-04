package engine

import (
	"errors"

	pr "tap/protocol"
)

func (e *Engine) get_combat_stats(player *Player) (string, any, error) {
	cs, exists := e.activeCombats[player.stats.CombatId]
	if !exists {
		return "", "", errors.New(pr.ErrInternalServer)
	}

	leaderName := player.name
	if group, ok := e.groups[player.group]; ok {
		leaderName = group.leader.name
	}

	currentTurnName := ""
	if cs.CurrentTurn >= 0 && cs.CurrentTurn < len(cs.Fighters) {
		currentTurnName = cs.Fighters[cs.CurrentTurn].getName()
	}

	team := make(map[string]pr.CombatPersonData)
	for _, p := range cs.Players {
		inventory := make([]string, 0, len(p.inventory))
		for _, item := range p.inventory {
			inventory = append(inventory, item.Id) // IDs stringified for inventory
		}
		team[p.name] = pr.CombatPersonData{
			Name:      p.name,
			Hp:        p.stats.Hp,
			Inventory: inventory,
		}
	}

	opponents := make(map[string]pr.CombatPersonData)
	for _, npc := range cs.Npcs {
		inventory := make([]string, 0, len(npc.ItemsReward))
		for _, item := range npc.ItemsReward {
			inventory = append(inventory, item.Id) // IDs stringified for inventory
		}
		opponents[npc.Id] = pr.CombatPersonData{
			Name:      npc.Name,
			Hp:        npc.Stats.Hp,
			Inventory: inventory,
		}
	}

	res := pr.CombatStatsCommandData{
		Leader:      leaderName,
		CurrentTurn: currentTurnName,
		Team:        team,
		Opponents:   opponents,
	}

	return "OK", res, nil
}
