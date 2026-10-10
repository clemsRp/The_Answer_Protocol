package panel

import (
	"bytes"
	"io"
	"strings"
	pr "tap/src/protocol"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	pseudo = ""
)

func SetInputText(input *tview.InputField, text string) {
	_, _, width, _ := input.GetRect()
	nb_spaces := (width - len(text)) / 2
	if nb_spaces < 0 {
		nb_spaces = 0
	}
	spaces := strings.Repeat(" ", nb_spaces)
	input.SetText(spaces + text)
}

func NewConnectComponent(m_pseudo *string, actionsChan chan<- Action) tview.Primitive {
	connect := createInputField(" Connect ", false, "")

	connect.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			r := event.Rune()
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				return event
			}
			return nil
		}
		return event
	})

	connect.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		currentText := connect.GetText()
		cleanText := strings.TrimSpace(currentText)

		nb_spaces := (width - len(cleanText)) / 2
		if nb_spaces < 0 {
			nb_spaces = 0
		}
		spaces := strings.Repeat(" ", nb_spaces)

		if currentText != spaces+cleanText {
			connect.SetText(spaces + cleanText)
		}
		return connect.GetInnerRect()
	})

	connect.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			input_pseudo := strings.TrimSpace(connect.GetText())
			if input_pseudo == "" {
				return
			}

			pseudo = input_pseudo

			SetInputText(connect, "")
			if m_pseudo != nil {
				*m_pseudo = input_pseudo
			}
			actionsChan <- Action{
				Type:    ActionSendServer,
				Payload: pr.CmdConnect + " " + input_pseudo,
			}
		}
	})

	input := createVerticalInputField("\t\t\t\t\t   ENTER PSEUDO", tcell.ColorGreen, connect)
	return input
}

func NewImageComponent(img_data []byte) *tview.TextView {
	imgView := tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(false)

	reader := bytes.NewReader(img_data)
	_, _ = io.Copy(tview.ANSIWriter(imgView), reader)
	return imgView
}
