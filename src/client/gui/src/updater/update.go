package updater

import (
	"tap/src/client/gui/src/core"
	panel "tap/src/client/tui/panels"
)

type Updater struct {
	app         *core.App
	actionsChan chan panel.Action
}

func NewUpdater(app *core.App) *Updater {
	return &Updater{
		app:         app,
		actionsChan: app.ActionsChan,
	}
}
