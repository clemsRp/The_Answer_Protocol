package panel

import (
	"fmt"
	pr "tap/protocol"

	"github.com/rivo/tview"
)

type DatasComponent struct {
	Layout *tview.Flex
	View   *tview.TextView
}

func NewDatasComponent(app *tview.Application) *DatasComponent {
	src := &DatasComponent{}

	src.View = createTextView("", " Datas ", true, Default, Black)
	src.View.SetDynamicColors(true).SetWordWrap(true)

	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(src.View, 0, 1, false)

	return src
}

func (c *DatasComponent) ListenOutputs(app *tview.Application, datasChan <-chan pr.ServerResponse) {
	go func() {
		for res := range datasChan {
			color := GetResponseColor(res)
			app.QueueUpdateDraw(func() {
				c.View.Clear()

				if IsErrorResponse(res) {
					errRes := fmt.Sprintf("[%s] %s", color, res.Msg)
					c.View.SetText(errRes)
					return
				}
				formatMsg := res.Msg
				// if res.Datas != nil {
				// 	for _, item := range res.Datas {

				// 	}
				// }
				// var out bytes.Buffer
				// json.Indent(&out, []byte(res), "", "    ")
				// formatted := out.String()

				// re := regexp.MustCompile(`\[\s*([^\[\]\{\}]*?)\s*\]`)
				// formatMsg := re.ReplaceAllStringFunc(formatted, func(m string) string {
				// 	parts := strings.Fields(m)
				// 	return strings.Join(parts, " ")
				// })

				c.View.SetText(formatMsg)
			})
		}
	}()
}
