package panel

import (
	"fmt"
	pr "tap/protocol"

	"github.com/rivo/tview"
)

var (
	inventory = make([]string, 0)
)

func SetInventory(inv []string) {
	inventory = inv
}

func GetInventory() []string {
	return inventory
}

func NewItemsComponent(
	app *tview.Application,
	popupGrid *tview.Grid,
	room_items []string,
	inventory_items []string,
	actionsChan chan<- Action,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
) *ChoiceListComponent {
	if inventory_items != nil {
		inventory = inventory_items
	}
	options := ConvertItems(room_items, inventory, actionsChan)

	return NewChoiceListComponent(app, popupGrid, "Items", options, onOpenPopup, onClosePopup, false)
}

func ConvertItems(room_items, inventory_items []string, actionsChan chan<- Action) map[string]OptionsMap {
	res := make(map[string]OptionsMap)

	if len(room_items) != 0 {
		res["ROOM"] = ConvertRoomItemsList(room_items, actionsChan)
	}

	if len(inventory_items) != 0 {
		res["INVENTORY"] = ConvertInventoryItemsList(inventory_items, actionsChan)
	}

	return res
}

func ConvertRoomItemsList(items []string, actionsChan chan<- Action) OptionsMap {
	res := make(OptionsMap)

	for _, item := range items {
		it := item
		res[it] = map[string]func(){
			pr.CmdTake: func() {
				actionsChan <- Action{
					Type:    ActionSendServer,
					Payload: fmt.Sprintf("%s %s", pr.CmdTake, it),
				}
			},
		}
	}

	return res
}

func ConvertInventoryItemsList(items []string, actionsChan chan<- Action) OptionsMap {
	res := make(OptionsMap)

	for _, item := range items {
		it := item
		res[it] = map[string]func(){
			pr.CmdDrop: func() {
				actionsChan <- Action{
					Type:    ActionSendServer,
					Payload: fmt.Sprintf("%s %s", pr.CmdDrop, it),
				}
			},
		}
	}

	return res
}
