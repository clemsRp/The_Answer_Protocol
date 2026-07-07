package main

import (
	"fmt"
	"time"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ChatComponent struct {
	Layout  *tview.Flex
	History *tview.TextView
	Scope   *tview.DropDown
	Input   *tview.InputField
}

func NewChatComponent(app *tview.Application, pseudo string) *ChatComponent {
	chat := &ChatComponent{}

	chat.History = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true)
	chat.History.
		SetBorder(true).
		SetTitle(" Chat History ")

	chat.Scope = createSelectField("Canal: ", []string{"GLOBAL", "ROOM", "GROUP"}, 0)

	chat.Input = tview.NewInputField().
		SetLabel(" Message: ").
		SetFieldWidth(0)

	chat.Input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			text := chat.Input.GetText()
			if text == "" {
				return
			}

			_, canal := chat.Scope.GetCurrentOption()
			hour := time.Now().Format("15:04")

			formated_msg := fmt.Sprintf("[gray][%s] [yellow][%s] [green]%s: [white]%s\n", hour, canal, pseudo, text)

			fmt.Fprint(chat.History, formated_msg)

			chat.Input.SetText("")
			inputs <- fmt.Sprintf("CHAT %s %s", canal, text)
		}
	})

	inputRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(chat.Scope, 17, 1, false).
		AddItem(chat.Input, 0, 1, true)

	chat.Layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(chat.History, 0, 1, false).
		AddItem(inputRow, 1, 1, true)


	chat.Layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if app.GetFocus() == chat.Input {
				app.SetFocus(chat.Scope)
			} else {
				app.SetFocus(chat.Input)
			}
			return nil
		}
		return event
	})

	return chat
}

func (c *ChatComponent) ListenOutputs(app *tview.Application, outputs <-chan string) {
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
