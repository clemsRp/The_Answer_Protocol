package drawer

import vars "tap/client/gui/src/variables"

func (dr *Drawer) DrawLeftPanel() {
	dr.DrawWoodFrameAt(
		vars.Position{X: 1, Y: 4},
		vars.Position{X: 7, Y: 17},
		1, 0, false,
	)
}
