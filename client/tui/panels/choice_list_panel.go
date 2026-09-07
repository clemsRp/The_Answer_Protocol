package panel

import (
	"context"
	"fmt"
	"strings"
	"sync"
	pr "tap/protocol"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type OptionsMap map[string]map[string]func()
type AllowedOptions interface {
	OptionsMap | map[string]OptionsMap
}

type ChoiceListComponent struct {
	Layout *tview.Flex
	List   *tview.List
}

var (
	popupBgColor = AppTheme.PopupBackground
	btnRestBg    = tcell.GetColor("#474646")
	btnActiveBg  = tcell.GetColor("#7e7979")
)

type entryItem struct {
	isHeader bool
	locName  string
	areBtns  bool
	action   func()
}

func NewChoiceListComponent[T AllowedOptions](
	app *tview.Application,
	popupGrid *tview.Grid,
	title string,
	options T,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
	are_btns bool,
) *ChoiceListComponent {

	src := &ChoiceListComponent{}
	src.List = createListView(" "+title+" ", true, false, true)
	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(src.List, 0, 1, false)

	var entries []entryItem

	switch opts := any(options).(type) {
	case OptionsMap:
		entries = append(entries, collectOptions(app, src, popupGrid, opts, onOpenPopup, onClosePopup, are_btns)...)

	case map[string]OptionsMap:
		for subTitle, subOptions := range opts {
			entries = append(entries, entryItem{
				isHeader: true,
				locName:  subTitle,
			})
			entries = append(entries, collectOptions(app, src, popupGrid, subOptions, onOpenPopup, onClosePopup, are_btns)...)
		}
	}

	formatItem := func(locName string, areBtns bool, isSelected bool) (string, string) {
		if areBtns {
			if isSelected {
				return "[yellow][ " + locName + " ][-]", ""
			}
			return fmt.Sprintf("[yellow:%s][ %s ][-:-]", AppTheme.PopupBackgroundHexa, locName), ""
		} else {
			parts := strings.SplitN(locName, "\n", 2)
			mainContent := parts[0]
			secondaryContent := ""
			if len(parts) > 1 {
				secondaryContent = parts[1]
			}

			if isSelected {
				mText := " ●  " + mainContent
				sText := ""
				if secondaryContent != "" {
					sText = secondaryContent
				}
				return mText, sText
			} else {
				mText := "○  " + mainContent
				sText := ""
				if secondaryContent != "" {
					sText = secondaryContent
				}
				return mText, sText
			}
		}
	}

	for idx, entry := range entries {
		if entry.isHeader {
			src.List.AddItem("[yellow:#000000]- "+entry.locName+":", "", 0, nil)
		} else {
			isSelected := (idx == 0)
			mText, sText := formatItem(entry.locName, entry.areBtns, isSelected)
			src.List.AddItem(mText, sText, 0, entry.action)
		}
	}

	src.List.SetChangedFunc(func(i int, mainText, secondaryText string, shortcut rune) {
		for idx, entry := range entries {
			if entry.isHeader {
				continue
			}
			isSelected := (idx == i)
			mText, sText := formatItem(entry.locName, entry.areBtns, isSelected)
			src.List.SetItemText(idx, mText, sText)
		}
	})

	return src
}

func collectOptions(
	app *tview.Application,
	src *ChoiceListComponent,
	popupGrid *tview.Grid,
	options OptionsMap,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
	are_btns bool,
) []entryItem {
	var entries []entryItem

	for location, actions := range options {
		locActions := actions
		locName := location

		itemAction := func() {
			var selectedFunc func()

			optionsFlex := tview.NewFlex().SetDirection(tview.FlexRow)
			optionsFlex.SetBackgroundColor(popupBgColor)

			makeSpacer := func() *tview.Box {
				spacer := tview.NewBox()
				spacer.SetBackgroundColor(popupBgColor)
				return spacer
			}

			optionsFlex.AddItem(makeSpacer(), 1, 0, false)

			actionList := tview.NewList().
				SetMainTextColor(tcell.ColorWhite).
				SetSelectedBackgroundColor(tcell.ColorRed).
				SetSelectedTextColor(tcell.ColorWhite)

			actionList.SetBackgroundColor(popupBgColor)

			funcsIndices := make([]func(), 0, len(locActions))
			cmdNames := make([]string, 0, len(locActions))

			formatPopupItem := func(name string, isSelected bool) string {
				if isSelected {
					return "[white:#7e7979]" + transform_name(name, 30)
				}
				return "[white:#474646]" + transform_name(name, 30)
			}

			for cmdName, cmdFunc := range locActions {
				cFunc := cmdFunc
				funcsIndices = append(funcsIndices, cFunc)
				cmdNames = append(cmdNames, cmdName)

				isFirst := len(cmdNames) == 1
				actionList.AddItem(formatPopupItem(cmdName, isFirst), "", 0, nil)
			}

			if len(funcsIndices) > 0 {
				selectedFunc = funcsIndices[0]
			}

			actionList.SetChangedFunc(func(i int, mainText, secondaryText string, shortcut rune) {
				if i >= 0 && i < len(funcsIndices) {
					selectedFunc = funcsIndices[i]

					itemCount := actionList.GetItemCount()
					for idx, name := range cmdNames {
						if idx < itemCount {
							isSelected := (idx == i)
							actionList.SetItemText(idx, formatPopupItem(name, isSelected), "")
						}
					}
				}
			})

			actionList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyEnter {
					return nil
				}
				return event
			})

			optionsFlex.AddItem(actionList, len(locActions)*2, 1, true)

			// Bouton Cancel
			cancelBtn := tview.NewButton("Cancel").
				SetLabelColor(tcell.ColorWhite).
				SetBackgroundColorActivated(btnActiveBg).
				SetLabelColorActivated(tcell.ColorWhite).
				SetSelectedFunc(func() {
					if onClosePopup != nil {
						onClosePopup()
					}
				})
			cancelBtn.SetBackgroundColor(btnRestBg)

			validateBtn := tview.NewButton("Validate").
				SetLabelColor(tcell.ColorWhite).
				SetBackgroundColorActivated(btnActiveBg).
				SetLabelColorActivated(tcell.ColorWhite).
				SetSelectedFunc(func() {
					if selectedFunc == nil {
						return
					}

					fnToExecute := selectedFunc
					selectedFunc = nil

					if onClosePopup != nil {
						onClosePopup()
					}

					fnToExecute()
				})
			validateBtn.SetBackgroundColor(btnRestBg)

			buttons := []*tview.Button{cancelBtn, validateBtn}
			totalHeight := (len(locActions) * 2) + 11

			newPopup := NewPopupComponent(app, popupGrid, optionsFlex, totalHeight, buttons)

			popupGrid.Clear()
			popupGrid.AddItem(newPopup.Layout, 1, 1, 1, 1, 0, 0, true)

			if onOpenPopup != nil {
				onOpenPopup(newPopup)
			}
		}

		entries = append(entries, entryItem{
			isHeader: false,
			locName:  locName,
			areBtns:  are_btns,
			action:   itemAction,
		})
	}

	return entries
}

func transform_name(option string, width int) string {
	if len(option) >= width {
		return option
	}
	first_len := (width - len(option)) / 2
	if first_len < 0 {
		first_len = 0
	}
	remaining := width - first_len - len(option)
	if remaining < 0 {
		remaining = 0
	}
	return strings.Repeat(" ", first_len) + option + strings.Repeat(" ", remaining)
}

func (c *ChoiceListComponent) ListenOutputs(ctx context.Context, wg *sync.WaitGroup, app *tview.Application, Chan <-chan pr.ServerResponse, function func(pr.ServerResponse)) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case res, ok := <-Chan:
				if !ok {
					return
				}
				response := res
				app.QueueUpdateDraw(func() {
					function(response)
				})
			}
		}
	}()
}
