package panel

import (
	pr "tap/protocol"

	"github.com/rivo/tview"
)

func NewNavigationComponent(
	app *tview.Application,
	popupGrid *tview.Grid,
	room_name string,
	exits map[string]string,
	actionsChan chan<- Action,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
) *ChoiceListComponent {
	options := ConvertExits(room_name, exits, actionsChan)

	src := NewChoiceListComponent(app, popupGrid, "Rooms", options, onOpenPopup, onClosePopup, false)

	return src
}

func ConvertExits(room_name string, exits map[string]string, actionsChan chan<- Action) map[string]OptionsMap {
	res := make(OptionsMap)

	for exit, room := range exits {
		dir := exit
		res[exit+": "+room] = map[string]func(){
			pr.CmdMove: func() {
				actionsChan <- Action{
					Type:    ActionSendServer,
					Payload: pr.CmdMove + " " + dir,
				}
				actionsChan <- Action{
					Type:    ActionSendServer,
					Payload: pr.CmdLook,
				}
			},
		}
	}

	return map[string]OptionsMap{room_name: res}
}
