package core

import (
	"encoding/json"
	"fmt"
	"os"
	"tap/client/gui/src/parser"
	"tap/client/gui/src/ui"
	vars "tap/client/gui/src/variables"
	panel "tap/client/tui/panels"
	"tap/engine"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type App struct {
	Textures     *parser.Textures
	Variables    *vars.Variables
	Colors       map[string]rl.Color
	Manager      *ui.Manager
	Rooms        map[string]*parser.Map
	ScreenWidth  int
	ScreenHeight int
	ActionsChan  chan panel.Action
	Running      bool

	updateQueue chan func()
}

func NewApp(actionsChan chan panel.Action) *App {
	rl.SetTraceLogLevel(rl.LogNone)
	rl.SetTargetFPS(60)

	rl.InitWindow(0, 0, "TAP")
	rl.HideCursor()
	rl.SetExitKey(0)

	monitor := rl.GetCurrentMonitor()
	screenWidth := rl.GetMonitorWidth(monitor)
	screenHeight := rl.GetMonitorHeight(monitor)

	rl.SetWindowSize(screenWidth, screenHeight)

	// Init App
	app := &App{
		Textures:     parser.LoadTextures(),
		Variables:    vars.GetVariables(),
		Colors:       vars.GetColors(),
		Manager:      ui.NewManager(),
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		ActionsChan:  actionsChan,
		updateQueue:  make(chan func(), 256),
		Running:      true,
	}

	app.ParseNpcNames("./world.json")

	// Parse maps
	var err error
	maps_folder_path := "./client/gui/maps/"
	app.Rooms, err = parser.ParseRooms(
		[]string{
			maps_folder_path + engine.RoomEntrance,
			maps_folder_path + engine.RoomVillageSquare,
			maps_folder_path + engine.RoomAbandonedFarm,
			maps_folder_path + engine.RoomMerchantTent,
			maps_folder_path + engine.RoomOldBarn,
			// maps_folder_path + engine.RoomChickenCoop,
			maps_folder_path + engine.RoomNorthBridge,
			// maps_folder_path + engine.RoomWindMill,
			// maps_folder_path + engine.RoomForestEdge,
			// maps_folder_path + engine.RoomRiverside,
			// maps_folder_path + engine.RoomFishingDock,
			// maps_folder_path + engine.RoomOrchard,
		},
	)

	// Handle error
	if err != nil {
		fmt.Println("Error parsing rooms:", err)
		app.Stop()
		return nil
	}

	app.AddMissingVariables(screenWidth)
	app.AddItemPositions()

	return app
}

func (app *App) ParseNpcNames(filepath string) {
	// Parse json file
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	var world []engine.Map
	if err := json.Unmarshal(data, &world); err != nil {
		fmt.Printf("JSON parsing error: %v\n", err)
	}

	// Create NpcConvertor
	app.Variables.NpcConvertor = make(map[string]string)
	for npc_id, npc := range world[0].Npcs {
		app.Variables.NpcConvertor[npc.Name] = npc_id
	}
}

func (app *App) AddMissingVariables(screenWidth int) {
	app.Variables.Tileset_size = float32(screenWidth / 32)
	app.Variables.FontSize = 0.4 * app.Variables.Tileset_size
	app.Variables.Player.Speed = int(0.12 * app.Variables.Tileset_size)
	app.Variables.Zoom = float32(app.Variables.Tileset_size) / float32(vars.FRAME_WIDTH)
	app.Variables.StartTime = time.Now()
	app.Variables.Player.LastTimeTyped = time.Now()
	app.Variables.StartingPosX = float32(15.5 * app.Variables.Tileset_size)
	app.Variables.StartingPosY = float32(8.5 * app.Variables.Tileset_size)
	app.Variables.Player.Position = &vars.Position{
		X: app.Variables.StartingPosX,
		Y: app.Variables.StartingPosY,
	}
	app.Variables.PanelsVariables.Chat.Rect = rl.NewRectangle(
		vars.CHAT_START_X*app.Variables.Tileset_size,
		vars.CHAT_START_Y*app.Variables.Tileset_size,
		vars.CHAT_WIDTH*app.Variables.Tileset_size,
		vars.CHAT_HEIGHT*app.Variables.Tileset_size,
	)
	app.Variables.PanelsVariables.Chat.ScrollRect = rl.NewRectangle(
		(vars.CHAT_START_X+vars.CHAT_WIDTH-0.6)*app.Variables.Tileset_size,
		(vars.CHAT_START_Y+0.25)*app.Variables.Tileset_size,
		0.25*app.Variables.Tileset_size,
		(vars.CHAT_HEIGHT-0.5)*app.Variables.Tileset_size,
	)

	app.Variables.PanelsVariables.Group.Rect = rl.NewRectangle(
		vars.GROUP_START_X*app.Variables.Tileset_size,
		vars.GROUP_START_Y*app.Variables.Tileset_size,
		vars.GROUP_WIDTH*app.Variables.Tileset_size,
		vars.GROUP_HEIGHT*app.Variables.Tileset_size,
	)
	app.Variables.PanelsVariables.Group.ScrollRect = rl.NewRectangle(
		(vars.GROUP_START_X+vars.GROUP_WIDTH-0.6)*app.Variables.Tileset_size,
		(vars.GROUP_START_Y+0.25)*app.Variables.Tileset_size,
		0.25*app.Variables.Tileset_size,
		(vars.GROUP_HEIGHT-0.5)*app.Variables.Tileset_size,
	)

	app.Variables.PanelsVariables.Talk.Rect = rl.NewRectangle(
		float32(vars.TALK_START_X)*app.Variables.Tileset_size,
		float32(vars.TALK_START_Y)*app.Variables.Tileset_size,
		float32(vars.TALK_WIDTH)*app.Variables.Tileset_size,
		float32(vars.TALK_HEIGHT)*app.Variables.Tileset_size,
	)

	app.Variables.Player.Zoom = app.Variables.Zoom

	remote_players := make(map[string]*vars.Player)
	app.Variables.RemotePlayers = &remote_players
}

func (app *App) AddItemPositions() {
	positions := make(map[string]*vars.Position)

	// Set positions
	// TODO Define all items
	positions["turnip_seed"] = &vars.Position{X: 4, Y: 11}
	positions["rusty_hoe"] = &vars.Position{X: 4, Y: 9}

	// Scale positions to map size
	for _, pos := range positions {
		(*pos).X *= app.Variables.Tileset_size
		(*pos).Y *= app.Variables.Tileset_size
	}

	app.Variables.ItemPositions = &positions
}

func (app *App) GroupContentStartY() float32 {
	tile := app.Variables.Tileset_size
	font := app.Variables.FontSize
	return (vars.GROUP_START_Y+2)*tile + 2*font
}

func (app *App) QueueUpdate(f func()) {
	if f == nil {
		return
	}
	select {
	case app.updateQueue <- f:
	default:
		fmt.Println("GUI: update queue full, dropping app UI update")
	}
}

func (app *App) DrainQueue() {
	for {
		select {
		case f := <-app.updateQueue:
			f()
		default:
			return
		}
	}
}
