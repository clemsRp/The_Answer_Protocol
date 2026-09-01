package panel

import (
	"context"
	"fmt"
	"strings"
	"sync"
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

	src.History = createTextView("", " Server Responses ", true)

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

func (c *ServerResponseComponent) AppendResponse(res pr.ServerResponse) {
	msg_type := res.Msg
	msg_text := ""

	if strings.ContainsRune(res.Msg, ' ') {
		split_msg := strings.SplitN(res.Msg, " ", 2)
		if len(split_msg) > 0 {
			msg_type = split_msg[0]
		}
		if len(split_msg) > 1 {
			msg_text = split_msg[1]
		}
	}

	color := GetResponseColor(res)
	fmt.Fprintf(c.History, "[%s]%s [white]%s\n", color, msg_type, msg_text)
}

func (c *ServerResponseComponent) ListenOutputs(ctx context.Context, wg *sync.WaitGroup, app *tview.Application, ServerChan <-chan pr.ServerResponse) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return

			case res, ok := <-ServerChan:
				if !ok {
					return
				}

				msg_type := res.Msg
				msg_text := ""

				if strings.ContainsRune(res.Msg, ' ') {
					split_msg := strings.SplitN(res.Msg, " ", 2)
					if len(split_msg) > 0 {
						msg_type = split_msg[0]
					}
					if len(split_msg) > 1 {
						msg_text = split_msg[1]
					}
				}

				color := GetResponseColor(res)
				app.QueueUpdateDraw(func() {
					fmt.Fprintf(c.History, "[%s]%s [white]%s\n", color, msg_type, msg_text)
				})
			}
		}
	}()
}
