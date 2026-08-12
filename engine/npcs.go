package engine

import (
	"slices"
)

type Npc struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Dialogue    []string     `json:"dialogue"`
	Role        string       `json:"role"`
	QuestId     string       `json:"quest_id"`
	Stats       *CombatStats `json:"stats"`
	Hostile     bool         `json:"hostile"`
	InCombat    bool         `json:"omitempty"`
	Fighter
}

func (n *Npc) isDead() bool {
	if n.Stats.Hp <= 0 {
		return true
	}
	return false
}

func (n *Npc) takeDamage(amount int) {
	n.Stats.Hp -= amount
}

func (n *Npc) Clone() *Npc {
	if n == nil {
		return nil
	}
	// copies only primitives (int strings bools)
	cp := *n
	// copies deep all elements like slices, objects CombatStats custom Clone method
	cp.Stats = n.Stats.Clone()
	cp.Dialogue = slices.Clone(n.Dialogue)
	return &cp
}

func isNpcInRoom(room *Room, npcName string) bool {
	for _, npc_name := range room.Npcs {
		if npc_name == npcName {
			return true
		}
	}
	return false
}
