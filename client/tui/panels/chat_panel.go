package panel

import (
	"fmt"
	"strings"

	pr "tap/protocol"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ChatComponent struct {
	Layout    *tview.Flex
	History   *tview.Pages
	Histories map[string]*tview.TextView
	Scope     *tview.DropDown
	Input     *tview.InputField
}

type SaveChat struct {
	Scope string
	Msg   string
}

var (
	last_chat = SaveChat{Scope: pr.GlobalChat, Msg: ""}
)

func NewChatComponent(app *tview.Application, inputs chan<- string) *ChatComponent {
	chat := &ChatComponent{}

	// Init Scope Histories
	chat.History = tview.NewPages()

	chat.Histories = map[string]*tview.TextView{
		pr.RoomChat:   NewHistoryComponent(pr.RoomChat),
		pr.GroupChat:  NewHistoryComponent(pr.GroupChat),
		pr.GlobalChat: NewHistoryComponent(pr.GlobalChat),
	}

	for scope, history := range chat.Histories {
		focus := false
		if scope == pr.GlobalChat {
			focus = true
		}
		chat.History.AddPage(scope, history, true, focus)
	}

	chat.Scope = createSelectField("Canal: ", []string{pr.GlobalChat, pr.RoomChat, pr.GroupChat}, 0)

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

	chat.Scope.SetSelectedFunc(func(scope string, index int) {
		chat.History.SwitchToPage(scope)
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

			last_chat = SaveChat{Scope: canal, Msg: text}

			chat.Input.SetText("")
			inputs <- fmt.Sprintf("CHAT %s %s", canal, text)
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

func NewHistoryComponent(scope string) *tview.TextView {
	history := createTextView("", "", true)
	history.
		SetBorder(true).
		SetTitle(fmt.Sprintf(" %s ", scope))

	return history
}

func (c *ChatComponent) ListenOutputs(app *tview.Application, chatChan <-chan pr.ServerResponse) {
	go func() {
		for res := range chatChan {
			// OK
			if strings.HasPrefix(res.Msg, "EVT") {
				split_msg := strings.SplitN(res.Msg, " ", 5)

				scope := split_msg[1]
				user := split_msg[3]
				message := split_msg[4]

				lineChat := fmt.Sprintf("[green]%s: [white]%s\n", user, message)

				app.QueueUpdateDraw(func() {
					if historyView, ok := c.Histories[scope]; ok {
						fmt.Fprint(historyView, lineChat)
					}
				})

			} else if res.Msg == "OK" {
				formated_msg := fmt.Sprintf("[green]%s: [white]%s\n", pseudo, last_chat.Msg)

				app.QueueUpdateDraw(func() {
					fmt.Fprint(c.Histories[last_chat.Scope], formated_msg)
				})
			}
		}
	}()
}
