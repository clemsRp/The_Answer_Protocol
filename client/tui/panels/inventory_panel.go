package panel

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InventoryComponent struct {
	Layout *tview.Flex
	List   *tview.List
}

func NewInventoryComponent(app *tview.Application) *InventoryComponent {
	src := &InventoryComponent{}

	src.List = createListView(" Inventory ", true, Default, Default, tcell.ColorPurple, Black, false, true)

	src.List.AddItem("Loading inventory...", "", 0, nil)

	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(src.List, 0, 1, false)

	return src
}

func (c *InventoryComponent) ListenOutputs(app *tview.Application, invChan <-chan string, inputs chan<- string) {
	go func() {
		for msg := range invChan {
			app.QueueUpdateDraw(func() {
				c.List.Clear()

				if strings.HasPrefix(msg, "ERROR:") {
					c.List.AddItem("Could not load inventory.", "", 0, nil)
					return
				}

				items := strings.Split(msg, ",")
				for i, item := range items {
					// Use index as a shortcut
					shortcut := rune('1' + i)
					itemName := item

					c.List.AddItem(itemName, "", shortcut, func() {
						// Send action to server
						inputs <- "USE_ITEM:" + itemName
					})
				}
			})
		}
	}()
}
