package engine

import "tap/protocol/parser"

type Player struct {
	ip         string
	name       string
	room       string
	hp         int
	hpMax      int
	group      string
	status     string
	inventory  []string
	quests     []string
	invitation []string
	promotion  bool
	combatId   string
	weapon     *parser.Weapon
}

