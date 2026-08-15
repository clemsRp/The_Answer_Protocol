package panel

import (
	"github.com/rivo/tview"
)

func NewNavigationComponent(
	app *tview.Application,
	popupGrid *tview.Grid,
	exits map[string]string,
	inputs chan<- string,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
) *ChoiceListComponent {
	options := ConvertExits(exits, inputs)

	src := NewChoiceListComponent(app, popupGrid, "Rooms", options, onOpenPopup, onClosePopup, false)

	return src
}

func ConvertExits(exits map[string]string, inputs chan<- string) OptionsMap {
	res := make(OptionsMap)

	for exit, room := range exits {
		res[exit+": "+room] = map[string]func(){
			"MOVE": func() {
				go func() {
					inputs <- "MOVE " + exit
					inputs <- "LOOK"
				}()
			},
		}
	}

	return res
}
