package panel

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CommandLineComponent struct {
	Layout  *tview.Flex
	History *tview.TextView
	Input   *tview.InputField
}

func NewCommandLineComponent(app *tview.Application, inputs chan<- string) *CommandLineComponent {
	command_line := &CommandLineComponent{}

	command_line.History = createTextView("", "", true, Default, Black)
	command_line.History.SetBorder(false)

	command_line.Input = tview.NewInputField().
		SetLabel("tap-cli> ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(Black)

	inputRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(command_line.Input, 0, 1, true)

	command_line.Layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(command_line.History, 0, 1, false).
		AddItem(inputRow, 1, 1, true)
	command_line.Layout.SetBorder(true).SetTitle(" CLI ")

	activeColor := tcell.ColorYellow
	inactiveColor := tcell.ColorDimGray
	command_line.Layout.SetBorderColor(inactiveColor)

	command_line.Input.SetFocusFunc(func() {
		command_line.Layout.SetBorderColor(activeColor)
		command_line.Layout.SetTitleColor(activeColor)

		command_line.Input.SetBorderColor(activeColor)
		command_line.Input.SetTitleColor(activeColor)
	})
	command_line.Input.SetBlurFunc(func() {
		command_line.Layout.SetBorderColor(inactiveColor)
		command_line.Layout.SetTitleColor(tcell.ColorWhite)

		command_line.Input.SetBorderColor(inactiveColor)
		command_line.Input.SetTitleColor(tcell.ColorWhite)
	})

	command_line.Input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			text := command_line.Input.GetText()
			if text == "" {
				return
			}

			fmt.Fprint(command_line.History, "tap-cli> "+text+"\n")

			command_line.Input.SetText("")
			inputs <- text
		}
	})

	command_line.Layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(command_line.Input)
			return nil
		}
		return event
	})

	return command_line
}

func (c *CommandLineComponent) ListenOutputs(app *tview.Application, commandLineChan <-chan string) {
	go func() {
		for msg := range commandLineChan {

			app.QueueUpdateDraw(func() {
				fmt.Fprint(c.History, msg+"\n")
			})
		}
	}()
}
