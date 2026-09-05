package tui

import (
	"context"
	"fmt"
	"sync"
	"tap/client/state"
	panel "tap/client/tui/panels"
	"tap/protocol"
	pr "tap/protocol"

	"github.com/gdamore/tcell/v2"
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
	Inspector      *panel.InspectorComponent
	Quest          *panel.QuestComponent
	PopupComponent *panel.PopupComponent
	navMatrix      [4][4]tview.Primitive
	ctx            context.Context
	wg             *sync.WaitGroup

	cliVisible     bool
	popupVisible   bool
	combatVisible  bool
	selectedPerson string
	inventory      []string
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
	m.Combat = panel.NewCombatComponent(m.app, m.actionsChan, panel.CombatDatas{Current_turn: "", Leader: "", MyPseudo: ""})
	m.Group = panel.NewGroupComponent(m.app, m.popup, panel.GroupDatas{}, m.actionsChan, m.OnOpenPopup, m.ShowGamePage)
	m.Navigation = panel.NewNavigationComponent(m.app, m.popup, "", map[string]string{}, m.actionsChan, m.OnOpenPopup, m.ShowGamePage)
	m.Items = panel.NewItemsComponent(m.app, m.popup, []string{}, []string{}, m.actionsChan, m.OnOpenPopup, m.ShowGamePage)
	m.Interaction = panel.NewInteractionComponent(m.app, m.popup, []string{}, []string{}, map[string]protocol.InspectNPCData{}, map[string]string{}, []string{}, m.actionsChan, m.OnOpenPopup, m.ShowGamePage)
	m.Inspector = panel.NewInspectorComponent(m.app, m.actionsChan)
	m.Quest = panel.NewQuestComponent(m.app)

	m.Server.CliBtn.
		SetSelectedFunc(func() {
			if m.cliVisible {
				m.grid.RemoveItem(m.CommandLine.Layout)
				m.grid.AddItem(m.Server.Layout, 3, 0, 1, 4, 0, 0, false)

				m.navMatrix = [4][4]tview.Primitive{
					{m.Navigation.List, m.Items.List, m.Interaction.List, m.Chat.Input},
					{m.Group.Layout, m.Items.List, m.Interaction.List, m.Chat.Input},
					{m.Inspector.View, m.Quest.View, m.Quest.View, m.Chat.Input},
					{m.Server.History, m.Server.History, m.Server.History, m.Server.History},
				}

				m.cliVisible = false

			} else {
				m.grid.AddItem(m.CommandLine.Layout, 3, 0, 1, 2, 0, 0, false)
				m.grid.AddItem(m.Server.Layout, 3, 2, 1, 2, 0, 0, false)

				m.navMatrix = [4][4]tview.Primitive{
					{m.Navigation.List, m.Items.List, m.Interaction.List, m.Chat.Input},
					{m.Group.Layout, m.Items.List, m.Interaction.List, m.Chat.Input},
					{m.Inspector.View, m.Quest.View, m.Quest.View, m.Chat.Input},
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
		{m.Quest.View, m.Quest.View, m.Quest.View, m.Chat.Input},
		{m.Server.History, m.Server.History, m.Server.History, m.Server.History},
	}
}

func (m *MyApp) setupGrid() {
	m.grid.AddItem(m.Navigation.Layout, 0, 0, 1, 1, 0, 0, false)
	m.grid.AddItem(m.Group.Layout, 1, 0, 1, 1, 0, 0, false)
	m.grid.AddItem(m.Items.Layout, 0, 1, 2, 1, 0, 0, false)
	m.grid.AddItem(m.Interaction.Layout, 0, 2, 2, 1, 0, 0, false)
	m.grid.AddItem(m.Inspector.Layout, 2, 0, 1, 1, 0, 0, false)
	m.grid.AddItem(m.Quest.Layout, 2, 1, 1, 2, 0, 0, false)
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

	panel.SetBlockedInputs(m.grid, false)
	m.combatVisible = true

	if m.Combat != nil {
		if len(m.Combat.ActionButtons) > 0 {
			m.app.SetFocus(m.Combat.ActionButtons[0])
		} else if m.Combat.Input != nil {
			m.app.SetFocus(m.Combat.Input)
		}
	}
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

// ShowCombatResultPopup affiche une popup VICTOIRE ou DÉFAITE avec un bouton OK.
// La page Combat est masquée, puis le clic sur OK ramène à la page de jeu.
func (m *MyApp) ShowCombatResultPopup(result string, rewards []string) {
	// Masquer la page Combat et la mettre en arrière-plan
	m.pages.HidePage("Combat")
	m.pages.SendToBack("Combat")
	m.combatVisible = false
	panel.SetBlockedInputs(m.grid, true)

	// Titre et couleur selon résultat
	emote := "⚔"
	msgColor := "[green]"
	if result == "DEFEAT" {
		emote = "☠"
		msgColor = "[red]"
	}

	msg := fmt.Sprintf(" %s  %s  %s ", emote, result, emote)

	body := msgColor + msg + "[-]"

	content := tview.NewTextView().
		SetText(body).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	content.SetBackgroundColor(panel.AppTheme.PopupBackground)

	// Style identique aux boutons Cancel / Validate
	okBtn := tview.NewButton("OK").
		SetSelectedFunc(func() {
			m.QueueUpdate(func() {
				m.ShowGamePage()
			})
		})

	m.popup.Clear()
	createdPopup := panel.NewPopupComponent(m.app, m.popup, content, 8, []*tview.Button{okBtn})
	createdPopup.FocusItem = okBtn

	m.popup.AddItem(createdPopup.Layout, 1, 1, 1, 1, 0, 0, true)
	m.PopupComponent = createdPopup

	// Afficher la Popup au premier plan, Combat bien derrière
	m.pages.ShowPage("Popup")
	m.pages.SendToFront("Popup")
	panel.SetBlockedInputs(m.grid, false)
	m.app.SetFocus(okBtn)
	m.popupVisible = true
}

// ShowQuestCompletedPopup affiche une popup quand une quête est complétée.
// Après le clic sur OK, la page de jeu est restaurée.
func (m *MyApp) ShowQuestCompletedPopup(questID, reward string) {
	body := "[yellow]Quête accomplie ![-]\n\n"
	body += "[white]" + questID + "[-]\n\n"
	if reward != "" {
		body += "[green]Récompense : " + reward + "[-]"
	} else {
		body += "[gray]Aucune récompense spécifiée.[-]"
	}

	content := tview.NewTextView().
		SetText(body).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	content.SetBorder(false)
	content.SetBackgroundColor(panel.AppTheme.PopupBackground)

	// Style identique aux boutons Cancel / Validate
	okBtn := tview.NewButton("  OK  ").
		SetLabelColor(tcell.ColorWhite).
		SetBackgroundColorActivated(tcell.GetColor("#7e7979")).
		SetLabelColorActivated(tcell.ColorWhite).
		SetSelectedFunc(func() {
			m.QueueUpdate(func() {
				m.ClosePopup()
			})
		})
	okBtn.SetBackgroundColor(tcell.GetColor("#474646"))

	m.popup.Clear()
	createdPopup := panel.NewPopupComponent(m.app, m.popup, content, 8, []*tview.Button{okBtn})
	createdPopup.FocusItem = okBtn
	createdPopup.LayoutTemp.SetTitle(" ✔  Quête complétée ").SetBorder(true).SetBorderColor(panel.AppTheme.BorderActive)
	createdPopup.LayoutTemp.SetTitleColor(panel.AppTheme.TitleActive)

	m.popup.AddItem(createdPopup.Layout, 1, 1, 1, 1, 0, 0, true)
	m.PopupComponent = createdPopup
	m.ShowPopupPage()
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

	m.inventory = append([]string{}, inventory...)

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

func (m *MyApp) UpdateInteraction(npcs, players []string, npcData map[string]protocol.InspectNPCData, npcDialogue map[string]string, groupMembers []string) {
	m.grid.RemoveItem(m.Interaction.Layout)

	m.Interaction = panel.NewInteractionComponent(
		m.app,
		m.popup,
		npcs,
		players,
		npcData,
		npcDialogue,
		groupMembers,
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
	// Conserver le focus du chat si on est en train d'écrire
	focusInput := false
	if m.Combat != nil && m.app.GetFocus() == m.Combat.Input {
		focusInput = true
	}

	m.combat.Clear()

	chats := make([]panel.Chat, len(combatState.Chats))
	for i, c := range combatState.Chats {
		chats[i] = panel.Chat{Pseudo: c.Pseudo, Msg: c.Msg}
	}

	lastCombatChat := combatState.LastCombatChat

	// Synchroniser l'inventaire depuis les données de combat de l'équipe si disponible
	playerInv := make([]string, 0)
	if myData, ok := combatState.Team[m.pseudo]; ok {
		playerInv = append([]string{}, myData.Inventory...)
		m.inventory = append([]string{}, myData.Inventory...)
	} else {
		playerInv = append([]string{}, m.inventory...)
	}

	// Conserver et valider selectedPerson (par défaut cibler le premier adversaire)
	_, inOpponents := combatState.Opponents[m.selectedPerson]
	_, inTeam := combatState.Team[m.selectedPerson]
	if !inOpponents && !inTeam {
		m.selectedPerson = ""
		for name := range combatState.Opponents {
			m.selectedPerson = name
			break
		}
		if m.selectedPerson == "" {
			m.selectedPerson = m.pseudo
		}
	}

	combatDatas := panel.CombatDatas{
		Chats:            chats,
		Last_combat_chat: &lastCombatChat,
		Current_turn:     combatState.CurrentTurn,
		Leader:           combatState.Leader,
		Team:             combatState.Team,
		Opponents:        combatState.Opponents,
		SelectedPerson:   &m.selectedPerson,
		Inventory:        playerInv,
		MyPseudo:         m.pseudo,
	}

	m.Combat = panel.NewCombatComponent(m.app, m.actionsChan, combatDatas)
	m.combat.AddItem(m.Combat.Layout, 1, 1, 1, 1, 0, 0, true)

	// Gestion intelligente du focus
	if m.combatVisible {
		if focusInput {
			if m.Combat.Input != nil {
				m.app.SetFocus(m.Combat.Input)
			}
		} else if combatDatas.Current_turn != "" && combatDatas.Current_turn == m.pseudo {
			if len(m.Combat.ActionButtons) > 0 {
				m.app.SetFocus(m.Combat.ActionButtons[0])
			}
		} else {
			if m.Combat.Input != nil {
				m.app.SetFocus(m.Combat.Input)
			}
		}
	}
}

func (m *MyApp) UpdateDatas(text string) {
	// Status data was previously pushed here.
}

func (m *MyApp) UpdateQuests(quests []protocol.TrackedQuestData) {
	if m.Quest != nil {
		m.Quest.SetQuests(quests)
	}
}

func (m *MyApp) UpdateInspector(text string) {
	if m.Inspector != nil {
		m.Inspector.SetDatas(text)
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

func (m *MyApp) AppendCliResponse(res pr.ServerResponse) {
	if m.CommandLine != nil {
		m.CommandLine.AppendResponse(res)
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
	m.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyCtrlD {
			m.actionsChan <- panel.Action{
				Type:    panel.ActionSendServer,
				Payload: pr.CmdQuit,
			}
			return nil
		}

		return event
	})

	return m.app.Run()
}

func (m *MyApp) Stop() {
	m.app.Stop()
}

func (m *MyApp) Draw() {
	m.app.Draw()
}
