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
