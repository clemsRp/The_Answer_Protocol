package drawer

import (
	vars "tap/client/gui/src/variables"
)

func (dr *Drawer) DrawCombatView() {
	dr.DrawBlurredGame()

	dr.DrawCombatChatZone()
	dr.DrawCombatFightersZone()
	dr.DrawCombatActionsZone()

	dr.DrawPlayerPanel()
	dr.DrawGameEmotes()
	dr.DrawInventoryEmotes()
	dr.DrawCombatInventoryButtons()
}

func (dr *Drawer) DrawCombatFrame(startX, startY, endX, endY float32) {
	dr.DrawWoodFrameAt(
		vars.Position{X: startX, Y: startY},
		vars.Position{X: endX, Y: endY},
		1, 0, false,
	)
}
