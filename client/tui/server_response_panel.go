package main

import (
	"fmt"

	"github.com/rivo/tview"
)

type ServerResponseComponent struct {
	Layout  *tview.Flex
	History *tview.TextView
}

func NewServerResponseComponent(app *tview.Application) *ServerResponseComponent {
	src := &ServerResponseComponent{}

	src.History = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true)

	src.History.SetBorder(true).SetTitle(" Server Responses ")

	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(src.History, 0, 1, false)

	return src
}

func (c *ServerResponseComponent) ListenOutputs(app *tview.Application, outputs <-chan string) {
	go func() {
		for msg := range outputs {

			lineChat := fmt.Sprintf("[gray][%s] [yellow][%s] [green]%s: [white]%s\n",
				msg, msg, msg, msg)

			app.QueueUpdateDraw(func() {
				fmt.Fprint(c.History, lineChat)
			})
		}
	}()
}
