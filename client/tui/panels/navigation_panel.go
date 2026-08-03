package panel

import (
	"github.com/rivo/tview"
)

func NewNavigationComponent(
	app *tview.Application,
	popupGrid *tview.Grid,
	inputs chan<- string,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
) *ChoiceListComponent {

	options := OptionsMap{
		"east: tavern": {
			"MOVE": func() {
				go func() {
					inputs <- "MOVE east"
				}()
			},
		},
		"north: entrance": {
			"MOVE": func() {
				go func() {
					inputs <- "MOVE north"
				}()
			},
			"MOVEE": func() {
				go func() {
					inputs <- "MOVE north"
				}()
			},
			"MOVEEE": func() {
				go func() {
					inputs <- "MOVE north"
				}()
			},
			"MOVEEEE": func() {
				go func() {
					inputs <- "MOVE north"
				}()
			},
		},
	}

	src := NewChoiceListComponent(app, popupGrid, options, onOpenPopup, onClosePopup)

	return src
}
