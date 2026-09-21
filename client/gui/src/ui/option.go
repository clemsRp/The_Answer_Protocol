package ui

import rl "github.com/gen2brain/raylib-go/raylib"

type Option struct {
	OptionName string
	X          int
	Y          int
	HasDecline bool
	DeclineBtn *Button
	AcceptBtn  *Button
}

func (m *Manager) UpdateOptions(view string) {
	mouse := rl.GetMousePosition()
	clicked := rl.IsMouseButtonPressed(rl.MouseButtonLeft)
	down := rl.IsMouseButtonDown(rl.MouseButtonLeft)

	for _, o := range m.views_options[view] {
		for _, b := range []*Button{o.DeclineBtn, o.AcceptBtn} {
			if b == nil {
				continue
			}
			b.hovered = rl.CheckCollisionPointRec(mouse, b.Rect())
			b.pressed = b.hovered && down
			if b.hovered && clicked && b.OnClick != nil {
				b.OnClick()
			}
		}
	}
}

func (m *Manager) SetViewOptions(view string, options []*Option) {
	m.views_options[view] = options
}

func (m *Manager) Options(view string) []*Option {
	return m.views_options[view]
}
