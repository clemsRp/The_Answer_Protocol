package panel

import (
	"context"
	"sync"
	pr "tap/protocol"

	"github.com/rivo/tview"
)

type InspectorComponent struct {
	Layout  *tview.Flex
	View    *tview.TextView
	SelfBtn *tview.Button
}

func NewInspectorComponent(app *tview.Application, actionsChan chan<- Action) *InspectorComponent {
	src := &InspectorComponent{}

	src.View = createTextView("", " Inspector ", true)
	src.View.SetDynamicColors(true).SetWordWrap(true)

	src.SelfBtn = tview.NewButton("Self").
		SetSelectedFunc(func() {
			actionsChan <- Action{
				Type:    ActionSendServer,
				Payload: pr.CmdInspectSelf,
			}
		})

	btnRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(src.SelfBtn, 10, 0, false).
		AddItem(nil, 1, 0, false)

	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(btnRow, 1, 0, false).
		AddItem(src.View, 0, 1, false)

	return src
}

func (c *InspectorComponent) SetDatas(text string) {
	c.View.Clear()
	c.View.SetText(text)
}

func (c *InspectorComponent) ListenOutputs(ctx context.Context, wg *sync.WaitGroup, app *tview.Application, datasChan <-chan pr.ServerResponse) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return

			case res, ok := <-datasChan:
				if !ok {
					return
				}

				app.QueueUpdateDraw(func() {
					c.View.Clear()
					c.View.SetText(res.Msg)
				})
			}
		}
	}()
}
