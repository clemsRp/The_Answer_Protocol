package tui

import (
	panel "tap/client/tui/panels"

	"github.com/rivo/tview"
)

type MyApp struct {
	app     *tview.Application
	pages   *tview.Pages
	connect *tview.Grid
	grid    *tview.Grid
	popup   *tview.Grid
	router  *Router

	Chat           *panel.ChatComponent
	Server         *panel.ServerResponseComponent
	CommandLine    *panel.CommandLineComponent
	Navigation     *panel.ChoiceListComponent
	Items          *panel.ChoiceListComponent
	Interaction    *panel.ChoiceListComponent
	Datas          *panel.DatasComponent
	PopupComponent *panel.PopupComponent
	navMatrix      [4][4]tview.Primitive
}

var (
	cli_visible   = false
	popup_visible = false
)

func NewMyApp(router *Router) *MyApp {
	tview.Styles.PrimitiveBackgroundColor = panel.Black

	m := &MyApp{
		app:     tview.NewApplication(),
		grid:    tview.NewGrid().SetRows(0, 0, 0, 0).SetColumns(0, 0, 0, 0),
		pages:   tview.NewPages(),
		connect: tview.NewGrid().SetRows(-1, 5, 20).SetColumns(-1, -1, -1, -1),
		popup:   tview.NewGrid().SetRows(0, 35, 0).SetColumns(0, 60, 0),
		router:  router,
	}

	m.connect.SetBackgroundColor(panel.Black)
	m.popup.SetBackgroundColor(panel.Black)

	m.setupComponents(router)
	m.setupGrid()
	m.setupMatrix()
	m.StartListeners()
	m.SetupFocusManager()
	m.app.EnableMouse(true)
	m.InitConnect()

	m.pages.AddPage("Connexion", m.connect, true, true)
	m.pages.AddPage("Game", m.grid, true, false)
	m.pages.AddPage("Popup", m.popup, true, false)

	m.app.SetRoot(m.pages, true)
	return m
}

func (m *MyApp) setupComponents(router *Router) {
	m.Chat = panel.NewChatComponent(m.app, router.Inputs)
	m.CommandLine = panel.NewCommandLineComponent(m.app, router.Inputs)
	m.Server = panel.NewServerResponseComponent(m.app)
	m.Navigation = panel.NewNavigationComponent(m.app, m.popup, router.Inputs, m.OnOpenPopup, m.ShowGamePage)
	m.Items = panel.NewItemsComponent(m.app, m.popup, router.Inputs, m.OnOpenPopup, m.ShowGamePage)
	m.Interaction = panel.NewInteractionComponent(m.app, m.popup, router.Inputs, m.OnOpenPopup, m.ShowGamePage)
	m.Datas = panel.NewDatasComponent(m.app)

	m.Server.CliBtn.
		SetSelectedFunc(func() {
			if cli_visible {
				m.grid.RemoveItem(m.CommandLine.Layout)
				m.grid.AddItem(m.Server.Layout, 3, 0, 1, 4, 0, 0, false)

				m.navMatrix = [4][4]tview.Primitive{
					{m.Navigation.List, m.Items.List, m.Interaction.List, m.Chat.Input},
					{m.Navigation.List, m.Items.List, m.Interaction.List, m.Chat.Input},
					{m.Datas.View, m.Datas.View, m.Datas.View, m.Chat.Input},
					{m.Server.History, m.Server.History, m.Server.History, m.Server.History},
				}

				cli_visible = false

			} else {
				m.grid.AddItem(m.CommandLine.Layout, 3, 0, 1, 2, 0, 0, false)
				m.grid.AddItem(m.Server.Layout, 3, 2, 1, 2, 0, 0, false)

				m.navMatrix = [4][4]tview.Primitive{
					{m.Navigation.List, m.Items.List, m.Interaction.List, m.Chat.Input},
					{m.Navigation.List, m.Items.List, m.Interaction.List, m.Chat.Input},
					{m.Datas.View, m.Datas.View, m.Datas.View, m.Chat.Input},
					{m.CommandLine.Input, m.CommandLine.Input, m.Server.History, m.Server.History},
				}

				cli_visible = true
			}
		})

	m.Server.QuitBtn.
		SetSelectedFunc(func() {
			router.Inputs <- "QUIT"
		})
}

func (m *MyApp) setupMatrix() {
	m.navMatrix = [4][4]tview.Primitive{
		{m.Navigation.List, m.Items.List, m.Interaction.List, m.Chat.Input},
		{m.Navigation.List, m.Items.List, m.Interaction.List, m.Chat.Input},
		{m.Datas.View, m.Datas.View, m.Datas.View, m.Chat.Input},
		{m.Server.History, m.Server.History, m.Server.History, m.Server.History},
	}
}

func (m *MyApp) setupGrid() {
	m.grid.AddItem(m.Navigation.Layout, 0, 0, 2, 1, 0, 0, false)
	m.grid.AddItem(m.Items.Layout, 0, 1, 2, 1, 0, 0, false)
	m.grid.AddItem(m.Interaction.Layout, 0, 2, 2, 1, 0, 0, false)
	m.grid.AddItem(m.Datas.Layout, 2, 0, 1, 3, 0, 0, false)
	m.grid.AddItem(m.Chat.Layout, 0, 3, 3, 1, 0, 0, true)
	m.grid.AddItem(m.Server.Layout, 3, 0, 1, 4, 0, 0, false)
}

func (m *MyApp) StartListeners() {
	m.Chat.ListenOutputs(m.app, m.router.ChatChan)
	m.CommandLine.ListenOutputs(m.app, m.router.CommandLineChan)
	m.Server.ListenOutputs(m.app, m.router.ServerChan)
	m.Navigation.ListenOutputs(m.app, m.router.NavChan, m.NavListenOutputs)
	m.Items.ListenOutputs(m.app, m.router.ItemsChan, m.ItemListenOutputs)
	m.Interaction.ListenOutputs(m.app, m.router.InteractionChan, m.InteractionListenOutputs)
	m.Datas.ListenOutputs(m.app, m.router.DatasChan)
}

func (m *MyApp) InitConnect() {
	input := panel.NewConnectComponent(m.router.Inputs)
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

func (m *MyApp) ShowGamePage() {
	m.popup.Clear()
	m.pages.SwitchToPage("Game")

	if popup_visible {
		panel.SetBlockedInputs(m.grid, true)
		popup_visible = false
	}

	m.app.SetFocus(m.Navigation.List)

	go m.app.QueueUpdateDraw(func() {
		m.app.Sync()
	})
}

func (m *MyApp) ShowPopupPage() {
	m.pages.ShowPage("Popup")
	m.pages.SendToFront("Popup")

	panel.SetBlockedInputs(m.grid, false)

	if m.PopupComponent != nil && m.PopupComponent.FocusItem != nil {
		m.app.SetFocus(m.PopupComponent.FocusItem)
	}

	popup_visible = true

	go m.app.QueueUpdateDraw(func() {
		m.app.Sync()
	})
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
