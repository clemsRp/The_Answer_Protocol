package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MyApp struct {
	app    *tview.Application
	grid   *tview.Grid
	router *Router

	Chat        *ChatComponent
	Server      *ServerResponseComponent
	Navigation  *NavigationComponent
	PlayersNPC  *PlayersNPCComponent
	Dialogue    *DialogueComponent
	Inventory   *InventoryComponent
	Quest       *QuestComponent
	ItemsInRoom *ItemsComponent
	Info        *InfoComponent
	navMatrix   [4][4]tview.Primitive
}

var (
	Black   = tcell.NewRGBColor(0, 0, 0)
	Default = tcell.ColorDefault
)

func NewMyApp(router *Router) *MyApp {
	m := &MyApp{
		app:    tview.NewApplication(),
		grid:   tview.NewGrid().SetRows(0, 0, 0, 0).SetColumns(0, 0, 0, 0),
		router: router,
	}

	m.Chat = NewChatComponent(m.app, "Player1", router.Inputs)
	m.Server = NewServerResponseComponent(m.app)
	m.Navigation = NewNavigationComponent(m.app)
	m.PlayersNPC = NewPlayersNPCComponent(m.app)
	m.Dialogue = NewDialogueComponent(m.app)
	m.Inventory = NewInventoryComponent(m.app)
	m.Quest = NewQuestComponent(m.app)
	m.ItemsInRoom = NewItemsRoomComponent(m.app)
	m.Navigation = NewNavigationComponent(m.app)
	m.Info = NewInfoComponent(m.app)
	m.setupGrid()
	m.navMatrix = [4][4]tview.Primitive{
		{m.Navigation.Navigation, m.PlayersNPC.List, m.Dialogue.View, m.Chat.Input},
		{m.ItemsInRoom.Items, m.PlayersNPC.List, m.Dialogue.View, m.Chat.Input},
		{m.Info.View, m.Quest.List, m.Inventory.List, m.Chat.Input},
		{m.Server.History, m.Server.History, m.Server.History, m.Chat.Input},
	}
	m.SetupFocusManager()

	m.app.EnableMouse(true)
	m.app.SetRoot(m.grid, true)
	return m
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

func (m *MyApp) moveFocusSpatial(dRow, dCol int) {
	currentFocus := m.app.GetFocus()
	startRow, startCol := -1, -1

	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if m.navMatrix[r][c] == currentFocus {
				startRow, startCol = r, c
				break
			}
		}
	}

	if startRow == -1 {
		return
	}

	targetRow, targetCol := startRow, startCol

	for {
		targetRow += dRow
		targetCol += dCol

		if targetRow < 0 || targetRow >= 4 || targetCol < 0 || targetCol >= 4 {
			return
		}

		nextComponent := m.navMatrix[targetRow][targetCol]

		if nextComponent != currentFocus && nextComponent != nil {
			m.app.SetFocus(nextComponent)
			return
		}
	}
}

func (m *MyApp) SetupFocusManager() {

	m.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == tcell.ModShift {
			switch event.Key() {

			case tcell.KeyUp:
				m.moveFocusSpatial(-1, 0)
				return nil

			case tcell.KeyDown:
				m.moveFocusSpatial(1, 0)
				return nil

			case tcell.KeyLeft:
				m.moveFocusSpatial(0, -1)
				return nil

			case tcell.KeyRight:
				m.moveFocusSpatial(0, 1)
				return nil
			}
		}

		return event
	})
}

func (m *MyApp) Run() error {
	return m.app.Run()
}

func (m *MyApp) Stop() {
	m.app.Stop()
}
