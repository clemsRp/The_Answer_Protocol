package panel

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PopupComponent struct {
	Layout      *tview.Flex
	LayoutTemp  *tview.Flex
	Buttons     *tview.Flex
	MainContent tview.Primitive
	FocusItem   tview.Primitive
	ButtonList  []*tview.Button
}

type InputCapturer interface {
	SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) *tview.Box
	SetMouseCapture(capture func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse)) *tview.Box
}

func NewPopupComponent(app *tview.Application, grid *tview.Grid, mainContent tview.Primitive, contentHeight int, buttons []*tview.Button) *PopupComponent {
	// Color constants
	bgColor := tcell.GetColor("#3a3838")

	popup := PopupComponent{
		LayoutTemp:  tview.NewFlex().SetDirection(tview.FlexRow),
		Layout:      tview.NewFlex().SetDirection(tview.FlexColumn),
		Buttons:     tview.NewFlex().SetDirection(tview.FlexColumn),
		MainContent: mainContent,
		FocusItem:   mainContent,
		ButtonList:  buttons,
	}

	// Apply background colors to containers
	popup.LayoutTemp.SetBackgroundColor(bgColor)
	popup.Buttons.SetBackgroundColor(bgColor)

	// Helper to create background-colored spacers
	createSpacer := func() *tview.Box {
		box := tview.NewBox()
		box.SetBackgroundColor(bgColor)
		return box
	}

	// Assemble buttons dynamically
	for i, btn := range buttons {
		popup.Buttons.AddItem(btn, 0, 1, i == 0)
		if i < len(buttons)-1 {
			popup.Buttons.AddItem(createSpacer(), 1, 1, false)
		}
	}

	// Assemble vertical inner layout
	popup.LayoutTemp.AddItem(createSpacer(), 5, 0, false)
	popup.LayoutTemp.AddItem(mainContent, 0, 1, true)
	popup.LayoutTemp.AddItem(createSpacer(), 3, 0, false)

	// Add buttons row only if buttons are provided
	if len(buttons) > 0 {
		popup.LayoutTemp.AddItem(popup.Buttons, 1, 0, false)
		popup.LayoutTemp.AddItem(createSpacer(), 5, 0, false)
	}

	// Build focusable primitives slice for Tab key navigation
	focusableItems := []tview.Primitive{mainContent}
	for _, btn := range buttons {
		focusableItems = append(focusableItems, btn)
	}

	// Handle focus navigation across main content and buttons
	popup.LayoutTemp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			currentFocus := app.GetFocus()
			for i, item := range focusableItems {
				if currentFocus == item {
					nextIndex := (i + 1) % len(focusableItems)
					app.SetFocus(focusableItems[nextIndex])
					return nil
				}
			}
			if len(focusableItems) > 0 {
				app.SetFocus(focusableItems[0])
			}
			return nil
		}
		return event
	})

	popup.Layout.AddItem(createSpacer(), 15, 0, false)
	popup.Layout.AddItem(popup.LayoutTemp, 0, 1, true)
	popup.Layout.AddItem(createSpacer(), 15, 0, false)

	// Total vertical height calculation for container grid
	totalHeight := contentHeight + 3
	if len(buttons) > 0 {
		totalHeight += 2
	}
	grid.SetRows(0, totalHeight-1, 0)

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
