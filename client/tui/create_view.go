package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createTextView(text string, title string, hasBorder bool, textColor tcell.Color, backgroundColor tcell.Color) *tview.TextView {
	tv := tview.NewTextView()

	tv.SetText(text)
	tv.SetTextColor(textColor)
	tv.SetBackgroundColor(backgroundColor)

	tv.SetBorder(hasBorder)
	if title != "" {
		tv.SetTitle(title)
	}

	return tv
}

func createListView(title string, hasBorder bool, mainColor tcell.Color, shortcutColor tcell.Color, selectedBgColor tcell.Color, backgroundColor tcell.Color) *tview.List {
	l := tview.NewList()

	l.SetMainTextColor(mainColor)
	l.SetShortcutColor(shortcutColor)
	l.SetSelectedBackgroundColor(selectedBgColor)
	l.SetBackgroundColor(backgroundColor)

	l.SetBorder(hasBorder)
	if title != "" {
		l.SetTitle(title)
	}

	return l
}

func createFormView(title string, hasBorder bool, labelColor tcell.Color, buttonColor tcell.Color, buttonBgColor tcell.Color, backgroundColor tcell.Color) *tview.Form {
	f := tview.NewForm()

	f.SetLabelColor(labelColor)
	f.SetButtonTextColor(buttonColor)
	f.SetButtonBackgroundColor(buttonBgColor)
	f.SetBackgroundColor(backgroundColor)

	f.SetBorder(hasBorder)
	if title != "" {
		f.SetTitle(title)
	}

	return f
}

func createInputField(title string, hasBorder bool, label string, labelColor tcell.Color, fieldColor tcell.Color, backgroundColor tcell.Color) *tview.InputField {
	input := tview.NewInputField()

	input.SetLabel(label)
	input.SetLabelColor(labelColor)
	input.SetFieldTextColor(fieldColor)
	input.SetBackgroundColor(backgroundColor)

	input.SetBorder(hasBorder)
	if title != "" {
		input.SetTitle(title)
	}

	return input
}
