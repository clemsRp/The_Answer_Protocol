package gui

import (
	"tap/client/gui/src/core"
	"tap/client/gui/src/drawer"
	"tap/client/gui/src/updater"
	panel "tap/client/tui/panels"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type App struct {
	*core.App
	Drawer  *drawer.Drawer
	Updater *updater.Updater
}

func NewApp(actionsChan chan panel.Action) *App {
	coreApp := core.NewApp(actionsChan)
	if coreApp == nil {
		return nil
	}

	return &App{
		App:     coreApp,
		Drawer:  drawer.NewDrawer(coreApp),
		Updater: updater.NewUpdater(coreApp),
	}
}

func (app *App) Update() {
	app.Updater.UpdateMouse()
	if app.Variables.Current_view == "Connect" {
		app.Updater.UpdateConnectView()

	} else if app.Variables.Current_view == "Game" {
		app.Updater.UpdateGameView()
	}
}

func (app *App) Draw() {
	if app.Variables.Current_view == "Connect" {
		app.Drawer.DrawConnectView()

	} else if app.Variables.Current_view == "Game" {
		app.Drawer.DrawGameView()
	}
	app.Drawer.DrawMouse()
}

func (app *App) Start() {
	for !rl.WindowShouldClose() {
		app.DrainQueue()

		app.Update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		app.Draw()
		rl.EndDrawing()
	}
}
