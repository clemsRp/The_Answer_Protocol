package variables

import (
	"tap/client/state"
	"tap/protocol"
)

type PanelVariables struct {
	Room           *protocol.LookCommandData
	RoomItems      *[]string
	InventoryItems *[]string
	Quests         *[]protocol.TrackedQuestData

	GroupState  *state.GroupState
	CombatState *state.CombatState

	Chat  *ChatPanel
	Group *GroupPanel
	Talk  *TalkPanel
}

var (
	results = make([]string, 0)
)

func GetPanelsVariables() *PanelVariables {
	return &PanelVariables{
		Room:           &protocol.LookCommandData{},
		RoomItems:      &[]string{},
		InventoryItems: &[]string{},
		Quests:         &[]protocol.TrackedQuestData{},
		GroupState:     &state.GroupState{},
		CombatState:    &state.CombatState{},

		Chat: &ChatPanel{
			Open:            false,
			CurrentScope:    "GLOBAL",
			ScopeChats:      make(map[string][]Chat),
			LastNbChats:     0,
			ScrollActive:    false,
			LastFrameScroll: false,
		},
		Group: &GroupPanel{
			Open:          false,
			InGroup:       false,
			Leader:        "",
			Grouped:       make([]string, 0),
			UnGrouped:     make([]string, 0),
			Invitations:   make([]string, 0),
			SendPromotion: false,
			Promote:       false,
		},
		Talk: &TalkPanel{
			Talking:  nil,
			LastTalk: "",
			Finished: false,
			Results:  &results,
		},
	}
}
