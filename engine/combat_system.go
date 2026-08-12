package engine

import (
	"errors"
	pr "tap/protocol"

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
	Id          string
	Fighters    []Fighter
	CurrentTurn int
	Win         bool
	Cancelled   bool
}

type CombatTurnResult struct {
	AttackerHp int    `json:"attacker_hp"`
	TargetHp   int    `json:"target_hp"`
	Damage     int    `json:"damage"`
	Status     string `json:"status"`
}

func (e *Engine) getValidTarget(room *Room, npcName string) (*Npc, error) {
	npcBase, exists := e.world.Npcs[npcName]
	if !exists || !isNpcInRoom(room, npcName) {
		return nil, errors.New(pr.ErrNpcNotFound)
	}
	if !npcBase.Hostile {
		return nil, errors.New(pr.ErrNpcNotHostile)
	}
	return npcBase, nil
}

func (e *Engine) initiateCombat(player *Player, npcName string) (*CombatSession, error) {
	npcBase, err := e.getValidTarget(player.room, npcName)
	if err != nil {
		return nil, err
	}

	combatID := uuid.New().String()
	combatNpc := npcBase.Clone()

	combat_session := &CombatSession{Id: combatID, Fighters: []Fighter{}}
	group, exists := e.getPlayerGroup(player)

	combat_session.addNpcToCombat(combatNpc)
	if !exists {
		combat_session.addPlayerToCombat(player)
	} else {
		for _, p := range group.players {
			combat_session.addPlayerToCombat(p)
		}
	}

	e.activeCombats[combatID] = combat_session

	return combat_session, nil
}

func (cs *CombatSession) addPlayerToCombat(player *Player) {
	player.inCombat = true
	player.stats.CombatId = cs.Id
	cs.Fighters = append(cs.Fighters, player)
}

func (cs *CombatSession) addNpcToCombat(npc *Npc) {
	npc.InCombat = true
	npc.Stats.CombatId = cs.Id
	cs.Fighters = append(cs.Fighters, npc)

}

func (e *Engine) processCombatTurn(player *Player, targetName string) (string, *CombatTurnResult, error) {
	if !player.inCombat || player.stats.CombatId == "" {
		return "OK", nil, errors.New(pr.ErrInternalServer)
	}

	combatSession, exists := e.activeCombats[player.stats.CombatId]
	if !exists {
		player.inCombat = false
		return "", nil, errors.New(pr.ErrInternalServer)
	}
	target, targetExists := combatSession.Enemies[targetName]
	if !targetExists {
		return "", nil, errors.New(pr.ErrNpcNotFound)
	}
	current_turn := combatSession.CurrentTurn
	if player.name != current_turn {
		return "", nil, errors.New(pr.ErrNotYourTurnToPlay)
	}

	return "" + targetName, *CombatTurnResult, nil
}
func (e *Engine) nextToPlay(cs *CombatSession) {

}
