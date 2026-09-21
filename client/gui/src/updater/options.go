package updater

import "tap/client/gui/src/ui"

func (up *Updater) buildGroupOptions() {
	up.app.Manager.SetViewOptions("Group", []*ui.Option{})
}
