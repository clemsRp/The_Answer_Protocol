package panel

import (
	"fmt"
	pr "tap/protocol"

	"github.com/rivo/tview"
)

type ServerResponseComponent struct {
	Layout  *tview.Flex
	CliBtn  *tview.Button
	QuitBtn *tview.Button
	Buttons *tview.Flex
	History *tview.TextView
}

func NewServerResponseComponent(app *tview.Application) *ServerResponseComponent {
	src := &ServerResponseComponent{}

	src.CliBtn = tview.NewButton("CLI")
	src.QuitBtn = tview.NewButton("QUIT")

	src.History = createTextView("", " Server Responses ", true, Default, Black)

	src.Buttons = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(nil, 1, 0, false).
		AddItem(src.CliBtn, 0, 1, false).
		AddItem(nil, 1, 0, false).
		AddItem(src.QuitBtn, 0, 1, false).
		AddItem(nil, 1, 0, false)

	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(src.Buttons, 1, 1, false).
		AddItem(src.History, 0, 1, false)

	return src
}

func (c *ServerResponseComponent) ListenOutputs(app *tview.Application, ServerChan <-chan pr.ServerResponse) {
	// TO CHANGE

	go func() {
		for res := range ServerChan {
			color := GetResponseColor(res)
			app.QueueUpdateDraw(func() {
				fmt.Fprintf(c.History, "[%s] %s\n", color, res.Msg)
			})
		}
	}()
}
