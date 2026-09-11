package gui

import (
	"fmt"
	"strings"
	"sync"
	"tap/client/gui/src/parser"
	vars "tap/client/gui/src/variables"
	"tap/client/state"
	panel "tap/client/tui/panels"
	"tap/engine"
	"tap/protocol"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type App struct {
	textures     *parser.Textures
	variables    *vars.Variables
	rooms        map[string]*parser.Map
	screenWidth  int
	screenHeight int
	actionsChan  chan panel.Action
	closeOnce    sync.Once
}

func NewApp(actionsChan chan panel.Action) *App {
	rl.SetTraceLogLevel(rl.LogNone)
	rl.SetTargetFPS(60)

	rl.InitWindow(0, 0, "TAP")

	monitor := rl.GetCurrentMonitor()
	screenWidth := rl.GetMonitorWidth(monitor)
	screenHeight := rl.GetMonitorHeight(monitor)
	rl.SetWindowSize(screenWidth, screenHeight)

	app := &App{
		textures:     parser.LoadTextures(),
		variables:    vars.GetVariables(),
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		actionsChan:  actionsChan,
	}

	var err error
	maps_folder_path := "./client/gui/maps/"
	app.rooms, err = parser.ParseRooms([]string{maps_folder_path + engine.RoomEntrance})

	if err != nil {
		fmt.Println("Error parsing rooms:", err)
		app.Stop()
		return nil
	}

	app.variables.Tileset_size = float32(screenWidth / 32)
	app.variables.Zoom = float32(app.variables.Tileset_size) / float32(vars.FRAME_WIDTH)
	app.variables.StartTime = time.Now()

	return app
}

func (app *App) Update() {
	if app.variables.Current_view == "Connect" {
		app.UpdateConnectView()

	} else if app.variables.Current_view == "Game" {
		app.UpdateGameView()
	}
}

func (app *App) Draw() {
	if app.variables.Current_view == "Connect" {
		app.DrawConnectView()

	} else if app.variables.Current_view == "Game" {
		app.DrawGameView()
	}
}

func (app *App) Start() {
	for !rl.WindowShouldClose() {
		app.Update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.LightGray)
		app.Draw()
		rl.EndDrawing()
	}
}

func (app *App) QueueUpdate(f func()) {
	if f != nil {
		f()
	}
}

func (app *App) ShowConnectPage() {
	app.variables.Current_view = "Connect"
}

func (app *App) ShowGamePage() {
	app.variables.Current_view = "Game"
}

func (app *App) ShowCombatPage()                                       {}
func (app *App) ShowPopupPage()                                        {}
func (app *App) ClosePopup()                                           {}
func (app *App) ShowCombatResultPopup(result string, rewards []string) {}
func (app *App) ShowQuestCompletedPopup(questID, reward string)        {}

func (app *App) UpdateNavigation(room *protocol.LookCommandData) {
	app.variables.PanelsVariables.Room = room
	app.variables.Current_room = strings.SplitN(room.Id, "room.", 2)[1]
}

func (app *App) UpdateItems(roomItems, inventory []string) {
	app.variables.PanelsVariables.RoomItems = &roomItems
	app.variables.PanelsVariables.Inventory = &inventory
}

func (app *App) UpdateDatas(text string) {}
func (app *App) UpdateInteraction(npcs, players []string, npcData map[string]protocol.InspectNPCData, npcDialogues map[string]string, groupMembers []string, quests []protocol.TrackedQuestData, completed_quests []string) {
}

func (app *App) UpdateGroup(groupState state.GroupState) {
	app.variables.PanelsVariables.GroupState = &groupState
}

func (app *App) UpdateCombat(combatState state.CombatState) {
	app.variables.PanelsVariables.CombatState = &combatState
}

func (app *App) UpdateQuests(quests []protocol.TrackedQuestData) {
	app.variables.PanelsVariables.Quests = &quests
}

func (app *App) UpdateInspector(text string)                      {}
func (app *App) AppendChat(scope, user, msg string)               {}
func (app *App) AppendCombatChat(user, msg string)                {}
func (app *App) AppendServerResponse(res protocol.ServerResponse) {}
func (app *App) AppendCliMessage(text string)                     {}

func (app *App) AppendCliResponse(res protocol.ServerResponse) {}

func (app *App) GetPseudo() string {
	return app.variables.Player.Pseudo
}

func (app *App) SetPseudo(pseudo string) {
	app.variables.Player.Pseudo = pseudo
}

func (app *App) Stop() {
	app.closeOnce.Do(func() {
		if app.textures != nil {
			app.textures.UnloadTextures()
		}
		rl.CloseWindow()
	})
}
