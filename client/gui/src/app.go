package gui

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"tap/client/gui/src/parser"
	vars "tap/client/gui/src/variables"
	"tap/client/state"
	panel "tap/client/tui/panels"
	"tap/engine"
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
	maps_folder_path := "./client/gui/tiled_maps/"
	app.rooms, err = parser.ParseRooms([]string{maps_folder_path + engine.RoomEntrance})

	if err != nil {
		fmt.Println("Error parsing rooms:", err)
		app.Stop()
		return nil
	}

	app.variables.Tileset_size = screenWidth / 30

	return app
}

func (app *App) Update() {
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

	app.DrawMap(engine.RoomEntrance)
}
func (app *App) DrawMap(map_name string) {
	cur_room := app.rooms[map_name]
	if cur_room == nil {
		return
	}

	cur_tilesets := cur_room.Tilesets

	const (
		FLIPPED_HORIZONTALLY_FLAG = 0x80000000
		FLIPPED_VERTICALLY_FLAG   = 0x40000000
		FLIPPED_DIAGONALLY_FLAG   = 0x20000000
	)

	for _, layer := range cur_room.Layers {
		if !layer.Visible {
			continue
		}

		for i, rawTile := range layer.Data {
			if rawTile == 0 {
				continue
			}

			flipH := (rawTile & FLIPPED_HORIZONTALLY_FLAG) != 0
			flipV := (rawTile & FLIPPED_VERTICALLY_FLAG) != 0
			flipD := (rawTile & FLIPPED_DIAGONALLY_FLAG) != 0

			tile := rawTile & 0x0FFFFFFF
			if tile == 0 {
				continue
			}

			gridX := i % layer.Width
			gridY := i / layer.Width

			posX := float32(gridX * app.variables.Tileset_size)
			posY := float32(gridY * app.variables.Tileset_size)

			var activeTileset parser.Tileset
			for j := len(cur_tilesets) - 1; j >= 0; j-- {
				if tile >= cur_tilesets[j].FirstGID {
					activeTileset = cur_tilesets[j]
					break
				}
			}

			if activeTileset.FirstGID == 0 {
				continue
			}

			sourcePath := activeTileset.Source
			baseName := filepath.Base(sourcePath)
			textureName := strings.TrimSuffix(baseName, filepath.Ext(baseName))

			texture, ok := (*app.textures)[textureName]
			if !ok {
				continue
			}

			localID := tile - activeTileset.FirstGID

			cols := int(texture.Width) / vars.FRAME_WIDTH
			if cols == 0 {
				cols = 1
			}

			indX := float32(localID % cols)
			indY := float32(localID / cols)

			var rotation float32 = 0
			var rX, rY float32 = 1, 1

			if flipD {
				rotation = 90
				rX = 1
				rY = -1
				if flipH {
					rX = -rX
				}
				if flipV {
					rY = -rY
				}
			} else {
				if flipH {
					rX = -1
				}
				if flipV {
					rY = -1
				}
			}

			zoom := float32(app.variables.Tileset_size) / float32(vars.FRAME_WIDTH)

			parser.DrawImage(texture, posX, posY, indX, indY, rX, rY, zoom, rotation)
		}
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

func (app *App) QueueUpdate(f func())                                  {}
func (app *App) ShowConnectPage()                                      {}
func (app *App) ShowGamePage()                                         {}
func (app *App) ShowCombatPage()                                       {}
func (app *App) ShowPopupPage()                                        {}
func (app *App) ClosePopup()                                           {}
func (app *App) ShowCombatResultPopup(result string, rewards []string) {}
func (app *App) ShowQuestCompletedPopup(questID, reward string)        {}
func (app *App) UpdateNavigation(room *protocol.LookCommandData)       {}
func (app *App) UpdateItems(roomItems, inventory []string)             {}
func (app *App) UpdateInteraction(npcs, players []string, npcData map[string]protocol.InspectNPCData, npcDialogues map[string]string, groupMembers []string, quests []protocol.TrackedQuestData, completed_quests []string) {
}
func (app *App) UpdateGroup(groupState state.GroupState)          {}
func (app *App) UpdateCombat(combatState state.CombatState)       {}
func (app *App) UpdateDatas(text string)                          {}
func (app *App) UpdateQuests(quests []protocol.TrackedQuestData)  {}
func (app *App) UpdateInspector(text string)                      {}
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
