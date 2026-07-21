package panel

import (
	"strings"

	"github.com/rivo/tview"
)

type DialogueComponent struct {
	Layout *tview.Flex
	View   *tview.TextView
}

func NewDialogueComponent(app *tview.Application) *DialogueComponent {
	src := &DialogueComponent{}

	src.View = createTextView("Select a character to start a dialogue.", " NPC Dialogues ", true, Default, Black)
	src.View.SetDynamicColors(true).SetWordWrap(true)

	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(src.View, 0, 1, false)
	src.View.SetText("Hello ! Welcome to Old Assault Party !")

	return src
}

func (c *DialogueComponent) ListenOutputs(app *tview.Application, dialogueChan <-chan string) {
	go func() {
		for msg := range dialogueChan {
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
