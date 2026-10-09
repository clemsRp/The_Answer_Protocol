package ui

import (
	vars "tap/client/gui/src/variables"

	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Interaction struct {
	ID      string
	Emote   *Emote
	Buttons []*Button

	hovered bool

	Inspect func(item string)
	Name    string
}

func (i *Interaction) Rect() rl.Rectangle {
	if i.Emote == nil || len(i.Emote.Frames) == 0 {
		return rl.NewRectangle(0, 0, 0, 0)
	}

	frame := i.Emote.Frames[0].Frame
	w := float32(vars.FRAME_WIDTH) * i.Emote.Zoom * frame.RatioX
	h := float32(vars.FRAME_HEIGHT) * i.Emote.Zoom * frame.RatioY

	return rl.NewRectangle(i.Emote.X, i.Emote.Y, w, h)
}

func (i *Interaction) IsHovered() bool {
	return i.hovered
}

func (m *Manager) SetViewInteractions(view string, items []*Interaction) {
	m.views_interactions[view] = items
}

func (m *Manager) Interactions(view string) []*Interaction {
	return m.views_interactions[view]
}

func (m *Manager) UpdateInteractions(view string) {
	items := m.views_interactions[view]
	if len(items) == 0 {
		return
	}

	mouse := rl.GetMousePosition()
	clicked := rl.IsMouseButtonPressed(rl.MouseButtonLeft)
	down := rl.IsMouseButtonDown(rl.MouseButtonLeft)

	for _, item := range items {
		// Check emote hover
		hover := rl.CheckCollisionPointRec(mouse, item.Rect())
		near := m.IsPlayerNear(item.Emote.X, item.Emote.Y)
		item.hovered = hover && near

		if hover && down {
			item.Inspect(item.ID)
		}

		button_hovered := false

		// Check buttons interactions
		for _, b := range item.Buttons {
			b.hovered = rl.CheckCollisionPointRec(mouse, b.Rect()) && near
			b.pressed = b.hovered && down

			if b.hovered {
				button_hovered = true
				if clicked && b.OnClick != nil {
					b.OnClick()
				}
			}
		}

		// Keep item hovered
		item.hovered = item.hovered || button_hovered

		// Reset buttons state
		if !item.hovered {
			for _, b := range item.Buttons {
				b.hovered = false
				b.pressed = false
			}
		}
	}
}

func (m *Manager) IsPlayerNear(itemX, itemY float32) bool {
	side_x := math.Pow(float64(itemX-m.PlayerPos.X), 2)
	side_y := math.Pow(float64(itemY-m.PlayerPos.Y), 2)

	dist := math.Sqrt(side_x + side_y)

	return dist <= float64(4*m.Tile)
}
