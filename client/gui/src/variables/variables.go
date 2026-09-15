package variables

import (
	"tap/client/state"
	"tap/protocol"
	"time"
)

type Size struct {
	Width  float32
	Height float32
}

type Direction struct {
	X float32
	Y float32
}

type Position struct {
	X float32
	Y float32
}

type Player struct {
	Direction     *Direction
	Position      *Position
	Speed         int
	Pseudo        string
	EmoteIndex    int
	LastTimeTyped time.Time
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
	Player        *Player
	RemotePlayers map[string]*Player
	Collisions    [][]bool
	Tileset_size  float32

	Current_room    string
	Current_view    string
	PanelsVariables *PanelVariables

	Zoom         float32
	FontSize     float32
	StartTime    time.Time
	MapStart     *Position
	MapSize      *Size
	StartingPosX float32
	StartingPosY float32
}

func GetVariables() *Variables {
	return &Variables{
		Player: &Player{
			Direction: &Direction{
				X: 0,
				Y: 1,
			},
			EmoteIndex: 2,
		},
		RemotePlayers: make(map[string]*Player),
		Current_room:  "entrance",
		Current_view:  "Connect",
		PanelsVariables: &PanelVariables{
			Room:        &protocol.LookCommandData{},
			RoomItems:   &[]string{},
			Inventory:   &[]string{},
			Quests:      &[]protocol.TrackedQuestData{},
			GroupState:  &state.GroupState{},
			CombatState: &state.CombatState{},
		},
		MapStart: &Position{
			X: float32(0),
			Y: float32(0),
		},
		MapSize: &Size{
			Width:  float32(1),
			Height: float32(1),
		},
	}
}
