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
	Id          string `json:"-" validate:"-"`
	Description string `json:"description"`
	Reward      string `json:"reward"`
	Status      string `json:"status" validate:"required,valid_quest_status"`
	Progress    string `json:"progress"`
	TargetItem  string `json:"target_item,omitempty"`
	TargetNpc   string `json:"target_npc,omitempty"`
}

func (q *Quest) Clone() *Quest {
	cp := *q
	return &cp
}

type Room struct {
	Id          string            `json:"-" validate:"-"`
	Name        string            `json:"name" validate:"required"`
	Description string            `json:"description"`
	Exits       map[string]string `json:"exits" validate:"dive,keys,valid_exit,endkeys,room_exists"`
	Items       []string          `json:"items" validate:"dive,item_exists"`
	Npcs        []string          `json:"npcs" validate:"dive,npc_exists"`
}

type Map struct {
	Rooms  map[string]*Room  `json:"rooms" validate:"dive"`
	Items  map[string]*Item  `json:"items" validate:"dive"`
	Npcs   map[string]*Npc   `json:"npcs" validate:"dive"`
	Quests map[string]*Quest `json:"quests" validate:"dive"`
}
