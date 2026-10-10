package variables

import (
	"tap/src/client/state"
	"tap/src/protocol"
	"time"
)

type PanelVariables struct {
	Room           *protocol.LookCommandData
	RoomItems      *[]string
	InventoryItems *[]string
	Quests         *[]protocol.TrackedQuestData

	GroupState  *state.GroupState
	CombatState *state.CombatState

	Chat         *ChatPanel
	Datas        *DatasPanel
	Group        *GroupPanel
	CombatResult *CombatResultPanel
	Inspect      *InspectPanel
	Talk         *TalkPanel
	Emotes       *EmotesPanel

	MiniMapOpen bool
}

var (
	results = make([]string, 0)
)

func GetPanelsVariables() *PanelVariables {
	return &PanelVariables{
		MiniMapOpen:    true,
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
			UnreadByScope:   make(map[string]int),
			LastNbChats:     0,
			ScrollActive:    false,
			LastFrameScroll: false,
		},
		Datas: &DatasPanel{
			Open:            false,
			NbServerPlayers: 0,
		},
		Group: &GroupPanel{
			Open:            false,
			InGroup:         false,
			Leader:          "None",
			Grouped:         make([]string, 0),
			UnGrouped:       make([]string, 0),
			Invitations:     make([]string, 0),
			SendInvitations: make([]string, 0),
			SendPromotion:   "",
			Promote:         false,

			ScrollActive:    false,
			LastFrameScroll: false,
		},
		CombatResult: &CombatResultPanel{},
		Inspect: &InspectPanel{
			Open:        false,
			LastInspect: "ROOM",
		},
		Talk: &TalkPanel{
			Talking:  nil,
			LastTalk: "",
			Finished: false,
			Results:  &results,
		},
		Emotes: &EmotesPanel{
			Open:           false,
			LastEmoteIndex: -1,
			LastEmoteTime:  time.Now(),
		},
	}
}
