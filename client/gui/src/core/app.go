package core

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"tap/client/gui/src/parser"
	"tap/client/gui/src/ui"
	vars "tap/client/gui/src/variables"
	panel "tap/client/tui/panels"
	"tap/engine"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type App struct {
	Textures     *parser.Textures
	Variables    *vars.Variables
	Colors       map[string]rl.Color
	Manager      *ui.Manager
	Rooms        map[string]*parser.Map
	WorldItems   map[string]*engine.Item
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

	// monitor := rl.GetCurrentMonitor()
	// screenWidth := rl.GetMonitorWidth(monitor)
	// screenHeight := rl.GetMonitorHeight(monitor)
	screenWidth := 800
	screenHeight := 500

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

	// Get all maps
	maps_paths := make([]string, 0)
	for _, room := range engine.ValidMaps {
		maps_paths = append(maps_paths, maps_folder_path+room)
	}

	app.Rooms, err = parser.ParseRooms(maps_paths)

	// Handle error
	if err != nil {
		fmt.Println("Error parsing rooms:", err)
		app.Stop()
		return nil
	}

	app.AddMissingVariables(screenWidth)
	app.AddItemPositions()
	app.AddNpcPositions()

	return app
}

func (app *App) ParseNpcNames(filepath string) {
	// Parse json file
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	var world []engine.Map
	if err := json.Unmarshal(data, &world); err != nil {
		fmt.Printf("JSON parsing error: %v\n", err)
		return
	}

	if len(world) > 0 {
		app.WorldItems = world[0].Items

		// Init NpcConvertor
		app.Variables.NpcConvertor = make(map[string]string)
		for npc_id, npc := range world[0].Npcs {
			app.Variables.NpcConvertor[npc.Name] = npc_id
			app.Variables.NpcConvertor[npc_id] = npc.Name
		}

		// Init ItemConvertor
		app.Variables.ItemConvertor = make(map[string]string)
		for item_id, item := range world[0].Items {
			app.Variables.ItemConvertor[item.Name] = item_id
			app.Variables.ItemConvertor[item_id] = item.Name
		}
	}
}

func (app *App) IsItemUsable(itemName string) bool {
	if app.WorldItems == nil {
		app.ParseNpcNames("./world.json")
	}

	cleanName := strings.TrimSpace(itemName)
	lowerName := strings.ToLower(cleanName)
	snakeName := strings.ReplaceAll(lowerName, " ", "_")

	if app.WorldItems != nil {
		for k, item := range app.WorldItems {
			if item == nil {
				continue
			}
			itemK := strings.ToLower(k)
			itemN := strings.ToLower(item.Name)
			if itemK == lowerName || itemK == snakeName || itemN == lowerName {
				return item.Type == "weapon" || item.Type == "consumable"
			}
		}
	}

	return true
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
