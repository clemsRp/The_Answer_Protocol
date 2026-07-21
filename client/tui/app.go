package main

import (
	panel "tap/client/tui/panels"

	"github.com/rivo/tview"
)

type MyApp struct {
	app     *tview.Application
	pages   *tview.Pages
	connect *tview.Grid
	grid    *tview.Grid
	router  *Router

	Chat        *panel.ChatComponent
	Server      *panel.ServerResponseComponent
	Navigation  *panel.NavigationComponent
	PlayersNPC  *panel.PlayersNPCComponent
	Dialogue    *panel.DialogueComponent
	Inventory   *panel.InventoryComponent
	Quest       *panel.QuestComponent
	ItemsInRoom *panel.ItemsComponent
	Info        *panel.InfoComponent
	navMatrix   [4][4]tview.Primitive
}

func NewMyApp(router *Router) *MyApp {
	tview.Styles.PrimitiveBackgroundColor = panel.Black

	m := &MyApp{
		app:     tview.NewApplication(),
		grid:    tview.NewGrid().SetRows(0, 0, 0, 0).SetColumns(0, 0, 0, 0),
		pages:   tview.NewPages(),
		connect: tview.NewGrid().SetRows(-1, 5, 20).SetColumns(-1, -1, -1, -1),
		router:  router,
	}

	m.connect.SetBackgroundColor(panel.Black)

	m.setupComponents(router)
	m.setupGrid()
	m.setupMatrix()
	m.StartListeners()
	m.SetupFocusManager()
	m.app.EnableMouse(true)
	m.InitConnect()
	m.pages.AddPage("Connexion", m.connect, true, true)

	m.pages.AddPage("Game", m.grid, true, false)

	m.app.SetRoot(m.pages, true)
	return m
}

func (m *MyApp) setupComponents(router *Router) {
	m.Chat = panel.NewChatComponent(m.app, "Player1", router.Inputs)
	m.Server = panel.NewServerResponseComponent(m.app)
	m.Navigation = panel.NewNavigationComponent(m.app)
	m.PlayersNPC = panel.NewPlayersNPCComponent(m.app)
	m.Dialogue = panel.NewDialogueComponent(m.app)
	m.Inventory = panel.NewInventoryComponent(m.app)
	m.Quest = panel.NewQuestComponent(m.app)
	m.ItemsInRoom = panel.NewItemsRoomComponent(m.app)
	m.Navigation = panel.NewNavigationComponent(m.app)
	m.Info = panel.NewInfoComponent(m.app)
}

func (m *MyApp) setupMatrix() {
	m.navMatrix = [4][4]tview.Primitive{
		{m.Navigation.Navigation, m.PlayersNPC.List, m.Dialogue.View, m.Chat.Input},
		{m.ItemsInRoom.Items, m.PlayersNPC.List, m.Dialogue.View, m.Chat.Input},
		{m.Info.View, m.Quest.List, m.Inventory.List, m.Chat.Input},
		{m.Server.History, m.Server.History, m.Server.History, m.Chat.Input},
	}
}

func (m *MyApp) setupGrid() {
	m.grid.AddItem(m.Navigation.Layout, 0, 0, 1, 1, 0, 0, false)
	m.grid.AddItem(m.ItemsInRoom.Layout, 1, 0, 1, 1, 0, 0, false)
	m.grid.AddItem(m.PlayersNPC.Layout, 0, 1, 2, 1, 0, 0, false)
	m.grid.AddItem(m.Quest.Layout, 2, 1, 1, 1, 0, 0, false)
	m.grid.AddItem(m.Dialogue.Layout, 0, 2, 2, 1, 0, 0, false)
	m.grid.AddItem(m.Info.Layout, 2, 0, 1, 1, 0, 0, false)
	m.grid.AddItem(m.Inventory.Layout, 2, 2, 1, 1, 0, 0, false)
	m.grid.AddItem(m.Chat.Layout, 0, 3, 4, 1, 0, 0, true)
	m.grid.AddItem(m.Server.Layout, 3, 0, 1, 3, 0, 0, false)
}

func (m *MyApp) StartListeners() {
	m.Chat.ListenOutputs(m.app, m.router.ChatChan)
	m.Server.ListenOutputs(m.app, m.router.ServerChan)
	m.Navigation.ListenOutputs(m.app, m.router.NavChan, m.router.Inputs)
	m.PlayersNPC.ListenOutputs(m.app, m.router.PlayersChan, m.router.Inputs)
	m.Dialogue.ListenOutputs(m.app, m.router.DialogueChan)
	m.Inventory.ListenOutputs(m.app, m.router.InventoryChan, m.router.Inputs)
	m.Quest.ListenOutputs(m.app, m.router.QuestChan, m.router.Inputs)
	m.ItemsInRoom.ListenOutputs(m.app, m.router.ItemsChan, m.router.Inputs)
	m.Navigation.ListenOutputs(m.app, m.router.NavChan, m.router.Inputs)
}

func (m *MyApp) InitConnect() {
	input := panel.NewConnectComponent(inputs)
	logoView := panel.NewImageComponent("assets/logo.ans")
	shopView := panel.NewImageComponent("assets/shopfront.ans")

	m.connect.AddItem(logoView, 0, 0, 1, 5, 0, 0, false)
	m.connect.AddItem(shopView, 2, 0, 2, 5, 0, 0, false)

	m.connect.AddItem(input, 1, 2, 1, 1, 0, 0, true)
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
