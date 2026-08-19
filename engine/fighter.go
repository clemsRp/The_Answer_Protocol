package engine

const (
	StatusNormal    string = "normal"
	StatusParalysed string = "paralysed"
	StatusSleeping  string = "sleeping"
	StatusSleepy    string = "sleepy"
)

type CombatStats struct {
	Hp         int    `json:"hp" validate:"gte=0"`
	HpMax      int    `json:"max_hp" validate:"gt=0,gtefield=Hp"`
	Mana       int    `json:"mana" validate:"gte=0"`
	Initiative int    `json:"initiative" validate:"gt=0"`
	CombatId   string `json:"combat_id,omitempty"`
	Status     string `json:"status" validate:"required,valid_npc_status"`
}

func (c *CombatStats) Clone() *CombatStats {
	if c == nil {
		return nil
	}
	cp := *c
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
