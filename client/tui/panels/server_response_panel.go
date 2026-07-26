package panel

import (
	"fmt"

	"github.com/rivo/tview"
)

type ServerResponseComponent struct {
	Layout  *tview.Flex
	CliBtn  *tview.Button
	History *tview.TextView
}

func NewServerResponseComponent(app *tview.Application) *ServerResponseComponent {
	src := &ServerResponseComponent{}

	src.CliBtn = tview.NewButton("CLI")

	src.History = createTextView("", " Server Responses ", true, Default, Black)

	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(src.CliBtn, 1, 1, false).
		AddItem(src.History, 0, 1, false)

	return src
}

func (c *ServerResponseComponent) ListenOutputs(app *tview.Application, ServerChan <-chan string) {
	// TO CHANGE

	go func() {
		for msg := range ServerChan {

			app.QueueUpdateDraw(func() {
				fmt.Fprint(c.History, msg+"\n")
			})
		}
	}()
}
