package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ItemsComponent struct {
	Layout *tview.Flex
	Items  *tview.List
}

func NewItemsRoomComponent(app *tview.Application) *ItemsComponent {
	src := &ItemsComponent{}

	src.Items = createListView(" Items in room ", true, Default, Default, tcell.ColorGreen, Black, false, true)

	// src.Items.AddItem("Looking for items in room", "", 0, nil)
	src.Items.AddItem("Ramasser : Épée", "", '1', nil).
		AddItem("Ramasser : Potion", "", '2', nil)
	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(src.Items, 0, 1, false)

	return src
}

func (c *ItemsComponent) ListenOutputs(app *tview.Application, itemsChan <-chan string, inputs chan<- string) {
	go func() {
		for msg := range itemsChan {
			app.QueueUpdateDraw(func() {
				c.Items.Clear()

				if strings.HasPrefix(msg, "ERROR:") {
					c.Items.AddItem("Erreur lors du chargement.", "", 0, nil)
					c.Items.AddItem("↻ Actualiser", "", 'r', func() {
						inputs <- "CMD:REFRESH_ITEMS"
					})
					return
				}

				c.Items.
					AddItem("Ramasser : Épée", "", '1', func() { inputs <- "TAKE:SWORD" }).
					AddItem("Ramasser : Potion", "", '2', func() { inputs <- "TAKE:POTION" })
			})
		}
	}()
}
