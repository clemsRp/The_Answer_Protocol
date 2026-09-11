package variables

import (
	"tap/client/state"
	"tap/protocol"
	"time"
)

type Direction struct {
	X float32
	Y float32
}

type Position struct {
	X float32
	Y float32
}

type Player struct {
	Direction *Direction
	Position  *Position
}

type PanelVariables struct {
	Room      *protocol.LookCommandData
	RoomItems *[]string
	Inventory *[]string
	Quests    *[]protocol.TrackedQuestData

	GroupState  *state.GroupState
	CombatState *state.CombatState
}

type Variables struct {
	Player       *Player
	Collisions   [][]bool
	Tileset_size int

	Current_room    string
	Current_view    string
	PanelsVariables *PanelVariables

	Zoom      float32
	StartTime time.Time
}

func GetVariables() *Variables {
	return &Variables{
		Player: &Player{
			Direction: &Direction{
				X: 0,
				Y: 1,
			},
			Position: &Position{
				X: float32(500),
				Y: float32(400),
			},
		},
		Current_room: "entrance",
		Current_view: "Connect",
		PanelsVariables: &PanelVariables{
			Room:        &protocol.LookCommandData{},
			RoomItems:   &[]string{},
			Inventory:   &[]string{},
			Quests:      &[]protocol.TrackedQuestData{},
			GroupState:  &state.GroupState{},
			CombatState: &state.CombatState{},
		},
	}
}
