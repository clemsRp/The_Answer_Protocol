package panel

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PlayersNPCComponent struct {
	Layout *tview.Flex
	List   *tview.List
}

func NewPlayersNPCComponent(app *tview.Application) *PlayersNPCComponent {
	src := &PlayersNPCComponent{}

	src.List = createListView(" Players & NPC ", true, Default, Default, tcell.ColorYellow, Black, false, true)

	src.List.AddItem("Charging entities...", "", 0, nil)
	src.List.
		AddItem("Robert [NPC]", "", 0, nil).
		AddItem("Ali [Player]", "", 0, nil)

	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(src.List, 0, 1, false)

	return src
}

func (c *PlayersNPCComponent) ListenOutputs(app *tview.Application, playersChan <-chan string, inputs chan<- string) {
	go func() {
		for msg := range playersChan {
			app.QueueUpdateDraw(func() {
				c.List.Clear()

				if strings.HasPrefix(msg, "ERROR:") {
					c.List.AddItem("Error of synchronisation.", "", 0, nil)
					c.List.AddItem("↻ Refresh", "", 'r', func() {
						inputs <- "CMD:REFRESH_ENTITIES"
					})
					return
				}

				c.List.
					AddItem("Robert [NPC]", "", 0, func() {
						inputs <- "INSPECT:ROBERT"
					}).
					AddItem("Ali [Player]", "", 0, nil)
			})
		}
	}()
}
