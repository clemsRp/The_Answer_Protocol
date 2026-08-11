package panel

import (
	pr "tap/protocol"

	"github.com/rivo/tview"
)

type DatasComponent struct {
	Layout *tview.Flex
	View   *tview.TextView
}

func NewDatasComponent(app *tview.Application) *DatasComponent {
	src := &DatasComponent{}

	src.View = createTextView("", " Datas ", true)
	src.View.SetDynamicColors(true).SetWordWrap(true)

	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(src.View, 0, 1, false)

	return src
}

func (c *DatasComponent) ListenOutputs(app *tview.Application, datasChan <-chan pr.ServerResponse) {
	go func() {
		for res := range datasChan {
			app.QueueUpdateDraw(func() {
				c.View.Clear()

				c.View.SetText(res.Msg)
			})
		}
	}()
}
