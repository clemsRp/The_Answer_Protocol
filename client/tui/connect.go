package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"io"
	"fmt"
)

func NewConnectComponent(m *MyApp) *tview.InputField {
	connect := createInputField(" Connect ", true, "", Default, Default, Black)

	connect.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			pseudo := connect.GetText()
			if pseudo == "" {
				return
			}

			connect.SetText("")
			inputs <- "CONNECT " + pseudo
		}
	})

	return connect
}

func NewImageComponent(m *MyApp) (*tview.TextView, *tview.TextView) {
    imgView := tview.NewTextView().
        SetDynamicColors(true).
        SetWordWrap(false)
    imgView.SetBorder(true).SetTitle(" [ Front Entrance ] ")

    file, err := os.Open("shopfront.ans")
    if err != nil {
        fmt.Fprintf(imgView, "Error loading image: %v", err)
    } else {
        defer file.Close()
        _, _ = io.Copy(tview.ANSIWriter(imgView), file)
    }

    inputField := tview.NewInputField().SetLabel("Command: ").SetFieldWidth(30)
    inputField.SetBorder(true).SetTitle(" [ Actions ] ")

    logView := tview.NewTextView().SetText("Welcome to Old Assault Party!\nType 'enter' to crash the doors.")
    logView.SetBorder(true).SetTitle(" [ Game Log ] ")

    return imgView, logView
}
