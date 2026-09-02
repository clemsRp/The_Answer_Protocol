package gui

import (
	"fmt"
	"sync"
	"tap/client/gui/src/parser"
	vars "tap/client/gui/src/variables"
	"tap/client/state"
	panel "tap/client/tui/panels"
	"tap/protocol"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type App struct {
	textures    *parser.Textures
	variables   *vars.Variables
	rooms       map[string]*parser.Map
	actionsChan chan panel.Action
	closeOnce   sync.Once
}

func NewApp(actionsChan chan panel.Action) *App {
	rl.SetTraceLogLevel(rl.LogNone)
	rl.SetTargetFPS(60)

	rl.InitWindow(0, 0, "TAP: OAP")

	monitor := rl.GetCurrentMonitor()
	screenWidth := rl.GetMonitorWidth(monitor)
	screenHeight := rl.GetMonitorHeight(monitor)
	rl.SetWindowSize(screenWidth, screenHeight)

	app := &App{
		textures:    parser.LoadTextures(),
		variables:   vars.GetVariables(),
		actionsChan: actionsChan,
	}

	var err error
	app.rooms, err = parser.ParseRooms([]string{})
	if err != nil {
		fmt.Println("Error parsing rooms:", err)
		app.Stop()
		return nil
	}

	app.variables.Tileset_size = screenWidth / 30

	return app
}

func (app *App) Update() {
	// Update player coordinates
	if rl.IsKeyPressed(rl.KeyEnter) {
		app.actionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: "CONNECT ali"}
	}
	if rl.IsKeyDown(rl.KeyRight) {
		app.variables.X += 4
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		app.variables.X -= 4
	}
	if rl.IsKeyDown(rl.KeyDown) {
		app.variables.Y += 4
	}
	if rl.IsKeyDown(rl.KeyUp) {
		app.variables.Y -= 4
	}
}

func (app *App) Draw() {

	/* app.DrawMap("entrance_floor.json")
	app.DrawMap("entrance_furnitures.json") */
}

func (app *App) DrawMap(map_name string) {
	/* cur_room := app.rooms[app.variables.Current_room]
	cur_tilesets := cur_room.Tilesets

	for y := range len(cur_tilesets) {
		for x := range len(cur_tilesets[y]) {
			// Get tileset datas
			// texture_name, indX, indY := get_tileset_datas(cur_tilesets[y][x])
			texture_name, indX, indY := "/home/crappo/Documents/M5/The_Answer_Protocol/pre-v1/client/gui/assets/Objects/Basic Furniture.png", 0, 0

			// Calculate tileset coordinates
			coor_x := x * app.variables.Tileset_size
			coor_y := y * app.variables.Tileset_size

			parser.DrawImage(
				(*app.textures)[texture_name],
				float32(coor_x), float32(coor_y),
				float32(indX), float32(indY),
				1, 1, float32(app.variables.Tileset_size/vars.FRAME_WIDTH),
			)
		}
	} */
}

func (app *App) Start() {
	for !rl.WindowShouldClose() {
		// Update game state
		app.Update()

		// Draw game
		rl.BeginDrawing()
		rl.ClearBackground(rl.LightGray)
		app.Draw()
		rl.EndDrawing()
	}
}

func (app *App) QueueUpdate(f func())                             {}
func (app *App) ShowConnectPage()                                 {}
func (app *App) ShowGamePage()                                    {}
func (app *App) ShowCombatPage()                                  {}
func (app *App) ShowPopupPage()                                   {}
func (app *App) ClosePopup()                                      {}
func (app *App) UpdateNavigation(room *protocol.LookCommandData)  {}
func (app *App) UpdateItems(roomItems, inventory []string)        {}
func (app *App) UpdateInteraction(npcs, players []string)         {}
func (app *App) UpdateGroup(groupState state.GroupState)          {}
func (app *App) UpdateCombat(combatState state.CombatState)       {}
func (app *App) UpdateDatas(text string)                          {}
func (app *App) AppendChat(scope, user, msg string)               {}
func (app *App) AppendCombatChat(user, msg string)                {}
func (app *App) AppendServerResponse(res protocol.ServerResponse) {}
func (app *App) AppendCliMessage(text string)                     {}
func (app *App) AppendCliResponse(res protocol.ServerResponse)    {}
func (app *App) GetPseudo() string                                { return "" }
func (app *App) SetPseudo(pseudo string)                          {}
func (app *App) Stop() {
	app.closeOnce.Do(func() {
		if app.textures != nil {
			app.textures.UnloadTextures()
		}
		rl.CloseWindow()
	})
}
