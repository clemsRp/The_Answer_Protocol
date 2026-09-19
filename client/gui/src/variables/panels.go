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

	Chat      *ChatPanel
	LeftPanel *LeftPanel
	Inventory *InventoryPanel
}

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
		LeftPanel: &LeftPanel{
			Open: false,
		},
		Inventory: &InventoryPanel{
			Open:    true,
			NbItems: 0,
		},
	}
}
