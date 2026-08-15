package panel

import (
	"fmt"
	pr "tap/protocol"

	"github.com/rivo/tview"
)

var (
	inventory = make([]string, 0)
)

func NewItemsComponent(
	app *tview.Application,
	popupGrid *tview.Grid,
	room_items []string,
	inputs chan<- string,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
) *ChoiceListComponent {
	options := ConvertItems(room_items, inventory, inputs)

	return NewChoiceListComponent(app, popupGrid, "Items", options, onOpenPopup, onClosePopup, false)
}

func ConvertItems(room_items, inventory_items []string, inputs chan<- string) map[string]OptionsMap {
	res := make(map[string]OptionsMap)

	if len(room_items) != 0 {
		res["ROOM"] = ConvertItemsList(room_items, pr.CmdTake, inputs)
	}

	if len(inventory_items) != 0 {
		res["INVENTORY"] = ConvertItemsList(inventory_items, pr.CmdDrop, inputs)
	}

	return res
}

func ConvertItemsList(items []string, cmd string, inputs chan<- string) OptionsMap {
	res := make(OptionsMap)

	for _, item := range items {
		res[item] = map[string]func(){
			cmd: func() {
				go func() {
					ItemFunc(item, cmd, inputs)
					inputs <- "LOOK"
				}()
			},
		}
	}

	return res
}

func ItemFunc(item, cmd string, inputs chan<- string) {
	if cmd == pr.CmdTake {
		inventory = append(inventory, item)

	} else if cmd == pr.CmdDrop {
		// Get cur_item index inside inventory
		item_index := -1
		for i, cur_item := range inventory {
			if cur_item == item {
				item_index = i
				break
			}
		}

		// Remove cur_item
		if item_index != -1 {
			inventory = append(inventory[:item_index], inventory[item_index+1:]...)
		}
	}
	inputs <- fmt.Sprintf("%s %s", cmd, item)
}
