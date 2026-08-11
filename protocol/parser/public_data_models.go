package parser

type PlayerInfo struct {
	Name      string   `json:"name"`
	Location  string   `json:"location"`
	Hp        int      `json:"hp"`
	HpMax     int      `json:"hp_max"`
	Room      string   `json:"room"`
	Group     string   `json:"group"`
	Inventory []string `json:"inventory"`
	Quests    []Quest  `json:"quests"`
}

type Stats struct {
	Hp     int    `json:"hp"`
	MaxHp  int    `json:"max_hp"`
	Status string `json:"status"`
}

type Quest struct {
	Description string `json:"description"`
	Reward      string `json:"reward"`
	Status      string `json:"status"`
}

type Npc struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Dialogue    []string `json:"dialogue"`
	Role        string   `json:"role"`
	QuestId     string   `json:"quest_id"`
	Stats       Stats    `json:"stats"`
	Hostile     bool     `json:"hostile"`
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
type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Obtainable  bool   `json:"obtainable"`
	Type        string `json:"type"`
	Tradable    bool   `json:"tradable"`
	Worth       int    `json:"worth"`
}

type Weapon struct {
	Damage   int
	Worth    int
	Tradable bool
}

type Ressource struct {
	Worth    int
	Tradable bool
}
