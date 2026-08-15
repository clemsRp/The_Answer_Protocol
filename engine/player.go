package engine

type Player struct {
	id             string
	name           string
	room           string
	hp             int
	hpMax          int
	group          string
	status         string
	inventory      []string
	quests         []string
	invitation     []string
	promotion      bool
	send_promotion bool
}
