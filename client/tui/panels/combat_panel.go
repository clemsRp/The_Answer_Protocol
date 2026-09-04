package panel

import (
	"context"
	"fmt"
	"sync"
	pr "tap/protocol"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Chat struct {
	Pseudo string
	Msg    string
}

type CombatDatas struct {
	Chats            []Chat
	Last_combat_chat *string
	Opponents        map[string]pr.CombatPersonData
	Team             map[string]pr.CombatPersonData
	Leader           string
	Current_turn     string
	SelectedPerson   *string
	Inventory        []string
	MyPseudo         string // Ajout du pseudo pour vérifier si c'est notre tour
}

type CombatComponent struct {
	Layout  *tview.Flex
	Input   *tview.InputField
	History *tview.TextView

	OpponentsList *tview.List
	TeamList      *tview.List
	StatsList     *tview.List
	ActionButtons []*tview.Button
}

var (
	CombatHeight = 45
	CombatWidth  = 160
)

func NewCombatComponent(app *tview.Application, actionsChan chan<- Action, combat_datas CombatDatas) *CombatComponent {
	src := CombatComponent{
		Layout: tview.NewFlex().SetDirection(tview.FlexColumn),
	}
	src.Layout.SetBackgroundColor(AppTheme.PopupBackground)

	makeSpacer := func() *tview.Box {
		spacer := tview.NewBox()
		spacer.SetBackgroundColor(AppTheme.PopupBackground)
		return spacer
	}

	// Get main combat content
	main_content, main_content_width, main_content_height := src.NewMainCombatComponent(actionsChan, combat_datas)

	// Calculate padding
	padding_left := (CombatWidth - main_content_width) / 2
	padding_up := (CombatHeight - main_content_height) / 2
	if padding_up <= 2 {
		padding_up = 3
	}

	main_flex := tview.NewFlex().SetDirection(tview.FlexRow)

	title := tview.NewTextView().SetText("COMBAT").SetTextAlign(tview.AlignCenter)
	title.SetBackgroundColor(AppTheme.PopupBackground)
	mid_padding_up := (padding_up - 1) / 2

	main_flex.AddItem(makeSpacer(), mid_padding_up, 0, false)
	main_flex.AddItem(title, 1, 1, false)
	main_flex.AddItem(makeSpacer(), padding_up-mid_padding_up-1, 0, false)

	// Set paddings
	main_flex.AddItem(main_content, main_content_height, 0, true)
	main_flex.AddItem(makeSpacer(), CombatHeight-padding_up-main_content_height, 0, false)

	src.Layout.AddItem(makeSpacer(), padding_left, 0, false)
	src.Layout.AddItem(main_flex, main_content_width, 1, true)
	src.Layout.AddItem(makeSpacer(), CombatWidth-padding_left-main_content_width, 0, false)

	var focusableItems []tview.Primitive
	if src.OpponentsList != nil {
		focusableItems = append(focusableItems, src.OpponentsList)
	}
	if src.TeamList != nil {
		focusableItems = append(focusableItems, src.TeamList)
	}
	for _, b := range src.ActionButtons {
		focusableItems = append(focusableItems, b)
	}
	if src.Input != nil {
		focusableItems = append(focusableItems, src.Input)
	}

	src.Layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			cur := app.GetFocus()
			for i, item := range focusableItems {
				if cur == item {
					nextIdx := (i + 1) % len(focusableItems)
					app.SetFocus(focusableItems[nextIdx])
					return nil
				}
			}
			if len(focusableItems) > 0 {
				app.SetFocus(focusableItems[0])
			}
			return nil
		}
		if event.Key() == tcell.KeyBacktab {
			cur := app.GetFocus()
			for i, item := range focusableItems {
				if cur == item {
					prevIdx := (i - 1 + len(focusableItems)) % len(focusableItems)
					app.SetFocus(focusableItems[prevIdx])
					return nil
				}
			}
			if len(focusableItems) > 0 {
				app.SetFocus(focusableItems[len(focusableItems)-1])
			}
			return nil
		}
		return event
	})

	return &src
}

func (c *CombatComponent) NewMainCombatComponent(actionsChan chan<- Action, combat_datas CombatDatas) (tview.Primitive, int, int) {
	main_content_width := 125
	main_content_height := 35

	// Left part
	// Add elements lists
	opponents := c.GetElementList("Opponents", combat_datas, combat_datas.Opponents, nil)
	team_mates := c.GetElementList("Team", combat_datas, combat_datas.Team, nil)
	c.OpponentsList = opponents
	c.TeamList = team_mates

	top_part := SeparateElements([]tview.Primitive{opponents, team_mates}, []int{0, 0}, true)

	// Add infos and stats
	infos := c.GetInfos(map[string]string{
		"Current turn": combat_datas.Current_turn,
		"Leader":       combat_datas.Leader,
	})
	stats := c.GetStats(combat_datas)
	c.StatsList = stats
	mid_part := SeparateElements([]tview.Primitive{infos, stats}, []int{0, 0}, true)

	// Add buttons
	all_btns, btns_height := c.GetAllButtons(actionsChan, combat_datas)
	btns := tview.NewFlex().SetDirection(tview.FlexRow)
	btns.SetBackgroundColor(AppTheme.PopupBackground)

	isMyTurn := (combat_datas.Current_turn != "" && combat_datas.Current_turn == combat_datas.MyPseudo)
	if !isMyTurn {
		btns.AddItem(all_btns[0], btns_height, 0, false)
		btns_box := tview.NewBox()
		btns_box.SetBackgroundColor(AppTheme.PopupBackground)
		btns.AddItem(btns_box, 0, 1, false)
	} else {
		for ind, btn := range all_btns {
			btns.AddItem(btn, 1, 0, true)

			if ind != len(all_btns)-1 {
				box := tview.NewBox()
				box.SetBackgroundColor(AppTheme.PopupBackground)
				btns.AddItem(box, 1, 0, false)
			}
		}
		btns_box := tview.NewBox()
		btns_box.SetBackgroundColor(AppTheme.PopupBackground)
		btns.AddItem(btns_box, 0, 1, false)
	}

	left_flex := SeparateElements([]tview.Primitive{top_part, mid_part, btns}, []int{0, 0, btns_height}, false)

	// Right part
	right_flex := c.GetCombatChat(actionsChan, combat_datas)

	main_content := SeparateElements([]tview.Primitive{left_flex, right_flex}, []int{0, 0}, true)

	return main_content, main_content_width, main_content_height
}

// fillElementList (re)populates an element list, highlighting the currently
// selected person (if any) with a yellow background.
func (c *CombatComponent) fillElementList(list *tview.List, element_type string, combat_datas CombatDatas, elements map[string]pr.CombatPersonData) {
	list.Clear()

	// Add elements type
	list.AddItem(fmt.Sprintf("[yellow:%s]%s:", AppTheme.PopupBackgroundHexa, element_type), "", 0, nil)

	// Add elements
	for element := range elements {
		person := element

		background := AppTheme.PopupBackgroundHexa
		text_color := "white"
		if combat_datas.SelectedPerson != nil && *combat_datas.SelectedPerson == person {
			background = AppTheme.TextHighlightHexa
			text_color = "black"
		}

		list.AddItem(fmt.Sprintf("[%s:%s]- %s", text_color, background, person), "", 0, func() {
			if combat_datas.SelectedPerson != nil {
				*combat_datas.SelectedPerson = person
			}
			c.RefreshCombatSelection(combat_datas)
		})
	}
}

// RefreshCombatSelection redraws the opponents/team lists (to move the yellow
// highlight) and the stats panel (to show the newly selected person's stats).
func (c *CombatComponent) RefreshCombatSelection(combat_datas CombatDatas) {
	if c.OpponentsList != nil {
		c.fillElementList(c.OpponentsList, "Opponents", combat_datas, combat_datas.Opponents)
	}

	if c.TeamList != nil {
		c.fillElementList(c.TeamList, "Team", combat_datas, combat_datas.Team)
	}

	if c.StatsList != nil {
		c.fillStats(c.StatsList, combat_datas)
	}
}

func (c *CombatComponent) GetElementList(element_type string, combat_datas CombatDatas, elements map[string]pr.CombatPersonData, function func()) *tview.List {
	list := createListView("", false, false, true)
	list.SetBackgroundColor(AppTheme.PopupBackground)

	c.fillElementList(list, element_type, combat_datas, elements)

	return list
}

func (c *CombatComponent) GetInfos(datas map[string]string) *tview.List {
	list := createListView("", false, false, true)
	list.SetBackgroundColor(AppTheme.PopupBackground)

	for info_name, info := range datas {
		list.AddItem(fmt.Sprintf("[white:%s]%s: [white:blue] %s ", AppTheme.PopupBackgroundHexa, info_name, info), "", 0, nil)
	}

	return list
}

// fillStats (re)populates the stats list based on the currently selected person.
func (c *CombatComponent) fillStats(list *tview.List, combat_datas CombatDatas) {
	list.Clear()

	// Get datas
	var datas pr.CombatPersonData
	var person_color string
	if combat_datas.SelectedPerson != nil {
		for person_name, person_datas := range combat_datas.Team {
			if person_name == *combat_datas.SelectedPerson {
				person_color = "green"
				datas = person_datas
				break
			}
		}

		if person_color != "green" {
			for person_name, person_datas := range combat_datas.Opponents {
				if person_name == *combat_datas.SelectedPerson {
					person_color = "red"
					datas = person_datas
					break
				}
			}
		}

		selected := *combat_datas.SelectedPerson
		list.AddItem(fmt.Sprintf("[%s:%s]%s [white:%s]stats:", person_color, AppTheme.PopupBackgroundHexa, selected, AppTheme.PopupBackgroundHexa), "", 0, nil)
	}

	list.AddItem(fmt.Sprintf("[white:%s]Health: [yellow:%s]%d", AppTheme.PopupBackgroundHexa, AppTheme.PopupBackgroundHexa, datas.Hp), "", 0, nil)

	if len(datas.Inventory) != 0 {
		list.AddItem(fmt.Sprintf("[white:%s]Inventory:", AppTheme.PopupBackgroundHexa), "", 0, nil)

		for _, item := range datas.Inventory {
			list.AddItem(fmt.Sprintf("[white:%s]- %s", AppTheme.PopupBackgroundHexa, item), "", 0, nil)
		}
	}
}

func (c *CombatComponent) GetStats(combat_datas CombatDatas) *tview.List {
	list := createListView("", false, false, true)
	list.SetBackgroundColor(AppTheme.PopupBackground)

	c.fillStats(list, combat_datas)

	return list
}

// GetAllButtons builds the bottom action row.
// Displays waiting text if it's not our turn.
func (c *CombatComponent) GetAllButtons(actionsChan chan<- Action, combat_datas CombatDatas) ([]tview.Primitive, int) {
	// Vérification de l'état du tour
	if combat_datas.Current_turn == "" || combat_datas.Current_turn != combat_datas.MyPseudo {
		turnName := combat_datas.Current_turn
		if turnName == "" {
			turnName = "..."
		}
		msg := tview.NewTextView().
			SetDynamicColors(true).
			SetTextAlign(tview.AlignCenter)
		msg.SetText(fmt.Sprintf("\n[yellow]Waiting for [-][red]%s[-][yellow] to play...[-]", turnName))
		msg.SetBackgroundColor(AppTheme.PopupBackground)
		c.ActionButtons = nil
		return []tview.Primitive{msg}, 3
	}

	attack_btn := tview.NewButton("Attack")
	attack_btn.SetSelectedFunc(func() {
		// Evalué au moment du clic, garantit qu'on prend la cible actuellement sélectionnée
		target := ""
		if combat_datas.SelectedPerson != nil {
			if _, ok := combat_datas.Opponents[*combat_datas.SelectedPerson]; ok {
				target = *combat_datas.SelectedPerson
			}
		}
		// Fallback sur le premier ennemi si aucune cible ou allié sélectionné
		if target == "" {
			for name := range combat_datas.Opponents {
				target = name
				break
			}
		}
		if target == "" {
			return
		}
		actionsChan <- Action{
			Type:    ActionSendServer,
			Payload: fmt.Sprintf("%s %s", pr.CmdAttack, target),
		}
	})

	flee_btn := tview.NewButton("Flee")
	flee_btn.SetSelectedFunc(func() {
		actionsChan <- Action{
			Type:    ActionSendServer,
			Payload: pr.CmdFlee,
		}
	})

	buttons := []*tview.Button{attack_btn, flee_btn}

	// Regrouper les items par ID avec compteur si plusieurs exemplaires
	type itemCounter struct {
		id    string
		count int
	}
	var itemsList []itemCounter
	itemIndices := make(map[string]int)

	for _, item := range combat_datas.Inventory {
		if idx, exists := itemIndices[item]; exists {
			itemsList[idx].count++
		} else {
			itemIndices[item] = len(itemsList)
			itemsList = append(itemsList, itemCounter{id: item, count: 1})
		}
	}

	for _, entry := range itemsList {
		itemName := entry.id
		btnLabel := fmt.Sprintf("Use %s", itemName)
		if entry.count > 1 {
			btnLabel = fmt.Sprintf("Use %s (x%d)", itemName, entry.count)
		}
		use_btn := tview.NewButton(btnLabel)
		use_btn.SetSelectedFunc(func() {
			actionsChan <- Action{
				Type:    ActionSendServer,
				Payload: fmt.Sprintf("%s %s", pr.CmdUseItem, itemName),
			}
		})
		buttons = append(buttons, use_btn)
	}

	c.ActionButtons = buttons
	var res []tview.Primitive
	for _, b := range buttons {
		res = append(res, b)
	}

	return res, 2 * len(res)
}

func (c *CombatComponent) GetCombatChat(actionsChan chan<- Action, combat_datas CombatDatas) *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Chat history
	c.History = createTextView("", "", true)
	c.History.
		SetBorder(true).
		SetTitle(" CHAT ")
	c.History.SetBackgroundColor(AppTheme.PopupBackground)
	c.History.SetDynamicColors(true)

	// Add chats
	for _, chat := range combat_datas.Chats {
		fmt.Fprintf(c.History, "[yellow]%s[-]: %s\n", chat.Pseudo, chat.Msg)
	}

	// Chat input config
	c.Input = tview.NewInputField().
		SetLabel(" Message: ").
		SetFieldWidth(0)

	c.Input.SetLabelStyle(tcell.StyleDefault.
		Background(AppTheme.PopupBackground).
		Foreground(AppTheme.TextHighlight))

	c.Input.SetFieldTextColor(AppTheme.TextPrimary)

	c.Input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			text := c.Input.GetText()
			if text == "" {
				return
			}

			if combat_datas.Last_combat_chat != nil {
				*combat_datas.Last_combat_chat = text
			}

			c.Input.SetText("")
			actionsChan <- Action{
				Type:    ActionSendServer,
				Payload: fmt.Sprintf("%s %s", pr.CmdChatCombat, text),
			}
		}
	})

	flex.AddItem(c.History, 0, 1, false)
	flex.AddItem(c.Input, 1, 0, true)

	return flex
}

func (c *CombatComponent) AppendChat(user, msg string) {
	if c.History != nil {
		fmt.Fprintf(c.History, "[yellow]%s[-]: %s\n", user, msg)
	}
}

func (c *CombatComponent) ListenOutputs(ctx context.Context, wg *sync.WaitGroup, app *tview.Application, Chan <-chan pr.ServerResponse, function func(pr.ServerResponse)) {
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
