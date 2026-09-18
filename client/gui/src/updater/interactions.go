package updater

import (
	"tap/client/gui/src/ui"
)

var nb_click = 0

func (up *Updater) buildGameInteractions() {
	up.app.Manager.SetViewInteractions("Game", []*ui.Interaction{})
}
