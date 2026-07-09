package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io"
	"os"
	"strings"
	"unicode"
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

func NewConnectComponent(m *MyApp) tview.Primitive {
	connect := createInputField(" Connect ", false, "", tcell.ColorGreen, tcell.ColorGreen, tcell.NewRGBColor(0, 0, 0))

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
			texteSaisi := strings.TrimSpace(connect.GetText())
			if texteSaisi == "" {
				return
			}

			pseudo = texteSaisi
			SetInputText(connect, "")
			inputs <- "CONNECT " + pseudo
		}
	})

	input := createVerticalInputField("\t\t\t\t\t   ENTER PSEUDO", tcell.ColorGreen, connect)
	return input
}

func NewImageComponent(m *MyApp, img_path string) *tview.TextView {
	imgView := tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(false)

	file, err := os.Open(img_path)
	if err != nil {
		fmt.Fprintf(imgView, "Error loading image: %v", err)
		return imgView
	}
	defer file.Close()
	
	_, _ = io.Copy(tview.ANSIWriter(imgView), file)
	return imgView
}