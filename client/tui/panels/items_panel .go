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

	switch opts := options.(type) {
	case OptionsMap:
		return NewChoiceListComponent(app, popupGrid, "Items", opts, onOpenPopup, onClosePopup)
	case map[string]OptionsMap:
		return NewChoiceListComponent(app, popupGrid, "Items", opts, onOpenPopup, onClosePopup)
	default:
		return nil
	}
}

func ConvertItems(room_items, inventory_items []string, inputs chan<- string) any {
	if len(room_items) == 0 {
		return ConvertItemsList(inventory_items, "DROP", inputs)

	} else if len(inventory_items) == 0 {
		return ConvertItemsList(room_items, "TAKE", inputs)
	}

	return map[string]OptionsMap{
		"ROOM":      ConvertItemsList(room_items, "TAKE", inputs),
		"INVENTORY": ConvertItemsList(inventory_items, "DROP", inputs),
	}
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
