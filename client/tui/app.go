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
		connect: tview.NewGrid().SetRows(0, 5).SetColumns(32, 0), // Modifié ici
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
	imgView, logView := NewImageComponent(m)

	// L'image prend toute la ligne du haut (Ligne 0, s'étale sur 2 colonnes)
	m.connect.AddItem(imgView, 0, 0, 1, 2, 0, 0, false)

	// L'input va en bas à gauche (Ligne 1, Colonne 0)
	m.connect.AddItem(input, 1, 0, 1, 1, 0, 0, true)

	// Le Log va en bas à droite (Ligne 1, Colonne 1)
	m.connect.AddItem(logView, 1, 1, 1, 1, 0, 0, false)
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
