package panel

import (
	"github.com/rivo/tview"
)

type InfoComponent struct {
	Layout *tview.Flex
	View   *tview.TextView
}

func NewInfoComponent(app *tview.Application) *InfoComponent {
	src := &InfoComponent{}

	src.View = createTextView("Select an entity or item to view details.", " Details ", true, Default, Black)
	src.View.SetDynamicColors(true).SetWordWrap(true)

	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(src.View, 0, 1, false)

	return src
}

func (c *InfoComponent) ListenOutputs(app *tview.Application, infoChan <-chan string) {
	go func() {
		for msg := range infoChan {
			app.QueueUpdateDraw(func() {
				c.View.SetText(msg)
			})
		}
	}()
}
