package ui

import (
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Button struct {
	ID      string
	Texture string

	X, Y     float32
	Zoom     float32
	Rotation float32

	Normal  Frame
	Hover   Frame
	Pressed Frame

	OnClick func()

	hovered bool
	pressed bool
}

func (b *Button) Rect() rl.Rectangle {
	w := float32(vars.FRAME_WIDTH) * b.Zoom * b.Normal.RatioX
	h := float32(vars.FRAME_HEIGHT) * b.Zoom * b.Normal.RatioY
	return rl.NewRectangle(b.X, b.Y, w, h)
}

func (b *Button) CurrentFrame() Frame {
	if b.pressed && b.Pressed.RatioX != 0 {
		return b.Pressed
	}

	return b.Normal
}

func (m *Manager) SetViewButtons(view string, buttons []*Button) {
	m.views_buttons[view] = buttons
}

func (m *Manager) Buttons(view string) []*Button {
	return m.views_buttons[view]
}

func (m *Manager) UpdateButtons(view string) {
	buttons := m.views_buttons[view]
	if len(buttons) == 0 {
		return
	}

	mouse := rl.GetMousePosition()
	clicked := rl.IsMouseButtonPressed(rl.MouseButtonLeft)
	down := rl.IsMouseButtonDown(rl.MouseButtonLeft)

	for _, b := range buttons {
		b.hovered = rl.CheckCollisionPointRec(mouse, b.Rect())
		b.pressed = b.hovered && down

		if b.hovered && clicked && b.OnClick != nil {
			b.OnClick()
		}
	}
}
