package ui

type Frame struct {
	IndX, IndY     float32
	RatioX, RatioY float32
}

type Manager struct {
	views_buttons      map[string][]*Button
	views_emotes       map[string][]*Emote
	views_interactions map[string][]*Interaction
}

func NewManager() *Manager {
	return &Manager{
		views_buttons:      map[string][]*Button{},
		views_emotes:       map[string][]*Emote{},
		views_interactions: map[string][]*Interaction{},
	}
}

func (m *Manager) SetViewButtons(view string, buttons []*Button) {
	m.views_buttons[view] = buttons
}

func (m *Manager) SetViewEmotes(view string, emotes []*Emote) {
	m.views_emotes[view] = emotes
}

func (m *Manager) SetViewInteractions(view string, items []*Interaction) {
	m.views_interactions[view] = items
}

func (m *Manager) Update(view string) {
	m.UpdateButtons(view)
	m.UpdateEmotes(view)
	m.UpdateInteractions(view)
}
