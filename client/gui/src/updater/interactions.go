package updater

import (
	"tap/client/gui/src/ui"
)

func (up *Updater) buildGameInteractions() {
	up.app.Manager.SetViewInteractions("Game", []*ui.Interaction{})
}

func (up *Updater) buildCombatInteractions() {
	up.app.Manager.SetViewInteractions("Combat", []*ui.Interaction{})
}
