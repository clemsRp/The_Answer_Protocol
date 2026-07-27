package panel

import (
	"strings"

	"github.com/rivo/tview"
)

type DatasComponent struct {
	Layout *tview.Flex
	View   *tview.TextView
}

func NewDatasComponent(app *tview.Application) *DatasComponent {
	src := &DatasComponent{}

	src.View = createTextView("Display of datas will appear here", " Datas ", true, Default, Black)
	src.View.SetDynamicColors(true).SetWordWrap(true)

	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(src.View, 0, 1, false)
	src.View.SetText("Display of datas will appear here")

	return src
}

func (c *DatasComponent) ListenOutputs(app *tview.Application, datasChan <-chan string) {
	go func() {
		for msg := range datasChan {
			app.QueueUpdateDraw(func() {
				c.View.Clear()

				if strings.HasPrefix(msg, "ERROR:") {
					c.View.SetText("[red]Error loading dialogue.")
					return
				}

				c.View.SetText(msg)
			})
		}
	}()
}
