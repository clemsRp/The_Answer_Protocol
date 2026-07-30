package panel

import (
	pr "tap/protocol"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PlayersNPCComponent struct {
	Layout *tview.Flex
	List   *tview.List
}

func NewPlayersNPCComponent(app *tview.Application) *PlayersNPCComponent {
	src := &PlayersNPCComponent{}

	src.List = createListView("NPCs", true, Default, Default, tcell.ColorYellow, Black, false, true)

	src.List.AddItem("Charging entities...", "", 0, nil)

	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(src.List, 0, 1, false)

	return src
}

func (c *PlayersNPCComponent) ListenOutputs(app *tview.Application, playersChan <-chan pr.ServerResponse, inputs chan<- string) {
	go func() {
		for res := range playersChan {
			app.QueueUpdateDraw(func() {
				c.List.Clear()

				if IsErrorResponse(res) {
					//TODO Handle error
					c.List.AddItem(pr.ErrInternalServer, "", 0, nil)
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
