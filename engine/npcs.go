package engine

import (
	"slices"
)

type Npc struct {
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Dialogue    []string     `json:"dialogue"`
	Role        string       `json:"role"`
	QuestId     string       `json:"quest_id"`
	Stats       *CombatStats `json:"stats,omitempty"`
	Hostile     bool         `json:"hostile"`
	Damage      int          `json:"damage"`
	InCombat    bool         `json:"omitempty"`
	XpReward    int          `json:"xp_reward,omitempty"`
	ItemsReward []*Item      `json:"items_reward,omitempty"`
	Fighter
}

func (n *Npc) isDead() bool {
	if n.Stats.Hp <= 0 {
		return true
	}
	return false
}

func (n *Npc) takeDamage(amount int) int {
	n.Stats.Hp -= amount
	return amount
	// maybe here substraction with defense of entity
}
func (n *Npc) getHp() int {
	return n.Stats.Hp
}

func (n *Npc) getInitiative() int {
	return n.Stats.Initiative
}

func (n *Npc) getDamage() int {
	return n.Damage
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

func (n *Npc) getName() string {
	return n.Name
}
