package panel

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PopupComponent struct {
	Layout      *tview.Flex
	ValidateBtn *tview.Button
	CancelBtn   *tview.Button
	Buttons     *tview.Flex
	Options     *tview.List
}

func NewPopupComponent(app *tview.Application, options []string, cancelFunc func(), validateFunc func()) *PopupComponent {
	popup := PopupComponent{
		Layout:      tview.NewFlex().SetDirection(tview.FlexRow),
		ValidateBtn: tview.NewButton("Validate"),
		CancelBtn:   tview.NewButton("Cancel"),
		Buttons:     tview.NewFlex().SetDirection(tview.FlexColumn),
		Options:     tview.NewList(),
	}

	popup.CancelBtn.
		SetSelectedFunc(func() {
			go cancelFunc()
		})

	popup.ValidateBtn.
		SetSelectedFunc(func() {
			go validateFunc()
		})

	popup.Layout.SetBorderPadding(2, 2, 5, 5)

	// Set bg colors
	bgColor := tcell.GetColor("#2c2b2b")

	popup.Layout.SetBackgroundColor(bgColor)
	popup.Options.SetBackgroundColor(bgColor)
	popup.Buttons.SetBackgroundColor(bgColor)
	popup.ValidateBtn.SetBackgroundColor(bgColor)
	popup.CancelBtn.SetBackgroundColor(bgColor)

	option_width := 15

	// Set options
	for index, option := range options {
		option_str := strings.Repeat(" ", (option_width-len(option))/2) + option
		popup.Options.
			AddItem(option_str+strings.Repeat(" ", option_width-len(option_str)), "", 0, nil)
		index++
	}

	// Set Buttons
	popup.Buttons.
		AddItem(popup.CancelBtn, 0, 1, true).
		AddItem(tview.NewBox(), 1, 1, false).
		AddItem(popup.ValidateBtn, 0, 1, true)

	// Set Layout
	popup.Layout.AddItem(popup.Options, 0, 1, true)
	popup.Layout.AddItem(popup.Buttons, 1, 1, true)

	popup.Layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if popup.Options.HasFocus() {
				app.SetFocus(popup.CancelBtn)
			} else if popup.CancelBtn.HasFocus() {
				app.SetFocus(popup.ValidateBtn)
			} else {
				app.SetFocus(popup.Options)
			}
			return nil
		}
		return event
	})

	return &popup
}
