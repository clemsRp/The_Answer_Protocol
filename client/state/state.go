package state

import (
	"sync"
	"tap/protocol"
)

type GameState struct {
	mu     sync.RWMutex
	Player *Player
	Combat *CombatState
}

type Player struct {
	Name           string
	Room           *protocol.LookCommandData
	Inventory      []string
	Quests         []string
	DefeatedNpcs   []string
	InCombat       bool
	EquippedWeapon string
	Stats          string
	GroupState     *GroupState
}

type GroupState struct {
	Group         string
	LastKick      *string
	Leader        bool
	Promotion     bool
	SendPromotion bool
	Grouped       []string
	UnGrouped     []string
	Invitations   []string
}

type CombatChat struct {
	Pseudo string
	Msg    string
}

type CombatState struct {
	InCombat       bool
	Chats          []CombatChat
	LastCombatChat string
	Opponents      map[string]protocol.CombatPersonData
	Team           map[string]protocol.CombatPersonData
	Leader         string
	CurrentTurn    string
	SelectedPerson string
}

func New() *GameState {
	return &GameState{
		Player: &Player{
			GroupState: &GroupState{},
			Room:       &protocol.LookCommandData{},
			Inventory:  make([]string, 0),
		},
		Combat: &CombatState{
			Chats:     make([]CombatChat, 0),
			Opponents: make(map[string]protocol.CombatPersonData),
			Team:      make(map[string]protocol.CombatPersonData),
		},
	}
}

// Read executes a function under a read lock.
func (gs *GameState) Read(fn func(state *GameState)) {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	fn(gs)
}

// Snapshot getters

func (gs *GameState) GetPlayerSnapshot() Player {
	gs.mu.RLock()
	defer gs.mu.RUnlock()

	if gs.Player == nil {
		return Player{}
	}

	snap := *gs.Player
	snap.Inventory = append([]string(nil), gs.Player.Inventory...)
	snap.Quests = append([]string(nil), gs.Player.Quests...)
	snap.DefeatedNpcs = append([]string(nil), gs.Player.DefeatedNpcs...)

	snap.Room = nil
	snap.GroupState = nil

	return snap
}

func (gs *GameState) GetGroupSnapshot() GroupState {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	if gs.Player == nil || gs.Player.GroupState == nil {
		return GroupState{}
	}
	snap := *gs.Player.GroupState

	snap.Grouped = append([]string(nil), gs.Player.GroupState.Grouped...)
	snap.UnGrouped = append([]string(nil), gs.Player.GroupState.UnGrouped...)
	snap.Invitations = append([]string(nil), gs.Player.GroupState.Invitations...)

	return snap
}

func (gs *GameState) GetCombatSnapshot() CombatState {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	if gs.Combat == nil {
		return CombatState{}
	}
	snap := *gs.Combat
	snap.Chats = append([]CombatChat(nil), gs.Combat.Chats...)

	snap.Opponents = make(map[string]protocol.CombatPersonData, len(gs.Combat.Opponents))
	for k, v := range gs.Combat.Opponents {
		snap.Opponents[k] = v
	}

	snap.Team = make(map[string]protocol.CombatPersonData, len(gs.Combat.Team))
	for k, v := range gs.Combat.Team {
		snap.Team[k] = v
	}

	return snap
}

// Updaters

func (gs *GameState) UpdatePlayer(updateFn func(p *Player)) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if gs.Player == nil {
		gs.Player = &Player{}
	}

	p := *gs.Player
	updateFn(&p)
	gs.Player = &p
}

func (gs *GameState) UpdateRoomLook(newRoom *protocol.LookCommandData) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if gs.Player == nil {
		gs.Player = &Player{}
	}

	gs.Player.Room = newRoom
}

func (gs *GameState) UpdateRoom(updateFn func(r *protocol.LookCommandData)) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if gs.Player == nil {
		gs.Player = &Player{}
	}
	if gs.Player.Room == nil {
		gs.Player.Room = &protocol.LookCommandData{}
	}

	r := *gs.Player.Room
	updateFn(&r)
	gs.Player.Room = &r
}

func (gs *GameState) UpdateGroupState(updateFn func(gs *GroupState)) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if gs.Player == nil {
		gs.Player = &Player{}
	}

	if gs.Player.GroupState == nil {
		gs.Player.GroupState = &GroupState{}
	}

	newState := *gs.Player.GroupState
	updateFn(&newState)
	gs.Player.GroupState = &newState
}

func (gs *GameState) UpdateCombatState(updateFn func(cs *CombatState)) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if gs.Combat == nil {
		gs.Combat = &CombatState{
			Chats:     make([]CombatChat, 0),
			Opponents: make(map[string]protocol.CombatPersonData),
			Team:      make(map[string]protocol.CombatPersonData),
		}
	}

	cs := *gs.Combat
	updateFn(&cs)
	gs.Combat = &cs
}

func (s *GameState) UpdatePlayerLeaveRoom(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Player != nil && s.Player.Room != nil {
		for i, name := range s.Player.Room.Players {
			if name == username {
				s.Player.Room.Players = append(s.Player.Room.Players[:i], s.Player.Room.Players[i+1:]...)
				break
			}
		}
	}
}

func (s *GameState) UpdatePlayerEnterRoom(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Player != nil && s.Player.Room != nil {
		s.Player.Room.Players = append(s.Player.Room.Players, username)
	}
}
