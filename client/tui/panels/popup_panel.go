package panel

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PopupComponent struct {
	Layout      *tview.Flex
	LayoutTemp  *tview.Flex
	ValidateBtn *tview.Button
	CancelBtn   *tview.Button
	Buttons     *tview.Flex
	OptionsGrid *tview.Grid
	OptionsList *tview.List
}

type InputCapturer interface {
	SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) *tview.Box
	SetMouseCapture(capture func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse)) *tview.Box
}

func NewPopupComponent(app *tview.Application, grid *tview.Grid, options map[string]func(), quitFunc func()) *PopupComponent {
	// Color constants
	bgColor := tcell.GetColor("#3a3838")
	selectedBgColor := tcell.GetColor("#6e6c6c")
	unselectedBgHex := "#2f2e2e"

	// Layout dimensions
	optionWidth := 15
	totalRows := 2*len(options) - 1

	popup := PopupComponent{
		LayoutTemp:  tview.NewFlex().SetDirection(tview.FlexRow),
		Layout:      tview.NewFlex().SetDirection(tview.FlexColumn),
		ValidateBtn: tview.NewButton("Validate"),
		CancelBtn:   tview.NewButton("Cancel"),
		Buttons:     tview.NewFlex().SetDirection(tview.FlexColumn),
		OptionsList: tview.NewList(),
		OptionsGrid: tview.NewGrid().SetRows(0, totalRows, 0).SetColumns(0, optionWidth, 0),
	}

	// Apply background colors to containers
	popup.LayoutTemp.SetBackgroundColor(bgColor)
	popup.Buttons.SetBackgroundColor(bgColor)
	popup.OptionsList.SetBackgroundColor(bgColor)
	popup.OptionsGrid.SetBackgroundColor(bgColor)

	// Configure list styling
	popup.OptionsList.ShowSecondaryText(false)
	popup.OptionsList.SetWrapAround(false)
	popup.OptionsList.SetMainTextColor(tcell.ColorWhite)
	popup.OptionsList.SetSecondaryTextColor(tcell.ColorGray)
	popup.OptionsList.SetSelectedBackgroundColor(selectedBgColor)
	popup.OptionsList.SetSelectedTextColor(tcell.ColorWhite)

	// Helpers for option formatting and cleaning
	tagRegex := regexp.MustCompile(`\[[^\]]*\]`)
	stripTags := func(text string) string {
		return tagRegex.ReplaceAllString(text, "")
	}

	formatOption := func(text string, isSelected bool) string {
		if isSelected {
			return text
		}
		return fmt.Sprintf("[:%s]%s[:-]", unselectedBgHex, text)
	}

	// Configure button actions
	popup.CancelBtn.SetSelectedFunc(func() {
		go quitFunc()
	})

	popup.ValidateBtn.SetSelectedFunc(func() {
		index := popup.OptionsList.GetCurrentItem()
		text, _ := popup.OptionsList.GetItemText(index)
		rawText := stripTags(text)

		function := options[strings.TrimSpace(rawText)]

		if function != nil {
			go function()
		}

		go quitFunc()
	})

	// Populate options list with separators
	index := 0
	for option := range options {
		paddingLeft := (optionWidth - len(option)) / 2
		paddingRight := optionWidth - len(option) - paddingLeft
		fullText := strings.Repeat(" ", paddingLeft) + option + strings.Repeat(" ", paddingRight)

		itemText := formatOption(fullText, index == 0)
		popup.OptionsList.AddItem(itemText, "", 0, nil)

		if index < len(options)-1 {
			popup.OptionsList.AddItem("", "", 0, nil)
		}

		index++
	}

	// Update items dynamically on selection change
	popup.OptionsList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		for i := 0; i < popup.OptionsList.GetItemCount(); i++ {
			if i%2 != 0 {
				continue
			}

			text, _ := popup.OptionsList.GetItemText(i)
			rawText := stripTags(text)
			isSelected := (i == index)

			popup.OptionsList.SetItemText(i, formatOption(rawText, isSelected), "")
		}
	})

	// Handle list navigation skipping empty separators
	popup.OptionsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		current := popup.OptionsList.GetCurrentItem()
		if event.Key() == tcell.KeyDown {
			if current+2 < popup.OptionsList.GetItemCount() {
				popup.OptionsList.SetCurrentItem(current + 2)
			}
			return nil
		} else if event.Key() == tcell.KeyUp {
			if current-2 >= 0 {
				popup.OptionsList.SetCurrentItem(current - 2)
			}
			return nil
		}
		return event
	})

	// Create reusable spacer box
	createSpacer := func() *tview.Box {
		box := tview.NewBox()
		box.SetBackgroundColor(bgColor)
		return box
	}

	// Assemble options grid
	popup.OptionsGrid.AddItem(createSpacer(), 1, 0, 1, 1, 0, 0, false)
	popup.OptionsGrid.AddItem(popup.OptionsList, 1, 1, 1, 1, 0, 0, true)
	popup.OptionsGrid.AddItem(createSpacer(), 1, 2, 1, 1, 0, 0, false)

	// Assemble buttons container
	popup.Buttons.
		AddItem(popup.CancelBtn, 0, 1, true).
		AddItem(createSpacer(), 1, 1, false).
		AddItem(popup.ValidateBtn, 0, 1, true)

	// Assemble vertical inner layout
	popup.LayoutTemp.AddItem(createSpacer(), 1, 0, false)
	popup.LayoutTemp.AddItem(popup.OptionsGrid, totalRows, 1, true)
	popup.LayoutTemp.AddItem(createSpacer(), 1, 0, false)
	popup.LayoutTemp.AddItem(popup.Buttons, 1, 0, true)
	popup.LayoutTemp.AddItem(createSpacer(), 1, 0, false)

	// Handle focus tab key navigation
	popup.LayoutTemp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if popup.OptionsList.HasFocus() {
				app.SetFocus(popup.CancelBtn)
			} else if popup.CancelBtn.HasFocus() {
				app.SetFocus(popup.ValidateBtn)
			} else {
				app.SetFocus(popup.OptionsList)
			}
			return nil
		}
		return event
	})

	// Assemble final outer layout
	popup.Layout.AddItem(createSpacer(), 0, 1, false)
	popup.Layout.AddItem(popup.LayoutTemp, 30, 1, true)
	popup.Layout.AddItem(createSpacer(), 0, 1, false)

	grid.SetRows(0, totalRows+4, 0)

	return &popup
}

func SetBlockedInputs(component InputCapturer, is_blocked bool) {
	// Open all inputs
	if is_blocked {
		component.SetInputCapture(nil)
		component.SetMouseCapture(nil)
		return
	}

	// Blocked all inputs
	component.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return nil
	})

	component.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		return action, nil
	})
}
