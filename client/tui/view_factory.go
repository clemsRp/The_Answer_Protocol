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
	inactiveColor := tcell.ColorDimGray
	activeColor := tcell.ColorYellow
	tv.SetBorderColor(inactiveColor)
	tv.SetFocusFunc(func() {
		tv.SetBorderColor(activeColor)
		tv.SetTitleColor(activeColor)
	})
	tv.SetBlurFunc(func() {
		tv.SetBorderColor(inactiveColor)
		tv.SetTitleColor(tcell.ColorWhite)
	})

	return tv
}

func createListView(title string, hasBorder bool, mainColor tcell.Color, shortcutColor tcell.Color, selectedBgColor tcell.Color, backgroundColor tcell.Color, highlightFullLine bool, showSecondaryText bool) *tview.List {
	l := tview.NewList()

	l.SetMainTextColor(mainColor)
	l.SetShortcutColor(shortcutColor)
	l.SetSelectedBackgroundColor(selectedBgColor)
	l.SetBackgroundColor(backgroundColor)
	l.ShowSecondaryText(showSecondaryText)
	l.SetHighlightFullLine(highlightFullLine)
	l.SetBorder(hasBorder)
	if title != "" {
		l.SetTitle(title)
	}
	inactiveColor := tcell.ColorDimGray
	activeColor := tcell.ColorYellow
	l.SetBorderColor(inactiveColor)

	l.SetFocusFunc(func() {
		l.SetBorderColor(activeColor)
		l.SetTitleColor(activeColor)
	})
	l.SetBlurFunc(func() {
		l.SetBorderColor(inactiveColor)
		l.SetTitleColor(tcell.ColorWhite)
	})
	return l
}

func createFormView(title string, hasBorder bool, labelColor tcell.Color, buttonColor tcell.Color, buttonBgColor tcell.Color, backgroundColor tcell.Color) *tview.Form {
	f := tview.NewForm()

	f.SetLabelColor(labelColor)
	f.SetButtonTextColor(buttonColor)
	f.SetButtonBackgroundColor(buttonBgColor)
	f.SetBackgroundColor(backgroundColor)
	f.SetBorder(hasBorder)
	inactiveColor := tcell.ColorDimGray
	activeColor := tcell.ColorYellow
	f.SetBorderColor(inactiveColor)
	if title != "" {
		f.SetTitle(title)
	}

	f.SetFocusFunc(func() {
		f.SetBorderColor(activeColor)
		f.SetTitleColor(activeColor)
	})
	f.SetBlurFunc(func() {
		f.SetBorderColor(inactiveColor)
		f.SetTitleColor(tcell.ColorWhite)
	})

	return f
}

func createInputField(title string, hasBorder bool, label string, labelColor tcell.Color, fieldColor tcell.Color, backgroundColor tcell.Color) *tview.InputField {
	input := tview.NewInputField()

	input.SetLabel(label)
	input.SetLabelColor(labelColor)
	input.SetFieldTextColor(fieldColor)
	input.SetBackgroundColor(backgroundColor)
	input.SetFieldBackgroundColor(backgroundColor)

	input.SetBorder(hasBorder)
	if title != "" {
		input.SetTitle(title)
	}

	return input
}

func createVerticalInputField(labelText string, labelColor tcell.Color, inputField *tview.InputField) tview.Primitive {
	label := tview.NewTextView().
		SetText(labelText).
		SetTextColor(labelColor)

	container := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(label, 1, 1, false).
		AddItem(inputField, 1, 1, true)

	return container
}

func createSelectField(label string, options []string, index int) *tview.DropDown {
	s := tview.NewDropDown().
		SetLabel(label).
		SetOptions(options, nil)

	s.SetCurrentOption(index)

	return s
}
