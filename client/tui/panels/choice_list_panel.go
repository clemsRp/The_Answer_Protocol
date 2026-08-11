package panel

import (
	"strings"
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
	popupBgColor = tcell.GetColor("#3a3838")
	btnRestBg    = tcell.GetColor("#474646")
	btnActiveBg  = tcell.GetColor("#7e7979")
)

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

	switch opts := any(options).(type) {
	case OptionsMap:
		createList(app, src, popupGrid, opts, onOpenPopup, onClosePopup, are_btns)

	case map[string]OptionsMap:
		for subTitle, subOptions := range opts {
			src.List.AddItem("[yellow:#000000]- "+subTitle+":", "", 0, nil)

			createList(app, src, popupGrid, subOptions, onOpenPopup, onClosePopup, are_btns)
		}
	}

	return src
}

func createList(
	app *tview.Application,
	src *ChoiceListComponent,
	popupGrid *tview.Grid,
	options OptionsMap,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
	are_btns bool,
) {
	index := 1
	for location, actions := range options {
		locActions := actions
		locName := location

		itemAction := func() {
			var selectedFunc func()

			optionsFlex := tview.NewFlex().SetDirection(tview.FlexRow)
			optionsFlex.SetBackgroundColor(popupBgColor)

			makeSpacer := func(height int) *tview.Box {
				spacer := tview.NewBox()
				spacer.SetBackgroundColor(popupBgColor)
				return spacer
			}

			optionsFlex.AddItem(makeSpacer(1), 1, 0, false)

			actionList := tview.NewList().
				SetMainTextColor(tcell.ColorWhite).
				SetSelectedBackgroundColor(tcell.ColorRed).
				SetSelectedTextColor(tcell.ColorWhite)

			actionList.SetBackgroundColor(popupBgColor)

			funcsIndices := make([]func(), 0, len(locActions))
			cmdNames := make([]string, 0, len(locActions))

			formatItem := func(name string, isSelected bool) string {
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
				actionList.AddItem(formatItem(cmdName, isFirst), "", 0, nil)
			}

			if len(funcsIndices) > 0 {
				selectedFunc = funcsIndices[0]
			}

			actionList.SetChangedFunc(func(i int, mainText, secondaryText string, shortcut rune) {
				if i >= 0 && i < len(funcsIndices) {
					selectedFunc = funcsIndices[i]

					for idx, name := range cmdNames {
						isSelected := (idx == i)
						actionList.SetItemText(idx, formatItem(name, isSelected), "")
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

		if are_btns {
			label := "[yellow][ " + locName + " ][-]"
			src.List.AddItem(label, "", 0, itemAction)

		} else {
			src.List.AddItem("  "+locName+"  ", "", 0, itemAction)
		}
		index++
	}
}

func transform_name(option string, width int) string {
	first_len := (width - len(option)) / 2
	return strings.Repeat(" ", first_len) + option + strings.Repeat(" ", width-first_len-len(option))
}

func (c *ChoiceListComponent) ListenOutputs(app *tview.Application, Chan <-chan pr.ServerResponse, function func(pr.ServerResponse)) {
	go func() {
		for res := range Chan {
			response := res
			app.QueueUpdateDraw(func() {
				function(response)
			})
		}
	}()
}
