package engine

import (
	"errors"
	"sort"
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

func (e *Engine) getValidTarget(player *Player, targetName string) (*Npc, error) {
	// if npc doesn't exist in word
	npcBase, exists := e.world.Npcs[targetName]
	if !exists {
		return nil, errors.New(pr.ErrNpcNotFound)
	}
	room, ok := e.world.Rooms[player.room.Name]
	if !ok {
		return nil, errors.New(pr.ErrInternalServer)
	}
	if !isNpcInRoom(room, targetName) {
		return nil, errors.New(pr.ErrNpcNotFound)
	}
	if !npcBase.Hostile {
		return nil, errors.New(pr.ErrNpcNotHostile)
	}

	// if player is not in combat return npc found that has been attacked
	if !player.inCombat {
		return npcBase.Clone(), nil
	}
	// if player in combat, check if targetName is in Fighters of the combat and is a NPC
	combat_session, exists := e.activeCombats[player.stats.CombatId]
	if !exists {
		return nil, errors.New(pr.ErrInternalServer)
	}
	for _, fighter := range combat_session.Fighters {
		switch verifTarget := fighter.(type) {
		case *Player:
			return nil, errors.New(pr.ErrNoAllyAttack)
		case *Npc:
			// return target clone
			return verifTarget.Clone(), nil
		default:
			return nil, errors.New(pr.ErrInternalServer)
		}
	}
	return npcBase, nil
}

func (e *Engine) initiateCombat(player *Player, npcName string) (string, *CombatTurnResult, error) {
	npcBase, err := e.getValidTarget(player, npcName)
	if err != nil {
		return "", nil, err
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

	//TODO
	combat_session.sortTurnsOrderByInitiative()
	// play a turn once here
	return combat_session.processCombatTurn(player, combatNpc)

}

func (cs *CombatSession) sortTurnsOrderByInitiative() {
	sort.Slice(cs.Fighters, func(i, j int) bool {
		initiativeI := cs.Fighters[i].getInitiative()
		initiativeJ := cs.Fighters[j].getInitiative()
		return initiativeI > initiativeJ
	})
}

// infinite cycle on even on current Turn == 100
func (cs *CombatSession) currentTurnIndex() int {
	return cs.CurrentTurn % len(cs.Fighters)
}

func (cs *CombatSession) isFighterTurn(fighter Fighter) bool {
	// currentTurnIndex() renvoie l'index exact de celui qui doit jouer
	expectedFighter := cs.Fighters[cs.currentTurnIndex()]

	// On compare simplement les noms via l'interface
	return expectedFighter.getName() == fighter.getName()
}

func (cs *CombatSession) processCombatTurn(attacker Fighter, target Fighter) (string, *CombatTurnResult, error) {
	if !cs.isFighterTurn(attacker) {
		return "", nil, errors.New(pr.ErrNotYourTurnToPlay)
	}

	turn_result := attacker.playCombatTurn(target)
	if target.isDead() {
		cs.Win = true
	} else {
		cs.CurrentTurn++
	}
	return "OK", turn_result, nil
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
