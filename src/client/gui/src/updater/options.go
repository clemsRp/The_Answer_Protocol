package updater

import "tap/src/client/gui/src/ui"

func (up *Updater) buildGroupOptions() {
	up.app.Manager.SetViewOptions("Group", []*ui.Option{})
}
