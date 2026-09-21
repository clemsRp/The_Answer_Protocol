package ui

import (
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Select struct {
	ID      string
	Texture string

	X, Y     float32
	Zoom     float32
	Rotation float32

	Normal Frame
	Hover  Frame

	ItemNormal Frame
	ItemHover  Frame

	Options       []string
	CurrentOption string
	Open          bool

	OnChange func(option string)

	hovered      bool
	hoveredIndex int
}

func (s *Select) MainRect() rl.Rectangle {
	// Get main rectangle
	w := float32(vars.FRAME_WIDTH) * s.Zoom * s.Normal.RatioX
	h := float32(vars.FRAME_HEIGHT) * s.Zoom * s.Normal.RatioY

	return rl.NewRectangle(s.X, s.Y, w, h)
}

func (s *Select) ItemRect(index int) rl.Rectangle {
	// Get item rectangle
	w := float32(vars.FRAME_WIDTH) * s.Zoom * s.ItemNormal.RatioX
	h := float32(vars.FRAME_HEIGHT) * s.Zoom * s.ItemNormal.RatioY
	mainH := float32(vars.FRAME_HEIGHT) * s.Zoom * s.Normal.RatioY

	y := s.Y + mainH + (float32(index) * h)

	return rl.NewRectangle(s.X, y, w, h)
}

func (s *Select) CurrentMainFrame() Frame {
	// Get current frame
	if s.hovered || s.Open {
		return s.Hover
	}

	return s.Normal
}

func (m *Manager) SetViewSelects(view string, selects []*Select) {
	m.views_selects[view] = selects
}

func (m *Manager) Selects(view string) []*Select {
	return m.views_selects[view]
}

func (m *Manager) UpdateSelects(view string) {
	selects := m.views_selects[view]
	if len(selects) == 0 {
		return
	}

	mouse := rl.GetMousePosition()
	clicked := rl.IsMouseButtonPressed(rl.MouseButtonLeft)

	for i := len(selects) - 1; i >= 0; i-- {
		s := selects[i]

		// Check main collision
		s.hovered = rl.CheckCollisionPointRec(mouse, s.MainRect())
		s.hoveredIndex = -1

		if s.Open {
			// Check items collisions
			for j := range s.Options {
				if rl.CheckCollisionPointRec(mouse, s.ItemRect(j)) {
					s.hoveredIndex = j
					break
				}
			}

			// Handle click
			if clicked {
				if s.hoveredIndex != -1 {
					// Update current option
					s.CurrentOption = s.Options[s.hoveredIndex]
					s.Open = false
					if s.OnChange != nil {
						s.OnChange(s.CurrentOption)
					}
				} else {
					// Close select
					s.Open = false
				}
			}
		} else {
			// Open select
			if s.hovered && clicked {
				s.Open = true
			}
		}
	}
}

func (s *Select) HoveredIndex() int {
	return s.hoveredIndex
}

func (s *Select) ItemFrame(index int) Frame {
	if s.Open && s.hoveredIndex == index {
		return s.ItemHover
	}
	return s.ItemNormal
}
