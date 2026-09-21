package ui

type Frame struct {
	IndX, IndY     float32
	RatioX, RatioY float32
}

type Manager struct {
	views_buttons      map[string][]*Button
	views_emotes       map[string][]*Emote
	views_interactions map[string][]*Interaction
	views_selects      map[string][]*Select
	views_options      map[string][]*Option
}

func NewManager() *Manager {
	return &Manager{
		views_buttons:      map[string][]*Button{},
		views_emotes:       map[string][]*Emote{},
		views_interactions: map[string][]*Interaction{},
		views_selects:      map[string][]*Select{},
		views_options:      map[string][]*Option{},
	}
}

func (m *Manager) Update(view string) {
	m.UpdateButtons(view)
	m.UpdateEmotes(view)
	m.UpdateInteractions(view)
	m.UpdateSelects(view)
	m.UpdateOptions(view)
}
