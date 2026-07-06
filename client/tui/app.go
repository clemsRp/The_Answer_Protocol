/* package main

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

func NewMyApp() *MyApp {
	m := &MyApp{
		app:  tview.NewApplication(),
		grid: tview.NewGrid().SetRows(0, 5).SetColumns(0),
	}

	m.app.EnableMouse(true)

	returns := createTextView("", "Responses", true, Default, Black)
	commandInput := createInputField(
        " [Saisie] ", 
        true, 
        "> ", 
        tcell.ColorGreen, 
        tcell.ColorOrange, 
        Black,
    )

    commandInput.SetDoneFunc(func(key tcell.Key) {
        if key == tcell.KeyEnter {
            // inputs <- commandInput.GetText()
            
            commandInput.SetText("") 
        }
    })

	m.grid.AddItem(returns, 0, 0, 1, 1, 0, 0, false)
	m.grid.AddItem(commandInput, 1, 0, 1, 1, 0, 0, true)

	m.app.SetRoot(m.grid, true)

	return m
}

func (m *MyApp) Run() error {
	return m.app.Run()
}
 */

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

func NewMyApp() *MyApp {
	m := &MyApp{
		app:  tview.NewApplication(),
		grid: tview.NewGrid().SetRows(0, 10).SetColumns(0),
	}

	m.app.EnableMouse(true)

	responses := createListView( "Server responses", true, Default, Default, tcell.ColorBlue, Black)
	commands := createInputField("Enter command", true, ">", Default, Default, Black)

	commands.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			inputs <- commands.GetText()
			commands.SetText("") 
		}
	})

	m.grid.AddItem(responses, 0, 0, 1, 1, 15, 0, false)
	m.grid.AddItem(commands, 1, 0, 1, 1, 15, 0, true)

	m.app.SetRoot(m.grid, true)

	go func() {
		for output := range outputs {
			responses.AddItem(output, "", 0, nil)
			m.app.Draw()
		}
	}()

	return m
}

func (m *MyApp) Run() error {
	return m.app.Run()
}
