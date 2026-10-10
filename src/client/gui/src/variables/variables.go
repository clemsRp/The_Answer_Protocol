package variables

import (
	"time"
)

type Size struct {
	Width  float32
	Height float32
}

type Variables struct {
	Player             *Player
	RemotePlayers      *map[string]*Player
	RemotePlayerEmotes *map[string]*PlayerEmote
	ItemPositions      *map[string]*Position
	NpcPositions       *map[string]*Position
	Collisions         [][]bool
	Tileset_size       float32
	LastRoomChange     time.Time

	Current_room      string
	LastRoomSpanwTime time.Time
	Current_view      string
	PanelsVariables   *PanelVariables

	Npcs          map[string]*NpcDatas
	Items         map[string]*ItemDatas
	NpcConvertor  map[string]string
	ItemConvertor map[string]string

	ConnectionError string

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
		Current_room:    "place_du_village",
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
