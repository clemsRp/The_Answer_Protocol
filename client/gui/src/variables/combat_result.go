package variables

import "time"

type CombatResultPanel struct {
	LastTime time.Time
	Result   string
	Rewards  []string
}
