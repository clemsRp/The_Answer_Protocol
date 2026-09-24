package variables

import (
	"time"
)

type Size struct {
	Width  float32
	Height float32
}

type Variables struct {
	Player        *Player
	RemotePlayers *map[string]*Player
	ItemPositions *map[string]*Position
	Collisions    [][]bool
	Tileset_size  float32

	Current_room    string
	Current_view    string
	PanelsVariables *PanelVariables

	Npcs         map[string]*NpcDatas
	NpcConvertor map[string]string

	Zoom         float32
	FontSize     float32
	StartTime    time.Time
	MapStart     *Position
	MapSize      *Size
	StartingPosX float32
	StartingPosY float32
	Mouse        *Mouse
}

func GetVariables() *Variables {
	return &Variables{
		Player:          GetPlayerVariables(),
		Current_room:    "entrance",
		Current_view:    "Connect",
		PanelsVariables: GetPanelsVariables(),
		Npcs:            make(map[string]*NpcDatas),
		MapStart: &Position{
			X: float32(0),
			Y: float32(0),
		},
		MapSize: &Size{
			Width:  float32(1),
			Height: float32(1),
		},
		FontSize: 40,
		Mouse: &Mouse{
			Clicked: false,
			Down:    false,
		},
	}
}
