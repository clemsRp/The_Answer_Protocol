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

func NewChatComponent(app *tview.Application, pseudo string, inputs chan<- string) *ChatComponent {
	chat := &ChatComponent{}

	chat.History = createTextView("", "", false, Default, Black)

	chat.Scope = createSelectField("Canal: ", []string{"GLOBAL", "ROOM", "GROUP"}, 1)

	chat.Input = tview.NewInputField().
		SetLabel(" Message: ").
		SetFieldWidth(0)
	inputRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(chat.Scope, 17, 1, false).
		AddItem(chat.Input, 0, 1, true)

	chat.Layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(chat.History, 0, 1, false).
		AddItem(inputRow, 1, 1, true)
	chat.Layout.SetBorder(true).SetTitle(" Chat")

	activeColor := tcell.ColorYellow
	inactiveColor := tcell.ColorDimGray
	chat.Layout.SetBorderColor(inactiveColor)

	chat.Input.SetFocusFunc(func() {
		chat.Layout.SetBorderColor(activeColor)
		chat.Layout.SetTitleColor(activeColor)

		chat.Input.SetBorderColor(activeColor)
		chat.Input.SetTitleColor(activeColor)
	})
	chat.Input.SetBlurFunc(func() {
		chat.Layout.SetBorderColor(inactiveColor)
		chat.Layout.SetTitleColor(tcell.ColorWhite)

		chat.Input.SetBorderColor(inactiveColor)
		chat.Input.SetTitleColor(tcell.ColorWhite)
	})
	chat.Scope.SetFocusFunc(func() {
		chat.Layout.SetBorderColor(activeColor)
		chat.Layout.SetTitleColor(activeColor)
	})
	chat.Scope.SetBlurFunc(func() {
		chat.Layout.SetBorderColor(inactiveColor)
		chat.Layout.SetTitleColor(tcell.ColorWhite)
	})

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
			inputs <- fmt.Sprintf("%s %s", canal, text)
		}
	})

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

func (c *ChatComponent) ListenOutputs(app *tview.Application, chatChan <-chan string) {
	// TO CHANGE

	go func() {
		for msg := range chatChan {

			lineChat := fmt.Sprintf("[gray][%s] [yellow][%s] [green]%s: [white]%s\n",
				msg, msg, msg, msg)

			app.QueueUpdateDraw(func() {
				fmt.Fprint(c.History, lineChat)
			})
		}
	}()
}
