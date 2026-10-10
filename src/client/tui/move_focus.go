package tui

import "github.com/gdamore/tcell/v2"

func (m *MyApp) moveFocusSpatial(dRow, dCol int) {
	currentFocus := m.app.GetFocus()
	if currentFocus == nil {
		return
	}
	startRow, startCol := -1, -1

	numRows := len(m.navMatrix)
	if numRows == 0 {
		return
	}

	numCols := 1
	if len(m.navMatrix) > 0 {
		numCols = len(m.navMatrix[0])
	}

	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
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

		if targetRow < 0 || targetRow >= numRows || targetCol < 0 || targetCol >= numCols {
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
