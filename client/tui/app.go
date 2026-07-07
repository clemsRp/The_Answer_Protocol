package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MyApp struct {
	app  *tview.Application
	grid *tview.Grid
}

var (
	Black   = tcell.NewRGBColor(0, 0, 0)
	Default = tcell.ColorDefault
)

func NewMyApp(inputs chan<- string, outputs <-chan string) *MyApp {
	m := &MyApp{
		app:  tview.NewApplication(),
		grid: tview.NewGrid().SetRows(0, 10).SetColumns(0),
	}
	chat_component := NewChatComponent(m.app, "Player1", inputs)
	chat_component.ListenOutputs(m.app, outputs)
	server_response_component := NewServerResponseComponent(m.app)
	server_response_component.ListenOutputs(m.app, outputs)
	m.grid.AddItem(chat_component.Layout, 0, 3, 3, 1, 0, 0, false)
	m.grid.AddItem(server_response_component.Layout, 3, 0, 0, 3, 0, 0, false)

	m.app.EnableMouse(true)

	// responses := createListView("Server responses", true, Default, Default, tcell.ColorBlue, Black)
	// commands := createInputField("Enter command", true, ">", Default, Default, Black)

	// commands.SetDoneFunc(func(key tcell.Key) {
	// 	if key == tcell.KeyEnter {
	// 		inputs <- commands.GetText()
	// 		commands.SetText("")
	// 	}
	// })

	// m.grid.AddItem(responses, 0, 0, 1, 1, 15, 0, false)
	// m.grid.AddItem(commands, 1, 0, 1, 1, 15, 0, true)

	m.app.SetRoot(m.grid, true)

	return m
}

func (m *MyApp) Run() error {
	return m.app.Run()
}

func (m *MyApp) Stop() {
	m.app.Stop()
}
