package gui

import (
	"tap/src/client/gui/src/core"
	"tap/src/client/gui/src/drawer"
	"tap/src/client/gui/src/updater"
	panel "tap/src/client/tui/panels"

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
	// Handle quit
	if rl.IsKeyPressed(rl.KeyEscape) {
		app.Variables.PanelsVariables.Datas.Open = false
		app.Variables.PanelsVariables.Inspect.Open = false
		app.Variables.PanelsVariables.Group.Open = false
		app.Variables.PanelsVariables.Chat.Open = false
	}

	app.Updater.UpdateMouse()
	if app.Variables.Current_view == "Connect" {
		app.Updater.UpdateConnectView()

	} else if app.Variables.Current_view == "Game" {
		app.Updater.UpdateGameView()

	} else if app.Variables.Current_view == "Combat" {
		app.Updater.UpdateCombatView()

	} else if app.Variables.Current_view == "CombatResult" {
		app.Updater.UpdateCombatResultView()
	}

	app.CheckPlayerMovement()
}

func (app *App) Draw() {
	app.Drawer.RenderChatTexture()
	app.Drawer.RenderGroupTexture()
	if app.Variables.Current_view == "Connect" {
		app.Drawer.DrawConnectView()

	} else if app.Variables.Current_view == "Game" {
		app.Drawer.DrawGameView()

	} else if app.Variables.Current_view == "Combat" {
		app.Drawer.DrawCombatView()

	} else if app.Variables.Current_view == "CombatResult" {
		app.Drawer.DrawCombatResultView()
	}

	app.Drawer.DrawMouse()
}

func (app *App) Start() {
	for !rl.WindowShouldClose() && app.Running {
		app.DrainQueue()

		app.Update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		app.Draw()
		rl.EndDrawing()
	}

	app.Textures.UnloadTextures()
	rl.CloseWindow()
}

func (app *App) CheckPlayerMovement() {
	current_room := app.Rooms[app.Variables.Current_room]
	if current_room == nil || len(current_room.Collisions) == 0 || len(current_room.Collisions[0]) == 0 {
		return
	}

	curX := app.Variables.Player.Position.X
	curY := app.Variables.Player.Position.Y
	tile := app.Variables.Tileset_size

	inside := curX >= 2*tile && curX <= 29*tile && curY >= 2*tile && curY <= 15*tile

	if !app.Updater.CanMove(current_room, curX, curY, int(tile)) && inside {
		newX, newY, found := app.FindNearestFreeTile(
			app.Variables.Player.Position.X,
			app.Variables.Player.Position.Y,
		)
		if found {
			app.Variables.Player.Position.X = newX
			app.Variables.Player.Position.Y = newY
		}
	}
}
