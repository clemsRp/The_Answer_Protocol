package panel

import (
	"github.com/rivo/tview"
)

func NewNavigationComponent(
	app *tview.Application,
	popupGrid *tview.Grid,
	room_name string,
	exits map[string]string,
	inputs chan<- string,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
) *ChoiceListComponent {
	options := ConvertExits(room_name, exits, inputs)

	src := NewChoiceListComponent(app, popupGrid, "Rooms", options, onOpenPopup, onClosePopup, false)

	return src
}

func ConvertExits(room_name string, exits map[string]string, inputs chan<- string) map[string]OptionsMap {
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

	return map[string]OptionsMap{room_name: res}
}
