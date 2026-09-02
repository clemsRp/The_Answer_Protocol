package panel

import (
	"context"
	"sync"
	pr "tap/protocol"

	"github.com/rivo/tview"
)

type QuestComponent struct {
	Layout *tview.Flex
	View   *tview.TextView
}

func NewQuestComponent(app *tview.Application) *QuestComponent {
	src := &QuestComponent{}

	src.View = createTextView("", " Quests ", true)
	src.View.SetDynamicColors(true).SetWordWrap(true)

	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(src.View, 0, 1, false)

	return src
}

func (c *QuestComponent) SetDatas(text string) {
	c.View.Clear()
	c.View.SetText(text)
}

func (c *QuestComponent) ListenOutputs(ctx context.Context, wg *sync.WaitGroup, app *tview.Application, datasChan <-chan pr.ServerResponse) {
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
