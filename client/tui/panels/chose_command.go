package panel

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewCommandComponent(app *tview.Application, grid *tview.Grid, options map[string]func(), quitFunc func()) *PopupComponent {
	// Color constants
	bgColor := tcell.GetColor("#3a3838")
	selectedBgColor := tcell.GetColor("#6e6c6c")
	unselectedBgHex := "#2f2e2e"

	// Layout dimensions
	optionWidth := 15
	totalRows := 2*len(options) - 1

	// Options list setup
	optionsList := tview.NewList()
	optionsList.SetBackgroundColor(bgColor)
	optionsList.ShowSecondaryText(false)
	optionsList.SetWrapAround(false)
	optionsList.SetMainTextColor(tcell.ColorWhite)
	optionsList.SetSecondaryTextColor(tcell.ColorGray)
	optionsList.SetSelectedBackgroundColor(selectedBgColor)
	optionsList.SetSelectedTextColor(tcell.ColorWhite)

	// Format helpers
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

	// Populate options list with separators
	index := 0
	for option := range options {
		paddingLeft := (optionWidth - len(option)) / 2
		paddingRight := optionWidth - len(option) - paddingLeft
		fullText := strings.Repeat(" ", paddingLeft) + option + strings.Repeat(" ", paddingRight)

		itemText := formatOption(fullText, index == 0)
		optionsList.AddItem(itemText, "", 0, nil)

		if index < len(options)-1 {
			optionsList.AddItem("", "", 0, nil)
		}

		index++
	}

	// Update items dynamically on selection change
	optionsList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		for i := 0; i < optionsList.GetItemCount(); i++ {
			if i%2 != 0 {
				continue
			}

			text, _ := optionsList.GetItemText(i)
			rawText := stripTags(text)
			isSelected := (i == index)

			optionsList.SetItemText(i, formatOption(rawText, isSelected), "")
		}
	})

	// Handle list navigation skipping empty separators
	optionsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		current := optionsList.GetCurrentItem()
		if event.Key() == tcell.KeyDown {
			if current+2 < optionsList.GetItemCount() {
				optionsList.SetCurrentItem(current + 2)
			}
			return nil
		} else if event.Key() == tcell.KeyUp {
			if current-2 >= 0 {
				optionsList.SetCurrentItem(current - 2)
			}
			return nil
		}
		return event
	})

	// Wrap list inside horizontal grid spacers
	createSpacer := func() *tview.Box {
		box := tview.NewBox()
		box.SetBackgroundColor(bgColor)
		return box
	}

	optionsGrid := tview.NewGrid().
		SetRows(0, totalRows, 0).
		SetColumns(0, optionWidth, 0)

	optionsGrid.SetBackgroundColor(bgColor)
	optionsGrid.AddItem(createSpacer(), 1, 0, 1, 1, 0, 0, false)
	optionsGrid.AddItem(optionsList, 1, 1, 1, 1, 0, 0, true)
	optionsGrid.AddItem(createSpacer(), 1, 2, 1, 1, 0, 0, false)

	// Buttons setup
	cancelBtn := tview.NewButton("Cancel").SetSelectedFunc(func() {
		if quitFunc != nil {
			go quitFunc()
		}
	})

	validateBtn := tview.NewButton("Validate").SetSelectedFunc(func() {
		idx := optionsList.GetCurrentItem()
		text, _ := optionsList.GetItemText(idx)
		rawText := stripTags(text)

		if fn, exists := options[strings.TrimSpace(rawText)]; exists && fn != nil {
			go fn()
		}

		go quitFunc()
	})

	buttons := []*tview.Button{cancelBtn, validateBtn}

	// Create generic popup instance
	popup := NewPopupComponent(app, grid, optionsGrid, totalRows, buttons)
	popup.FocusItem = optionsList

	return popup
}
