package protocol

type ExitsData struct {
	North string `json:"north,omitempty"`
	East  string `json:"east,omitempty"`
	West  string `json:"west,omitempty"`
	South string `json:"south,omitempty"`
}

type LookCommandData struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Exits       ExitsData `json:"exits"`
	Players     []string  `json:"players"`
	Items       []string  `json:"items"`
	Npcs        []string  `json:"npcs"`
}

type StatusCommandData struct {
	Hp     int    `json:"hp"`
	MaxHp  int    `json:"max_hp"`
	Status string `json:"status"`
}

type CombatPersonData struct {
	Name      string   `json:"name"`
	Hp        int      `json:"hp"`
	Inventory []string `json:"inventory"`
}

type CombatStatsCommandData struct {
	Leader      string                      `json:"leader"`
	CurrentTurn string                      `json:"current_turn"`
	Team        map[string]CombatPersonData `json:"team"`
	Opponents   map[string]CombatPersonData `json:"opponents"`
}

type QuestData struct {
	Id          string `json:"quest_id"`
	Status      string `json:"status"`
	Description string `json:"description"`

	Reward string `json:"reward"`
}

type TrackedQuestData struct {
	Id       string `json:"quest_id"`
	Status   string `json:"status"`
	Progress string `json:"progress"`
}

type AttackCommandData struct {
	AttackerHp int    `json:"attacker_hp"`
	TargetHp   int    `json:"target_hp"`
	Damage     int    `json:"damage"`
	Status     string `json:"status"`
}

type InspectPlayerData struct {
	Name      string `json:"name"`
	IsInGroup bool   `json:"is_in_group"`
	InCombat  bool   `json:"in_combat"`
	Hp        int    `json:"hp"`
	MaxHp     int    `json:"max_hp"`
	Status    string `json:"status"`
}

type InspectNPCData struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Dialogue    []string `json:"dialogue"`
	Role        string   `json:"role"`
	QuestId     string   `json:"quest_id"`
	Hostile     bool     `json:"hostile"`
	Damage      int      `json:"damage"`
	Hp          int      `json:"hp"`
	HpMax       int      `json:"hp_max"`
	InCombat    bool     `json:"in_combat,omitempty"`
	XpReward    int      `json:"xp_reward,omitempty"`
	ItemsReward []string `json:"items_reward,omitempty"`
}

type InspectItemData struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Obtainable  bool   `json:"obtainable"`
	Type        string `json:"type"`
	Tradable    bool   `json:"tradable"`
	Worth       int    `json:"worth"`

	Damage int `json:"damage,omitempty"`

	TargetStat string `json:"target_stat,omitempty"`
	Amount     int    `json:"amount,omitempty"`
	EffectType string `json:"effect_type,omitempty"`
	Duration   int    `json:"duration,omitempty"`
	Charges    int    `json:"charges,omitempty"`
}
