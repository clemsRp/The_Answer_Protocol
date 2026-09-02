package panel

import (
	"context"
	"sync"
	pr "tap/protocol"

	"github.com/rivo/tview"
)

type InspectorComponent struct {
	Layout *tview.Flex
	View   *tview.TextView
}

func NewInspectorComponent(app *tview.Application) *InspectorComponent {
	src := &InspectorComponent{}

	src.View = createTextView("", " Inspector ", true)
	src.View.SetDynamicColors(true).SetWordWrap(true)

	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(src.View, 0, 1, false)

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
