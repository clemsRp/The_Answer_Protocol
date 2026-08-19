package engine

import (
	"slices"
)

type Npc struct {
	Name        string       `json:"name" validate:"required"`
	Description string       `json:"description"`
	Dialogue    []string     `json:"dialogue"`
	Role        string       `json:"role" validate:"required,valid_role"`
	QuestId     string       `json:"quest_id" validate:"omitempty,quest_exists"`
	Stats       *CombatStats `json:"stats,omitempty" validate:"required_if=Hostile true"`
	Hostile     bool         `json:"hostile"`
	Damage      int          `json:"damage"`
	InCombat    bool         `json:"omitempty"`
	Fighter     `validate:"-"`
}

func (n *Npc) isDead() bool {
	if n.Stats.Hp <= 0 {
		return true
	}
	return false
}

func (n *Npc) takeDamage(amount int) {
	n.Stats.Hp -= amount
	// maybe here substraction with defense of entity
}
func (n *Npc) getHp() int {
	return n.Stats.Hp
}

func (n *Npc) getInitiative() int {
	return n.Stats.Initiative
}

func (n *Npc) playCombatTurn(target Fighter) *CombatTurnResult {
	target.takeDamage(n.Damage)
	turn_result := &CombatTurnResult{AttackerHp: n.getHp(), TargetHp: target.getHp(), Damage: n.Damage, Status: "combat"}
	return turn_result
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
