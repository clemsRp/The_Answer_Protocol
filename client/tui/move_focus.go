package main

import "github.com/gdamore/tcell/v2"

func (m *MyApp) moveFocusSpatial(dRow, dCol int) {
	currentFocus := m.app.GetFocus()
	startRow, startCol := -1, -1

	for r := range 4 {
		for c := range 4 {
			if m.navMatrix[r][c] == currentFocus {
				startRow, startCol = r, c
				break
			}
		}
	}

	if startRow == -1 {
		return
	}

	targetRow, targetCol := startRow, startCol

	for {
		targetRow += dRow
		targetCol += dCol

		if targetRow < 0 || targetRow >= 4 || targetCol < 0 || targetCol >= 4 {
			return
		}

		nextComponent := m.navMatrix[targetRow][targetCol]

		if nextComponent != currentFocus && nextComponent != nil {
			m.app.SetFocus(nextComponent)
			return
		}
	}
}

func (m *MyApp) SetupFocusManager() {

	m.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == tcell.ModShift {
			switch event.Key() {

			case tcell.KeyUp:
				m.moveFocusSpatial(-1, 0)
				return nil

			case tcell.KeyDown:
				m.moveFocusSpatial(1, 0)
				return nil

			case tcell.KeyLeft:
				m.moveFocusSpatial(0, -1)
				return nil

			case tcell.KeyRight:
				m.moveFocusSpatial(0, 1)
				return nil
			}
		}

		return event
	})
}
