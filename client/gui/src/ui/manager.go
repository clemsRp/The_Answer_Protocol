package ui

type Frame struct {
	IndX, IndY     float32
	RatioX, RatioY float32
}

type Manager struct {
	views_buttons map[string][]*Button
	// views_emotes  map[string][]*Emote
}

func NewManager() *Manager {
	return &Manager{
		views_buttons: map[string][]*Button{},
		// views_emotes:  map[string][]*Emote{},
	}
}

func (m *Manager) SetViewButtons(view string, buttons []*Button) {
	m.views_buttons[view] = buttons
}

func (m *Manager) SetViewEmotes(view string) {
	// m.views_emotes[view] = emotes
}

func (m *Manager) Update(view string) {
	m.UpdateButtons(view)
	// m.UpdateEmotes(view)
}
