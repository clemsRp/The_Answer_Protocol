package panel

import (
	"github.com/rivo/tview"
)

func NewInteractionComponent(
	app *tview.Application,
	popupGrid *tview.Grid,
	inputs chan<- string,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
) *ChoiceListComponent {

	options := OptionsMap{
		"Adelina": {
			"TALK": func() {
				go func() {
					inputs <- "TALK Adelina"
				}()
			},
		},
		"Jeng Jong": {
			"TALK": func() {
				go func() {
					inputs <- "TALK Jeng Jong"
				}()
			},
			"ATTACK": func() {
				go func() {
					inputs <- "ATTACK Jeng Jong"
				}()
			},
		},
	}

	src := NewChoiceListComponent(app, popupGrid, "Interactions", options, onOpenPopup, onClosePopup)

	return src
}
