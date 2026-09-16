package variables

import (
	"tap/client/state"
	"tap/protocol"
)

type PanelVariables struct {
	Room      *protocol.LookCommandData
	RoomItems *[]string
	Inventory *[]string
	Quests    *[]protocol.TrackedQuestData

	GroupState  *state.GroupState
	CombatState *state.CombatState

	Chat *ChatPanel
}
