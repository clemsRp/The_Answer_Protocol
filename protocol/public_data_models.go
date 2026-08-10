package protocol

type PlayerInfo struct {
	Name      string   `json:"name"`
	Location  string   `json:"location"`
	Hp        int      `json:"hp"`
	HpMax     int      `json:"hp_max"`
	Room      string   `json:"room"`
	Group     string   `json:"group"`
	Inventory []string `json:"inventory"`
}
