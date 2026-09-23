package variables

type GroupPanel struct {
	Open            bool
	InGroup         bool
	IsLeader        bool
	Leader          string
	Grouped         []string
	UnGrouped       []string
	Invitations     []string
	SendInvitations []string
	SendPromotion   string
	Promote         bool
}
