package variables

type GroupPanel struct {
	Open          bool
	InGroup       bool
	Leader        string
	Grouped       []string
	UnGrouped     []string
	Invitations   []string
	SendPromotion bool
	Promote       bool
}
