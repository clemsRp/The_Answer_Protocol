package engine

import (
	"errors"
	pr "tap/protocol"
	"tap/protocol/parser"

	"github.com/google/uuid"
)

// =========================================
//         COMBAT SYSTEM DESIGN
// =========================================
//
//
// SYSTEM:
// - when ATTACK occurs for the 1st time, state of Player/s (if grouped) enters in Combat (leader or not).
// - a POPUP displays:
//     * displaying ascii art for the player/s.
//     * a panel for the actions and the turn of the person.
// - you have 20 seconds to play your turn.
//
// INITIATIVE:
// - the 1st who start dependS ON INITIATIVE formula.
//
// IN COMBAT PER TURN POSSIBILITES:
// - Can use an item (but can't attack) then passes your turn.
// - Can change weapon (illimited) and then attack the enemy.
// - Can use a SPELL (if you have the mana) then passes your turn.
// - You can try to escape, then you only escape even if you are in a group, then it sends EVENT on the others user in combat.
// - You can still CHAT (only command authorized during combat)
//
//
//
// ENDING:
// - if noone is in combat anymore (everyone fleed) or one PART died, the combat ends.

type CombatSession struct {
	Attackers   map[string]*Player
	Enemies     map[string]*parser.Npc
	CurrentTurn string
	Win         bool
	Cancelled   bool
}

type CombatTurnResult struct {
	AttackerHp int    `json:"attacker_hp"`
	TargetHp   int    `json:"target_hp"`
	Damage     int    `json:"damage"`
	Status     string `json:"status"`
}

func (e *Engine) initiateCombat(player *Player, npcName string) (string, *CombatSession, error) {
	current_room := e.world.Rooms[player.room]

	npc_base, exists := e.world.Npcs[npcName]
	if !exists || !isNpcInRoom(*current_room, npcName) {
		return "", nil, errors.New(pr.ErrNpcNotFound)
	}
	if !npc_base.Hostile {
		return "", nil, errors.New(pr.ErrNpcNotHostile)
	}

	combatID := uuid.New().String()
	copy_npc := *npc_base
	copy_npc.Stats.Status = "combat"

	attackers := make(map[string]*Player)
	enemies := make(map[string]*parser.Npc)
	enemies[copy_npc.Name] = &copy_npc

	current_group, is_in_group := e.groups[player.group]
	if !is_in_group {
		player.status = "combat"
		player.combatId = combatID
		attackers[player.name] = player
	} else {
		for _, p := range current_group {
			p.status = "combat"
			p.combatId = combatID
			attackers[p.name] = p
		}
	}

	combat_session := &CombatSession{
		Attackers: attackers,
		Enemies:   enemies,
	}
	e.activeCombats[combatID] = combat_session
	combat_result := CombatTurnResult{AttackerHp: player.hp, TargetHp: npc_base.Stats.Hp, Damage: }
	return "", combat_session, nil
}

func (e *Engine) processCombatTurn(player *Player, targetName string) (string, any, error) {
	combatSession, exists := e.activeCombats[player.combatId]
	if !exists {
		player.status = "idle"
		return "", "", errors.New(pr.ErrInternalServer)
	}
	target, targetExists := combatSession.Enemies[targetName]
	if !targetExists {
		return "", "", errors.New(pr.ErrNpcNotFound)
	}

	return "" + targetName, nil, nil
}
