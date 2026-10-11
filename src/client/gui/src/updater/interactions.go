package updater

import (
	"tap/src/client/gui/src/ui"
)

func (up *Updater) buildGameInteractions() {
	up.app.Manager.SetViewInteractions("Game", []*ui.Interaction{})
}
