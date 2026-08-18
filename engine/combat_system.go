package engine

import (
	"errors"
	"fmt"
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
	Npcs        []*Npc
	Players     []*Player
	CurrentTurn int
	Fighting    bool
	State       CombatState
}
type CombatState string

const (
	StateOngoing   CombatState = "ONGOING"
	StateVictory   CombatState = "VICTORY"
	StateDefeat    CombatState = "DEFEAT"
	StateCancelled CombatState = "CANCELLED"
)

type CombatTurnResult struct {
	AttackerHp int         `json:"attacker_hp"`
	TargetHp   int         `json:"target_hp"`
	Damage     int         `json:"damage"`
	Status     CombatState `json:"status"`
}

type ActionLog struct {
	ActorName  string            `json:"actor_name"`
	TargetName string            `json:"target_name"`
	Result     *CombatTurnResult `json:"result"`
}

type FullTurnResponse struct {
	PlayerAction ActionLog   `json:"player_action"`
	NpcReactions []ActionLog `json:"npc_reactions"`
	CombatState  CombatState `json:"combat_state"`
}

func (e *Engine) getValidTarget(player *Player, targetName string) (*Npc, error) {
	// if npc doesn't exist in word
	npcBase, exists := e.world.Npcs[targetName]
	if !exists {
		return nil, errors.New(pr.ErrNpcNotFound)
	}
	room, ok := e.world.Rooms[player.room.Id]
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
	cs, exists := e.activeCombats[player.stats.CombatId]
	if !exists {
		fmt.Println("combat pb")
		return nil, errors.New(pr.ErrInternalServer)
	}
	for _, fighter := range cs.Fighters {
		switch verifTarget := fighter.(type) {
		case *Player:
			return nil, errors.New(pr.ErrNoAllyAttack)
		case *Npc:
			// return target clone
			return verifTarget, nil
		default:
			return nil, errors.New(pr.ErrInternalServer)
		}
	}
	return npcBase, nil
}

func (e *Engine) initiateCombat(player *Player, npc_copy *Npc) (string, *FullTurnResponse, error) {

	combatID := uuid.New().String()

	cs := &CombatSession{Id: combatID, Fighters: []Fighter{}, State: StateOngoing}
	group, exists := e.groups[player.group]

	cs.addNpcToCombat(npc_copy)
	if !exists {
		cs.addPlayerToCombat(player)
	} else {
		for _, p := range group.players {
			cs.addPlayerToCombat(p)
		}
	}
	cs.sortTurnsOrderByInitiative()
	e.activeCombats[combatID] = cs

	return cs.processCombatTurn(player, npc_copy, true)

}

func (cs *CombatSession) processCombatTurn(attacker Fighter, target Fighter, isFirstAttack bool) (string, *FullTurnResponse, error) {
	response := &FullTurnResponse{
		NpcReactions: []ActionLog{},
		CombatState:  cs.State,
	}

	inflicted_damage := target.takeDamage(attacker.getDamage())
	if cs.checkIfPlayersAreDead() {
		cs.State = StateDefeat
		response.CombatState = cs.State
	}
	if cs.checkIfNpcsAreDead() {
		cs.State = StateVictory
		response.CombatState = cs.State
	}
	player_turn_result := &CombatTurnResult{AttackerHp: attacker.getHp(), TargetHp: target.getHp(), Damage: inflicted_damage, Status: cs.State}
	response.PlayerAction = ActionLog{
		ActorName:  attacker.getName(),
		TargetName: target.getName(),
		Result:     player_turn_result,
	}
	if !isFirstAttack {
		cs.CurrentTurn++
	}

	if cs.State != StateOngoing {
		return "OK", response, nil
	}
	for cs.State == StateOngoing {
		current_fighter := cs.Fighters[cs.currentTurnIndex()]

		npcFighter, ok := current_fighter.(*Npc)
		if !ok {
			break
		}
		var current_target Fighter
		for _, p := range cs.Players {
			if !p.isDead() {
				current_target = p
				break
			}
		}
		inflicted := current_target.takeDamage(npcFighter.getDamage())
		if cs.checkIfPlayersAreDead() {
			cs.State = StateDefeat
			response.CombatState = cs.State
		}
		if cs.checkIfNpcsAreDead() {
			cs.State = StateVictory
			response.CombatState = cs.State
		}
		response.CombatState = cs.State
		response.NpcReactions = append(response.NpcReactions, ActionLog{ActorName: npcFighter.Name, TargetName: current_target.getName(), Result: &CombatTurnResult{AttackerHp: npcFighter.getHp(), TargetHp: current_target.getHp(), Damage: inflicted, Status: cs.State}})

		cs.CurrentTurn++
	}
	return "OK", response, nil
}

func (cs *CombatSession) sortTurnsOrderByInitiative() {
	sort.Slice(cs.Fighters, func(i, j int) bool {
		initiativeI := cs.Fighters[i].getInitiative()
		initiativeJ := cs.Fighters[j].getInitiative()
		return initiativeI > initiativeJ
	})
}

func (cs *CombatSession) checkIfPlayersAreDead() bool {
	allDead := true
	for _, player := range cs.Players {
		if !player.isDead() {
			allDead = false
		}
	}
	return allDead
}

func (cs *CombatSession) checkIfNpcsAreDead() bool {
	allDead := true
	for _, npc := range cs.Npcs {
		if !npc.isDead() {
			allDead = false
		}
	}
	return allDead
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

func (cs *CombatSession) addPlayerToCombat(player *Player) {
	player.inCombat = true
	player.stats.CombatId = cs.Id
	cs.Fighters = append(cs.Fighters, player)
	cs.Players = append(cs.Players, player)
}

func (cs *CombatSession) addNpcToCombat(npc *Npc) {
	npc.InCombat = true
	npc.Stats.CombatId = cs.Id
	cs.Fighters = append(cs.Fighters, npc)
	cs.Npcs = append(cs.Npcs, npc)

}
