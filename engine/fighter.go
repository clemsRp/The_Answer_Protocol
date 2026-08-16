package engine

const (
	StatusNormal    string = "normal"
	StatusParalysed string = "paralysed"
	StatusSleeping  string = "sleeping"
	StatusSleepy    string = "sleepy"
)

type CombatStats struct {
	Hp         int    `json:"hp"`
	Initiative int    `json:"initiative"`
	HpMax      int    `json:"max_hp"`
	CombatId   string `json:"combat_id,omitempty"`
	Status     string `json:"status,omitempty"`
}

func (cs *CombatStats) Clone() *CombatStats {
	if cs == nil {
		return nil
	}
	cp := *cs
	return &cp
}

type Fighter interface {
	takeDamage(amount int)
	isDead() bool
	getName() string
	getHp() int
	getInitiative() int
	playCombatTurn(target Fighter) *CombatTurnResult
}
