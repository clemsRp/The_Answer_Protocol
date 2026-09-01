package tui

import (
	"context"
	"sync"
	"tap/client/state"
	panel "tap/client/tui/panels"
	"tap/protocol"
	pr "tap/protocol"

	"github.com/rivo/tview"
)

type MyApp struct {
	pseudo      string
	app         *tview.Application
	pages       *tview.Pages
	connect     *tview.Grid
	grid        *tview.Grid
	popup       *tview.Grid
	combat      *tview.Grid
	actionsChan chan panel.Action

	Chat           *panel.ChatComponent
	Server         *panel.ServerResponseComponent
	CommandLine    *panel.CommandLineComponent
	Navigation     *panel.ChoiceListComponent
	Combat         *panel.CombatComponent
	Group          *panel.GroupComponent
	Items          *panel.ChoiceListComponent
	Interaction    *panel.ChoiceListComponent
	Datas          *panel.DatasComponent
	PopupComponent *panel.PopupComponent
	navMatrix      [4][4]tview.Primitive
	ctx            context.Context
	wg             *sync.WaitGroup

	cliVisible     bool
	popupVisible   bool
	combatVisible  bool
	selectedPerson string
}

func NewMyApp(ctx context.Context, wg *sync.WaitGroup, actionsChan chan panel.Action) *MyApp {
	tview.Styles.PrimitiveBackgroundColor = panel.AppTheme.Background
	m := &MyApp{
		app:            tview.NewApplication(),
		grid:           tview.NewGrid().SetRows(0, 0, 0, 0).SetColumns(0, 0, 0, 0),
		pages:          tview.NewPages(),
		connect:        tview.NewGrid().SetRows(-1, 5, 20).SetColumns(-1, -1, -1, -1),
		popup:          tview.NewGrid().SetRows(0, 35, 0).SetColumns(0, 60, 0),
		combat:         tview.NewGrid().SetRows(0, panel.CombatHeight, 0).SetColumns(0, panel.CombatWidth, 0),
		actionsChan:    actionsChan,
		ctx:            ctx,
		wg:             wg,
		selectedPerson: "",
	}

	m.setupComponents()
	m.setupGrid()
	m.setupMatrix()
	m.SetupFocusManager()
	m.app.EnableMouse(true)
	m.InitConnect()

	m.pages.AddPage("Connexion", m.connect, true, true)
	m.pages.AddPage("Game", m.grid, true, false)
	m.pages.AddPage("Popup", m.popup, true, false)
	m.pages.AddPage("Combat", m.combat, true, false)

	m.app.SetRoot(m.pages, true)
	return m
}

func (m *MyApp) setupComponents() {
	m.Chat = panel.NewChatComponent(m.app, m.actionsChan)
	m.CommandLine = panel.NewCommandLineComponent(m.app, m.actionsChan)
	m.Server = panel.NewServerResponseComponent(m.app)
	m.Combat = panel.NewCombatComponent(m.actionsChan, panel.CombatDatas{Current_turn: "", Leader: ""})
	m.Group = panel.NewGroupComponent(m.app, m.popup, panel.GroupDatas{}, m.actionsChan, m.OnOpenPopup, m.ShowGamePage)
	m.Navigation = panel.NewNavigationComponent(m.app, m.popup, "", map[string]string{}, m.actionsChan, m.OnOpenPopup, m.ShowGamePage)
	m.Items = panel.NewItemsComponent(m.app, m.popup, []string{}, []string{}, m.actionsChan, m.OnOpenPopup, m.ShowGamePage)
	m.Interaction = panel.NewInteractionComponent(m.app, m.popup, []string{}, []string{}, m.actionsChan, m.OnOpenPopup, m.ShowGamePage)
	m.Datas = panel.NewDatasComponent(m.app)

	m.Server.CliBtn.
		SetSelectedFunc(func() {
			if m.cliVisible {
				m.grid.RemoveItem(m.CommandLine.Layout)
				m.grid.AddItem(m.Server.Layout, 3, 0, 1, 4, 0, 0, false)

				m.navMatrix = [4][4]tview.Primitive{
					{m.Navigation.List, m.Items.List, m.Interaction.List, m.Chat.Input},
					{m.Group.Layout, m.Items.List, m.Interaction.List, m.Chat.Input},
					{m.Datas.View, m.Datas.View, m.Datas.View, m.Chat.Input},
					{m.Server.History, m.Server.History, m.Server.History, m.Server.History},
				}

				m.cliVisible = false

			} else {
				m.grid.AddItem(m.CommandLine.Layout, 3, 0, 1, 2, 0, 0, false)
				m.grid.AddItem(m.Server.Layout, 3, 2, 1, 2, 0, 0, false)

				m.navMatrix = [4][4]tview.Primitive{
					{m.Navigation.List, m.Items.List, m.Interaction.List, m.Chat.Input},
					{m.Group.Layout, m.Items.List, m.Interaction.List, m.Chat.Input},
					{m.Datas.View, m.Datas.View, m.Datas.View, m.Chat.Input},
					{m.CommandLine.Input, m.CommandLine.Input, m.Server.History, m.Server.History},
				}

				m.cliVisible = true
			}
		})

	m.Server.QuitBtn.
		SetSelectedFunc(func() {
			m.actionsChan <- panel.Action{Type: panel.ActionQuit}
		})
}

func (m *MyApp) setupMatrix() {
	m.navMatrix = [4][4]tview.Primitive{
		{m.Navigation.List, m.Items.List, m.Interaction.List, m.Chat.Input},
		{m.Group.Layout, m.Items.List, m.Interaction.List, m.Chat.Input},
		{m.Datas.View, m.Datas.View, m.Datas.View, m.Chat.Input},
		{m.Server.History, m.Server.History, m.Server.History, m.Server.History},
	}
}

func (m *MyApp) setupGrid() {
	m.grid.AddItem(m.Navigation.Layout, 0, 0, 1, 1, 0, 0, false)
	m.grid.AddItem(m.Group.Layout, 1, 0, 1, 1, 0, 0, false)
	m.grid.AddItem(m.Items.Layout, 0, 1, 2, 1, 0, 0, false)
	m.grid.AddItem(m.Interaction.Layout, 0, 2, 2, 1, 0, 0, false)
	m.grid.AddItem(m.Datas.Layout, 2, 0, 1, 3, 0, 0, false)
	m.grid.AddItem(m.Chat.Layout, 0, 3, 3, 1, 0, 0, true)
	m.grid.AddItem(m.Server.Layout, 3, 0, 1, 4, 0, 0, false)
	m.combat.AddItem(m.Combat.Layout, 1, 1, 1, 1, 0, 0, true)
}

func (m *MyApp) InitConnect() {
	input := panel.NewConnectComponent(&m.pseudo, m.actionsChan)
	logoView := panel.NewImageComponent("client/tui/assets/logo.ans")
	shopView := panel.NewImageComponent("client/tui/assets/shopfront.ans")

	m.connect.AddItem(logoView, 0, 0, 1, 5, 0, 0, false)
	m.connect.AddItem(shopView, 2, 0, 2, 5, 0, 0, false)
	m.connect.AddItem(input, 1, 2, 1, 1, 0, 0, true)
}

func (m *MyApp) OnOpenPopup(createdPopup *panel.PopupComponent) {
	m.PopupComponent = createdPopup
	m.ShowPopupPage()
}

func (m *MyApp) ShowConnectPage() {
	m.pages.SwitchToPage("Connexion")
}

func (m *MyApp) ShowGamePage() {
	m.popup.Clear()
	m.pages.SwitchToPage("Game")

	if m.popupVisible {
		panel.SetBlockedInputs(m.grid, true)
		m.popupVisible = false
	} else if m.combatVisible {
		panel.SetBlockedInputs(m.grid, true)
		m.combatVisible = false
	}

	if m.Navigation != nil && m.Navigation.List != nil {
		m.app.SetFocus(m.Navigation.List)
	}

	m.app.Sync()
}

func (m *MyApp) ShowCombatPage() {
	m.pages.ShowPage("Combat")
	m.pages.SendToFront("Combat")
	if m.Combat != nil && m.Combat.Input != nil {
		m.app.SetFocus(m.Combat.Input)
	}

	panel.SetBlockedInputs(m.grid, false)
	m.combatVisible = true
	m.app.Sync()
}

func (m *MyApp) ShowPopupPage() {
	m.pages.ShowPage("Popup")
	m.pages.SendToFront("Popup")

	panel.SetBlockedInputs(m.grid, false)

	if m.PopupComponent != nil && m.PopupComponent.FocusItem != nil {
		m.app.SetFocus(m.PopupComponent.FocusItem)
	}

	m.popupVisible = true
	m.app.Sync()
}

func (m *MyApp) ClosePopup() {
	m.ShowGamePage()
}

// Controller update callbacks

func (m *MyApp) UpdateNavigation(room *protocol.LookCommandData) {
	m.grid.RemoveItem(m.Navigation.Layout)

	roomName := ""
	opts := make(map[string]string)
	if room != nil {
		roomName = room.Name
		exits := room.Exits
		if exits.North != "" {
			opts["north"] = exits.North
		}
		if exits.East != "" {
			opts["east"] = exits.East
		}
		if exits.West != "" {
			opts["west"] = exits.West
		}
		if exits.South != "" {
			opts["south"] = exits.South
		}
	}

	m.Navigation = panel.NewNavigationComponent(
		m.app,
		m.popup,
		roomName,
		opts,
		m.actionsChan,
		m.OnOpenPopup,
		m.ShowGamePage,
	)

	m.grid.AddItem(m.Navigation.Layout, 0, 0, 1, 1, 0, 0, false)
	m.setupMatrix()
}

func (m *MyApp) UpdateItems(roomItems, inventory []string) {
	m.grid.RemoveItem(m.Items.Layout)

	m.Items = panel.NewItemsComponent(
		m.app,
		m.popup,
		roomItems,
		inventory,
		m.actionsChan,
		m.OnOpenPopup,
		m.ShowGamePage,
	)

	m.grid.AddItem(m.Items.Layout, 0, 1, 2, 1, 0, 0, false)
	m.setupMatrix()
}

func (m *MyApp) UpdateInteraction(npcs, players []string) {
	m.grid.RemoveItem(m.Interaction.Layout)

	m.Interaction = panel.NewInteractionComponent(
		m.app,
		m.popup,
		npcs,
		players,
		m.actionsChan,
		m.OnOpenPopup,
		m.ShowGamePage,
	)

	m.grid.AddItem(m.Interaction.Layout, 0, 2, 2, 1, 0, 0, false)
	m.setupMatrix()
}

func (m *MyApp) UpdateGroup(groupState state.GroupState) {
	m.grid.RemoveItem(m.Group.Layout)

	groupDatas := panel.GroupDatas{
		Group:         groupState.Group,
		LastKick:      groupState.LastKick,
		Leader:        groupState.Leader,
		Promotion:     groupState.Promotion,
		SendPromotion: groupState.SendPromotion,
		Grouped:       &groupState.Grouped,
		UnGrouped:     &groupState.UnGrouped,
		Invitations:   &groupState.Invitations,
	}

	m.Group = panel.NewGroupComponent(
		m.app,
		m.popup,
		groupDatas,
		m.actionsChan,
		m.OnOpenPopup,
		m.ShowGamePage,
	)

	m.grid.AddItem(m.Group.Layout, 1, 0, 1, 1, 0, 0, false)
	m.setupMatrix()
}

func (m *MyApp) UpdateCombat(combatState state.CombatState) {
	m.combat.Clear()

	chats := make([]panel.Chat, len(combatState.Chats))
	for i, c := range combatState.Chats {
		chats[i] = panel.Chat{Pseudo: c.Pseudo, Msg: c.Msg}
	}

	lastCombatChat := combatState.LastCombatChat
	selected := m.selectedPerson
	if selected == "" {
		selected = m.pseudo
	}

	combatDatas := panel.CombatDatas{
		Chats:            chats,
		Last_combat_chat: &lastCombatChat,
		Current_turn:     combatState.CurrentTurn,
		Leader:           combatState.Leader,
		Team:             combatState.Team,
		Opponents:        combatState.Opponents,
		SelectedPerson:   &selected,
	}

	m.Combat = panel.NewCombatComponent(m.actionsChan, combatDatas)
	m.combat.AddItem(m.Combat.Layout, 1, 1, 1, 1, 0, 0, true)
}

func (m *MyApp) UpdateDatas(text string) {
	if m.Datas != nil {
		m.Datas.SetDatas(text)
	}
}

func (m *MyApp) AppendChat(scope, user, msg string) {
	if m.Chat != nil {
		m.Chat.AppendMessage(scope, user, msg)
	}
}

func (m *MyApp) AppendCombatChat(user, msg string) {
	if m.Combat != nil {
		m.Combat.AppendChat(user, msg)
	}
}

func (m *MyApp) AppendServerResponse(res pr.ServerResponse) {
	if m.Server != nil {
		m.Server.AppendResponse(res)
	}
}

func (m *MyApp) AppendCliMessage(text string) {
	if m.CommandLine != nil {
		m.CommandLine.AppendText(text)
	}
}

func (m *MyApp) GetPseudo() string {
	return m.pseudo
}

func (m *MyApp) SetPseudo(pseudo string) {
	m.pseudo = pseudo
}

func (m *MyApp) QueueUpdate(f func()) {
	m.app.QueueUpdateDraw(f)
}

func (m *MyApp) QueueUpdateDraw(f func()) *tview.Application {
	return m.app.QueueUpdateDraw(f)
}

func (m *MyApp) Run() error {
	return m.app.Run()
}

func (m *MyApp) Stop() {
	m.app.Stop()
}

func (m *MyApp) Draw() {
	m.app.Draw()
}
