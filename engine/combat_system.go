package engine

import (
	"errors"
	"fmt"
	"slices"
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
	Id           string
	Fighters     []Fighter
	Npcs         []*Npc
	Players      []*Player
	CurrentTurn  int
	State        CombatState
	TurnResponse *FullTurnResponse
	RoomId       string
	Engine       *Engine
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
	if !isNpcInRoom(room, targetName) || slices.Contains(player.DefeatedNpcs, targetName) {
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

func (e *Engine) initiateCombat(player *Player, npc_copy *Npc) (*CombatSession, string, *FullTurnResponse) {

	combatID := uuid.New().String()

	cs := &CombatSession{Id: combatID, Fighters: []Fighter{}, State: StateOngoing, RoomId: player.room.Id, Engine: e, CurrentTurn: -1}
	group, exists := e.groups[player.group]

	cs.addNpcToCombat(npc_copy)
	if !exists {
		cs.addPlayerToCombat(player)
	} else {
		for _, p := range group.players {
			// add players in the combat only if they are in the same room
			if p.room == player.room {
				if p.name != player.name && !slices.Contains(p.DefeatedNpcs, npc_copy.Id) {
					e.inform_user(p, "EVT COMBAT FIGHT_STARTED")
				}
				if !slices.Contains(p.DefeatedNpcs, npc_copy.Id) {

					cs.addPlayerToCombat(p)
				}
			} else {
				if p.name != player.name {
					// else inform the others that some people in the group entered in combat.

					e.inform_user(p, "EVT GROUP DISTANT_ALLIES_COMBAT_START")
				}
			}
		}
	}
	cs.sortTurnsOrderByInitiative()
	e.activeCombats[combatID] = cs

	res, full_turn_res := cs.processCombatTurn(player, npc_copy)

	return cs, res, full_turn_res

}

func (cs *CombatSession) processCombatTurn(attacker Fighter, target Fighter) (string, *FullTurnResponse) {
	response := &FullTurnResponse{
		NpcReactions: []ActionLog{},
		CombatState:  cs.State,
	}
	cs.TurnResponse = response

	inflicted_damage := target.takeDamage(attacker.getDamage())
	if cs.checkIfPlayersAreDead() {
		cs.State = StateDefeat
		cs.TurnResponse.CombatState = cs.State
	}
	if cs.checkIfNpcsAreDead() {
		cs.State = StateVictory
		cs.TurnResponse.CombatState = cs.State
	}
	player_turn_result := &CombatTurnResult{AttackerHp: attacker.getHp(), TargetHp: target.getHp(), Damage: inflicted_damage, Status: cs.State}
	cs.TurnResponse.PlayerAction = ActionLog{
		ActorName:  attacker.getName(),
		TargetName: target.getName(),
		Result:     player_turn_result,
	}
	cs.nextTurn()

	if cs.State != StateOngoing {
		return "OK", response
	}
	cs.processNpcsTurn()

	return "OK", cs.TurnResponse
}

func (cs *CombatSession) processNpcsTurn() {
	for cs.State == StateOngoing {
		current_fighter := cs.Fighters[cs.CurrentTurn]

		npcFighter, ok := current_fighter.(*Npc)
		if !ok {
			break
		}
		// choose target between players
		var current_target Fighter
		for _, p := range cs.Players {
			if p.inCombat && p.room.Id == cs.RoomId && !p.isDead() {
				current_target = p
				break
			}
		}
		if current_target == nil {
			return
		}
		inflicted := current_target.takeDamage(npcFighter.getDamage())
		if cs.checkIfPlayersAreDead() {
			cs.State = StateDefeat
			cs.TurnResponse.CombatState = cs.State
		}
		if cs.checkIfNpcsAreDead() {
			cs.State = StateVictory
			cs.TurnResponse.CombatState = cs.State
		}
		cs.nextTurn()
		cs.TurnResponse.CombatState = cs.State
		cs.TurnResponse.NpcReactions = append(cs.TurnResponse.NpcReactions, ActionLog{ActorName: npcFighter.Name, TargetName: current_target.getName(), Result: &CombatTurnResult{AttackerHp: npcFighter.getHp(), TargetHp: current_target.getHp(), Damage: inflicted, Status: cs.State}})
	}
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

func (cs *CombatSession) nextTurn() {
	maxLoops := len(cs.Fighters)

	for range maxLoops {
		cs.CurrentTurn++
		if cs.CurrentTurn >= len(cs.Fighters) {
			cs.CurrentTurn = 0
		}

		currentFighter := cs.Fighters[cs.CurrentTurn]
		canPlay := true

		switch f := currentFighter.(type) {
		case *Player:
			if f.isDead() || !f.inCombat {
				canPlay = false
			}
		case *Npc:
			if f.isDead() {
				canPlay = false
			}
		}
		if canPlay {
			current_player := cs.Fighters[cs.CurrentTurn]
			msg := fmt.Sprintf("EVT COMBAT TURN %s", current_player.getName())
			cs.Engine.inform_combat_players(cs, nil, msg)
			break
		}
	}
}

func (cs *CombatSession) isFighterTurn(fighter Fighter) bool {
	// currentTurnIndex() renvoie l'index exact de celui qui doit jouer
	expectedFighter := cs.Fighters[cs.CurrentTurn]

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

func (e *Engine) end_combat(cs *CombatSession) {
	if cs.State == StateVictory {
		e.inform_combat_players(cs, nil, "EVT COMBAT VICTORY")
	}
	for _, player := range cs.Players {
		player.stats.CombatId = ""
		player.inCombat = false
		// if lost, punished by reducing its HP and sending them to FIRST ROOM
		if cs.State == StateDefeat {
			player.stats.Hp = player.stats.HpMax / 2
			player.room = e.world.Rooms[RoomEntrance]
			msg := fmt.Sprintf("EVT COMBAT DEFEAT new_room=%s", RoomEntrance)
			e.inform_combat_players(cs, nil, msg)
		}
		if cs.State == StateVictory {
			//reward
			for _, npc := range cs.Npcs {
				fmt.Println("npc ID", npc.Id)
				player.DefeatedNpcs = append(player.DefeatedNpcs, npc.Id)
			}
		}
		delete(e.activeCombats, cs.Id)
	}
}

func (cs *CombatSession) leaveCombat(player *Player) error {
	index := slices.Index(cs.Fighters, Fighter(player))

	if index == -1 {
		return errors.New(pr.ErrInternalServer)
	}
	if cs.CurrentTurn != index {
		return errors.New(pr.ErrNotYourTurnToPlay)
	}
	player.inCombat = false
	player.stats.CombatId = ""
	cs.nextTurn()
	cs.processNpcsTurn()
	return nil
}
