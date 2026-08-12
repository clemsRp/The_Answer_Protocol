package engine

// type PlayerInfo struct {
// 	Name      string   `json:"name"`
// 	Location  string   `json:"location"`
// 	Hp        int      `json:"hp"`
// 	HpMax     int      `json:"hp_max"`
// 	Room      string   `json:"room"`
// 	Group     string   `json:"group"`
// 	Inventory []string `json:"inventory"`
// 	Quests    []Quest  `json:"quests"`
// }

type Quest struct {
	Description string `json:"description"`
	Reward      string `json:"reward"`
	Status      string `json:"status"`
}

type Room struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Exits       map[string]string `json:"exits"`
	Items       []string          `json:"items"`
	Npcs        []string          `json:"npcs"`
}

type Map struct {
	Rooms  map[string]*Room  `json:"rooms"`
	Items  map[string]*Item  `json:"items"`
	Npcs   map[string]*Npc   `json:"npcs"`
	Quests map[string]*Quest `json:"quests"`
}


