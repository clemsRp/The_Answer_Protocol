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

type GroupComponent struct {
	Layout *tview.Flex
}

type GroupDatas struct {
	Group         string
	LastKick      *string
	Leader        bool
	Promotion     bool
	SendPromotion bool
	Grouped       *[]string
	UnGrouped     *[]string
	Invitations   *[]string
}

func NewGroupComponent(
	app *tview.Application,
	popup *tview.Grid,
	groupDatas GroupDatas,
	actionsChan chan<- Action,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
) *GroupComponent {
	src := GroupComponent{
		Layout: tview.NewFlex().SetDirection(tview.FlexRow),
	}
	src.Layout.SetBorder(true).SetTitle(" Group ")

	datas := tview.NewList().
		SetSelectedBackgroundColor(tcell.GetColor("#000000")).
		SetSelectedTextColor(tcell.ColorWhite)

	cur_group := groupDatas.Group
	if cur_group == "" {
		cur_group = "None"
	}

	val_color := "yellow"

	// Display datas
	datas.AddItem(fmt.Sprintf("Group: [%s]%s", val_color, cur_group), "", 0, nil)
	datas.AddItem(fmt.Sprintf("Leader: [%s]%t", val_color, groupDatas.Leader), "", 0, nil)

	src.Layout.AddItem(datas, 0, 1, true)

	makeSpacer := func() *tview.Box {
		spacer := tview.NewBox()
		return spacer
	}

	src.Layout.AddItem(makeSpacer(), 1, 1, false)

	if groupDatas.Group == "" {
		// Create button
		create_btn := tview.NewButton("Create").SetSelectedFunc(func() {
			actionsChan <- Action{Type: ActionSendServer, Payload: pr.CmdCreateGroup}
		})
		src.Layout.AddItem(create_btn, 1, 1, false)

		if groupDatas.Invitations != nil && len(*groupDatas.Invitations) != 0 {
			// Join button
			join_func := func(invitations *[]string, inv string) {
				actionsChan <- Action{Type: ActionSendServer, Payload: pr.CmdJoinGroup + " " + inv}

				if invitations == nil {
					return
				}

				for i, v := range *invitations {
					if v == inv {
						*invitations = append((*invitations)[:i], (*invitations)[i+1:]...)
						break
					}
				}
			}
			src.Layout.AddItem(makeSpacer(), 1, 1, false)
			join_btn := CreateOptionBtn(app, popup, "Join", groupDatas.Invitations, onOpenPopup, onClosePopup, join_func)
			src.Layout.AddItem(join_btn, 1, 1, false)
		}

	} else {
		if groupDatas.Leader {
			// Invite button
			invite_func := func(users *[]string, user string) {
				actionsChan <- Action{Type: ActionSendServer, Payload: pr.CmdInviteGroup + " " + user}
				actionsChan <- Action{Type: ActionSendServer, Payload: pr.CmdUnGrouped}

				if users == nil {
					return
				}
			}

			if groupDatas.UnGrouped != nil && len(*groupDatas.UnGrouped) > 0 {
				invite_btn := CreateOptionBtn(app, popup, "Invite", groupDatas.UnGrouped, onOpenPopup, onClosePopup, invite_func)
				src.Layout.AddItem(invite_btn, 1, 1, false)
			}

			// Kick button
			kick_func := func(users *[]string, user string) {
				actionsChan <- Action{Type: ActionSendServer, Payload: pr.CmdKickGroup + " " + user}
				if groupDatas.LastKick != nil {
					*groupDatas.LastKick = user
				}

				if users == nil {
					return
				}

				for i, v := range *users {
					if v == user && i < len(*users) {
						*users = append((*users)[:i], (*users)[i+1:]...)
						break
					}
				}
			}

			if groupDatas.Grouped != nil && len(*groupDatas.Grouped) > 0 {
				kick_btn := CreateOptionBtn(app, popup, "Kick", groupDatas.Grouped, onOpenPopup, onClosePopup, kick_func)
				src.Layout.AddItem(makeSpacer(), 1, 1, false)
				src.Layout.AddItem(kick_btn, 1, 1, false)
			}

			if !groupDatas.SendPromotion && groupDatas.Grouped != nil && len(*groupDatas.Grouped) > 0 {
				// Promote button
				promote_flex := tview.NewFlex().SetDirection(tview.FlexColumn)

				selector := createSelectField("", modify_select_options(*groupDatas.Grouped), 0)
				promote_btn := tview.NewButton("Promote").SetSelectedFunc(func() {
					_, new_leader := selector.GetCurrentOption()
					actionsChan <- Action{Type: ActionSendServer, Payload: pr.CmdPromoteGroup + " " + strings.TrimSpace(new_leader)}
				})

				promote_flex.AddItem(promote_btn, 0, 1, false)
				promote_flex.AddItem(makeSpacer(), 1, 1, false)
				promote_flex.AddItem(selector, get_select_width(*groupDatas.Grouped), 1, false)

				src.Layout.AddItem(makeSpacer(), 1, 1, false)
				src.Layout.AddItem(promote_flex, 1, 1, true)
			}

		} else if groupDatas.Promotion {
			// Accept/Decline promotion
			choice_flex := tview.NewFlex().SetDirection(tview.FlexColumn)

			// Accept button
			accept_btn := tview.NewButton("✅")
			accept_btn.
				SetSelectedFunc(func() {
					actionsChan <- Action{Type: ActionSendServer, Payload: pr.CmdAcceptPromoteGroup}
					groupDatas.Promotion = false
				})
			accept_btn.SetBackgroundColor(tcell.GetColor("#000000"))
			accept_btn.SetBackgroundColorActivated(tcell.GetColor("#000000"))

			decline_btn := tview.NewButton("❌")
			decline_btn.SetSelectedFunc(func() {
				actionsChan <- Action{Type: ActionSendServer, Payload: pr.CmdDeclinePromoteGroup}
				groupDatas.Promotion = false
			})
			decline_btn.SetBackgroundColor(tcell.GetColor("#000000"))
			decline_btn.SetBackgroundColorActivated(tcell.GetColor("#000000"))

			choice_flex.AddItem(tview.NewTextView().SetText("Get promoted: ").SetTextColor(tcell.ColorYellow), 15, 0, false)
			choice_flex.AddItem(tview.NewBox(), 0, 0, false)
			choice_flex.AddItem(accept_btn, 1, 0, true)
			choice_flex.AddItem(tview.NewBox(), 2, 0, false)
			choice_flex.AddItem(decline_btn, 1, 0, true)

			src.Layout.AddItem(choice_flex, 1, 1, true)
		}

		src.Layout.AddItem(makeSpacer(), 1, 1, false)
		leave_btn := tview.NewButton("Leave").SetSelectedFunc(func() {
			actionsChan <- Action{Type: ActionSendServer, Payload: pr.CmdLeaveGroup}
		})
		src.Layout.AddItem(leave_btn, 1, 1, false)
	}

	return &src
}

func CreateOptionBtn(
	app *tview.Application,
	popupGrid *tview.Grid,
	btnLabel string,
	options *[]string,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
	onOptionSelected func(options *[]string, selected string),
) *tview.Button {

	btn := tview.NewButton(btnLabel).SetSelectedFunc(func() {
		if options == nil || len(*options) == 0 {
			return
		}

		optsList := *options
		var selectedOption string

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

		formatItem := func(name string, isSelected bool) string {
			if isSelected {
				return "[white:#7e7979]" + transform_name(name, 30)
			}
			return "[white:#474646]" + transform_name(name, 30)
		}

		// Ajout des options dans la liste
		for idx, opt := range optsList {
			isFirst := idx == 0
			actionList.AddItem(formatItem(opt, isFirst), "", 0, nil)
		}

		if len(optsList) > 0 {
			selectedOption = optsList[0]
		}

		// Mise à jour de la sélection au survol/changement d'index
		actionList.SetChangedFunc(func(i int, mainText, secondaryText string, shortcut rune) {
			if i >= 0 && i < len(optsList) {
				selectedOption = optsList[i]

				itemCount := actionList.GetItemCount()
				for idx, name := range optsList {
					if idx < itemCount {
						isSelected := (idx == i)
						actionList.SetItemText(idx, formatItem(name, isSelected), "")
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

		optionsFlex.AddItem(actionList, len(optsList)*2, 1, true)

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

		// Bouton Validate
		validateBtn := tview.NewButton("Validate").
			SetLabelColor(tcell.ColorWhite).
			SetBackgroundColorActivated(btnActiveBg).
			SetLabelColorActivated(tcell.ColorWhite).
			SetSelectedFunc(func() {
				if selectedOption == "" {
					return
				}

				chosenOpt := selectedOption

				if onClosePopup != nil {
					onClosePopup()
				}

				if onOptionSelected != nil {
					onOptionSelected(options, chosenOpt)
				}
			})
		validateBtn.SetBackgroundColor(btnRestBg)

		buttons := []*tview.Button{cancelBtn, validateBtn}
		totalHeight := (len(optsList) * 2) + 11

		newPopup := NewPopupComponent(app, popupGrid, optionsFlex, totalHeight, buttons)

		popupGrid.Clear()
		popupGrid.AddItem(newPopup.Layout, 1, 1, 1, 1, 0, 0, true)

		if onOpenPopup != nil {
			onOpenPopup(newPopup)
		}
	})

	return btn
}

func get_select_width(options []string) int {
	max_len := 0
	for _, option := range options {
		if len(option) > max_len {
			max_len = len(option)
		}
	}

	return max_len + 2
}

func modify_select_options(options []string) []string {
	res := make([]string, 0)

	for _, option := range options {
		res = append(res, fmt.Sprintf(" %s ", option))
	}

	return res
}

func (c *GroupComponent) ListenOutputs(ctx context.Context, wg *sync.WaitGroup, app *tview.Application, Chan <-chan pr.ServerResponse, function func(pr.ServerResponse)) {
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
