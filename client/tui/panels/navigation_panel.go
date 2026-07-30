package panel

import (
	"strings"
	pr "tap/protocol"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type NavigationComponent struct {
	Layout     *tview.Flex
	Navigation *tview.List
	onSelect   func()
}

func NewNavigationComponent(app *tview.Application, onSelect func()) *NavigationComponent {
	src := &NavigationComponent{
		onSelect: onSelect,
	}

	src.Navigation = createListView(" Actions ", true, Default, Default, tcell.ColorBlue, Black, false, true)
	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(src.Navigation, 0, 1, false)

	exits := map[string]string{
		"north": "health_aisle",
		"east":  "fresh_section",
	}

	index := 1
	for dir, room := range exits {
		action := func() {
			if onSelect != nil {
				go onSelect()
			}
		}

		src.Navigation.
			AddItem(dir+" : "+room, "", rune('0'+index), action)
		index++
	}
	return src
}

func (c *NavigationComponent) ListenOutputs(app *tview.Application, navChan <-chan pr.ServerResponse, inputs chan<- string) {
	go func() {
		for msg := range navChan {
			app.QueueUpdateDraw(func() {
				c.Navigation.Clear()

				if strings.HasPrefix(msg, "ERROR:") {
					return
				}
			})
		}
	}()
}
