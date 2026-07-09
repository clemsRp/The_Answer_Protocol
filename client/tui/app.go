package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MyApp struct {
	app     *tview.Application
	pages   *tview.Pages
	connect *tview.Grid
	grid    *tview.Grid
}

var (
	Black   = tcell.NewRGBColor(0, 0, 0)
	Default = tcell.ColorDefault
)

func NewMyApp(inputs chan<- string, outputs <-chan string) *MyApp {
	tview.Styles.PrimitiveBackgroundColor = Black

	m := &MyApp{
		app:     tview.NewApplication(),
		pages:   tview.NewPages(),
		connect: tview.NewGrid().SetRows(-1, 5, 20).SetColumns(-1, -1, -1, -1),
		grid:    tview.NewGrid().SetRows(0).SetColumns(0),
	}

	m.app.EnableMouse(true)

	m.InitConnect()
	m.InitGrid()

	m.pages.AddPage("Connexion", m.connect, true, true)

	m.pages.AddPage("Game", m.grid, true, false)

	go func() {
		for output := range outputs {
			if output == "OK connected" {
				m.pages.SwitchToPage("Game")
				m.Draw()
			}
		}
	}()

	m.app.SetRoot(m.pages, true)

	return m
}

func (m *MyApp) InitGrid() {
	chatPanel := NewChatComponent(m.app, "alice")

	m.grid.AddItem(chatPanel.Layout, 0, 0, 1, 1, 0, 0, true)
}

func (m *MyApp) InitConnect() {
	input := NewConnectComponent(m)
	logoView := NewImageComponent(m, "assets/logo.ans")
	shopView := NewImageComponent(m, "assets/shopfront.ans")

	m.connect.AddItem(logoView, 0, 0, 1, 5, 0, 0, false)
	m.connect.AddItem(shopView, 2, 0, 2, 5, 0, 0, false)

	m.connect.AddItem(input, 1, 2, 1, 1, 0, 0, true)

	// m.connect.AddItem(logView, , false)
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
