package panel

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Focusable interface {
	SetFocusFunc(callback func()) *tview.Box
	SetBlurFunc(callback func()) *tview.Box
}

func BindContainerFocus(container *tview.Flex, children ...Focusable) {
	for _, child := range children {
		child.SetFocusFunc(func() {
			container.SetBorderColor(AppTheme.BorderActive)
			container.SetTitleColor(AppTheme.TitleActive)
		})
		child.SetBlurFunc(func() {
			container.SetBorderColor(AppTheme.BorderInactive)
			container.SetTitleColor(AppTheme.TitleInactive)
		})
	}
}

func createTextView(text string, title string, hasBorder bool) *tview.TextView {
	tv := tview.NewTextView()
	tv.SetText(text)
	tv.SetTextColor(AppTheme.TextPrimary)
	tv.SetBackgroundColor(AppTheme.Background)
	tv.SetDynamicColors(true)
	tv.SetBorder(hasBorder)
	tv.SetBorderColor(AppTheme.BorderInactive)

	if title != "" {
		tv.SetTitle(title)
		tv.SetTitleColor(AppTheme.TitleInactive)
	}

	tv.SetFocusFunc(func() {
		tv.SetBorderColor(AppTheme.BorderActive)
		tv.SetTitleColor(AppTheme.TitleActive)
	})
	tv.SetBlurFunc(func() {
		tv.SetBorderColor(AppTheme.BorderInactive)
		tv.SetTitleColor(AppTheme.TitleInactive)
	})
	return tv
}

func createListView(title string, hasBorder bool, highlightFullLine bool, showSecondaryText bool) *tview.List {
	l := tview.NewList()
	l.SetMainTextColor(AppTheme.TextPrimary)
	l.SetShortcutColor(AppTheme.TextHighlight)
	l.SetSelectedBackgroundColor(AppTheme.ListSelectedBg)
	l.SetSelectedTextColor(AppTheme.ListSelectedTxt)
	l.SetBackgroundColor(AppTheme.Background)
	l.ShowSecondaryText(showSecondaryText)
	l.SetHighlightFullLine(highlightFullLine)
	l.SetBorder(hasBorder)
	l.SetBorderColor(AppTheme.BorderInactive)

	if title != "" {
		l.SetTitle(title)
		l.SetTitleColor(AppTheme.TitleInactive)
	}

	l.SetFocusFunc(func() {
		l.SetBorderColor(AppTheme.BorderActive)
		l.SetTitleColor(AppTheme.TitleActive)
	})
	l.SetBlurFunc(func() {
		l.SetBorderColor(AppTheme.BorderInactive)
		l.SetTitleColor(AppTheme.TitleInactive)
	})
	return l
}

func createInputField(title string, hasBorder bool, label string) *tview.InputField {
	input := tview.NewInputField()
	input.SetLabel(label)
	input.SetLabelColor(AppTheme.TextSecondary)
	input.SetFieldTextColor(AppTheme.TextPrimary)
	input.SetBackgroundColor(AppTheme.Background)
	input.SetFieldBackgroundColor(AppTheme.Background)
	input.SetBorder(hasBorder)
	input.SetBorderColor(AppTheme.BorderInactive)

	if title != "" {
		input.SetTitle(title)
		input.SetTitleColor(AppTheme.TitleInactive)
	}

	input.SetFocusFunc(func() {
		input.SetBorderColor(AppTheme.BorderActive)
		input.SetTitleColor(AppTheme.TitleActive)
	})
	input.SetBlurFunc(func() {
		input.SetBorderColor(AppTheme.BorderInactive)
		input.SetTitleColor(AppTheme.TitleInactive)
	})
	return input
}

func createSelectField(label string, options []string, index int) *tview.DropDown {
	s := tview.NewDropDown().
		SetLabel(label).
		SetOptions(options, nil)

	s.SetLabelColor(AppTheme.TextSecondary)
	s.SetFieldTextColor(AppTheme.TextPrimary)
	s.SetFieldBackgroundColor(AppTheme.Background)
	s.SetBackgroundColor(AppTheme.Background)
	s.SetCurrentOption(index)
	return s
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

func SeparateElements(elements []tview.Primitive, sizes []int, is_row bool) *tview.Flex {
	flex := tview.NewFlex()
	var dir int
	if is_row {
		dir = tview.FlexColumn
	} else {
		dir = tview.FlexRow
	}
	flex.SetDirection(dir)

	makeSpacer := func(backgroundColor tcell.Color) *tview.Box {
		spacer := tview.NewBox()
		spacer.SetBackgroundColor(backgroundColor)
		return spacer
	}

	makeLineSpacer := func() *tview.Box {
		box := tview.NewBox()
		box.SetBackgroundColor(AppTheme.PopupBackground)

		box.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
			style := tcell.StyleDefault.
				Background(AppTheme.PopupBackground).
				Foreground(AppTheme.Background)

			if is_row {
				for row := y; row < y+height; row++ {
					screen.SetContent(x, row, '│', nil, style)
				}
			} else {
				for col := x; col < x+width; col++ {
					screen.SetContent(col, y, '─', nil, style)
				}
			}

			return x, y, width, height
		})

		return box
	}

	for index, ele := range elements {
		size := 0
		if index < len(sizes) {
			size = sizes[index]
		}
		flex.AddItem(ele, size, 1, true)

		if index != len(elements)-1 {
			flex.AddItem(makeSpacer(AppTheme.PopupBackground), 1, 0, false)
			flex.AddItem(makeLineSpacer(), 1, 0, false)
			flex.AddItem(makeSpacer(AppTheme.PopupBackground), 1, 0, false)
		}
	}

	return flex
}
