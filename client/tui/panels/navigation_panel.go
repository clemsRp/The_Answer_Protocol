package panel

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type NavigationComponent struct {
	Layout     *tview.Flex
	Navigation *tview.List
}

func NewNavigationComponent(app *tview.Application) *NavigationComponent {
	src := &NavigationComponent{}

	src.Navigation = createListView(" Rooms ", true, Default, Default, tcell.ColorBlue, Black, false, true)
	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(src.Navigation, 0, 1, false)

	exits := map[string]string{
		"north": "health_aisle",
		"east":  "fresh_section",
	}

	index := 1
	for dir, room := range exits {
		src.Navigation.
			AddItem(dir+" : "+room, "", rune('0'+index), nil)
		index++
	}
	return src
}

func (c *NavigationComponent) ListenOutputs(app *tview.Application, navChan <-chan string, inputs chan<- string) {
	go func() {
		for msg := range navChan {
			app.QueueUpdateDraw(func() {
				c.Navigation.Clear()

				if strings.HasPrefix(msg, "ERROR:") {
					c.Navigation.AddItem("Failed to receive data for navigation.", "", 0, nil)

					c.Navigation.AddItem("↻ Actualiser", "", 'r', func() {
						inputs <- "CMD:REFRESH_ROOM"
					})
					return
				}

				c.Navigation.
					AddItem("NORTH : room3", "", 'n', func() { inputs <- "MOVE:NORTH" }).
					AddItem("EAST : room1", "", 'e', func() { inputs <- "MOVE:EAST" })
			})
		}
	}()
}
