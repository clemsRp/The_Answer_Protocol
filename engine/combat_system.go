package engine

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	pr "tap/protocol"

	"github.com/google/uuid"
)

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
	// If player is IN combat, verify target within active combat fighters
	if player.inCombat {
		cs, exists := e.activeCombats[player.stats.CombatId]
		if !exists {
			return nil, errors.New(pr.ErrInternalServer)
		}
		for _, fighter := range cs.Fighters {
			if p, ok := fighter.(*Player); ok {
				if p.name == targetName {
					return nil, errors.New(pr.ErrNoAllyAttack)
				}
			}
			if npc, ok := fighter.(*Npc); ok {
				if npc.Id == targetName || npc.Name == targetName {
					return npc, nil
				}
			}
		}
		return nil, errors.New(pr.ErrNpcNotFound)
	}

	// If player is NOT in combat, check if NPC is already in an active combat in the room
	for _, cs := range e.activeCombats {
		if cs.RoomId == player.room.Id {
			for _, n := range cs.Npcs {
				if n.Id == targetName || n.Name == targetName {
					return n, nil
				}
			}
		}
	}

	// Otherwise, verify if it's a valid base NPC in the room
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

	return npcBase.Clone(), nil
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
			if p.room == player.room {
				if p.name != player.name && !slices.Contains(p.DefeatedNpcs, npc_copy.Id) {
					e.inform_user(p, "EVT COMBAT FIGHT_STARTED")
				}
				if !slices.Contains(p.DefeatedNpcs, npc_copy.Id) {
					cs.addPlayerToCombat(p)
				}
			} else {
				if p.name != player.name {
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
		var current_target Fighter
		for _, p := range cs.Players {
			if p.inCombat && p.room.Id == cs.RoomId && !p.isDead() {
				current_target = p
				break
			}
		}

		if current_target == nil {
			cs.State = StateCancelled
			if cs.TurnResponse != nil {
				cs.TurnResponse.CombatState = cs.State
			}
			break
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
	for _, player := range cs.Players {
		if !player.isDead() {
			return false
		}
	}
	return true
}

func (cs *CombatSession) checkIfNpcsAreDead() bool {
	for _, npc := range cs.Npcs {
		if !npc.isDead() {
			return false
		}
	}
	return true
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
			msg := fmt.Sprintf("%s %s %s %s", pr.MsgEvt, pr.CategoryCombat, pr.TypeAllyTurn, current_player.getName())
			cs.Engine.inform_combat_players(cs, nil, msg)
			break
		}
	}
}

func (cs *CombatSession) isFighterTurn(fighter Fighter) bool {
	if cs.CurrentTurn < 0 || cs.CurrentTurn >= len(cs.Fighters) {
		return false
	}
	expectedFighter := cs.Fighters[cs.CurrentTurn]
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
	if cs.State == StateDefeat {
		msg := fmt.Sprintf("EVT COMBAT DEFEAT new_room=%s", RoomEntrance)
		for _, p := range cs.Players {
			e.inform_user(p, msg)
		}
	} else if cs.State == StateVictory {
		for _, p := range cs.Players {
			e.inform_user(p, "EVT COMBAT VICTORY")
		}
	}

	for _, player := range cs.Players {
		player.stats.CombatId = ""
		player.inCombat = false
		if cs.State == StateDefeat {
			player.stats.Hp = player.stats.HpMax / 2
			player.room = e.world.Rooms[RoomEntrance]
		}
		if cs.State == StateVictory {
			for _, npc := range cs.Npcs {
				if !slices.Contains(player.DefeatedNpcs, npc.Id) {
					player.DefeatedNpcs = append(player.DefeatedNpcs, npc.Id)
				}
			}
			// A freshly defeated npc may fulfil an active quest target, so
			// recompute progress right away instead of waiting for the
			// player to ask for it.
			e.refreshQuestProgress(player)
		}
	}
	delete(e.activeCombats, cs.Id)
}

func (cs *CombatSession) leaveCombat(player *Player) error {
	if !cs.isFighterTurn(player) {
		return errors.New(pr.ErrNotYourTurnToPlay)
	}

	player.inCombat = false
	player.stats.CombatId = ""

	if cs.TurnResponse == nil {
		cs.TurnResponse = &FullTurnResponse{
			NpcReactions: []ActionLog{},
			CombatState:  cs.State,
		}
	}

	playersLeft := false
	for _, p := range cs.Players {
		if p.inCombat {
			playersLeft = true
			break
		}
	}

	if !playersLeft {
		cs.State = StateCancelled
		cs.TurnResponse.CombatState = cs.State
		return nil
	}

	cs.nextTurn()
	cs.processNpcsTurn()
	return nil
}
