package panel

import (
	"fmt"

	"github.com/rivo/tview"
)

func NewItemsComponent(
	app *tview.Application,
	popupGrid *tview.Grid,
	inputs chan<- string,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
) *ChoiceListComponent {
	room_items := []string{"sword"}
	inventory_items := []string{}

	options := ConvertItems(room_items, inventory_items, inputs)

	return NewChoiceListComponent(app, popupGrid, "Items", options, onOpenPopup, onClosePopup)
}

func ConvertItems(room_items, inventory_items []string, inputs chan<- string) map[string]OptionsMap {
	res := make(map[string]OptionsMap)

	if len(room_items) != 0 {
		res["ROOM"] = ConvertItemsList(room_items, "TAKE", inputs)
	}

	if len(inventory_items) != 0 {
		res["INVENTORY"] = ConvertItemsList(inventory_items, "DROP", inputs)
	}

	return res
}

func ConvertItemsList(items []string, cmd string, inputs chan<- string) OptionsMap {
	res := make(OptionsMap)

	for _, item := range items {
		res[item] = map[string]func(){
			cmd: func() {
				go func() {
					inputs <- fmt.Sprintf("%s %s", cmd, item)
				}()
			},
		}
	}

	return res
}
