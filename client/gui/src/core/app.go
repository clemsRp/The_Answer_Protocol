package core

import (
	"fmt"
	"strings"
	"sync"
	"tap/client/gui/src/parser"
	"tap/client/gui/src/ui"
	vars "tap/client/gui/src/variables"
	"tap/client/state"
	panel "tap/client/tui/panels"
	"tap/engine"
	"tap/protocol"
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
	closeOnce    sync.Once

	updateQueue chan func()
}

func NewApp(actionsChan chan panel.Action) *App {
	rl.SetTraceLogLevel(rl.LogNone)
	rl.SetTargetFPS(60)

	rl.InitWindow(0, 0, "TAP")

	monitor := rl.GetCurrentMonitor()
	screenWidth := rl.GetMonitorWidth(monitor)
	screenHeight := rl.GetMonitorHeight(monitor)
	// screenWidth := 900
	// screenHeight := 400

	rl.SetWindowSize(screenWidth, screenHeight)

	app := &App{
		Textures:     parser.LoadTextures(),
		Variables:    vars.GetVariables(),
		Colors:       vars.GetColors(),
		Manager:      ui.NewManager(),
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		ActionsChan:  actionsChan,
		updateQueue:  make(chan func(), 256),
	}

	var err error
	maps_folder_path := "./client/gui/maps/"
	app.Rooms, err = parser.ParseRooms([]string{maps_folder_path + engine.RoomEntrance})

	if err != nil {
		fmt.Println("Error parsing rooms:", err)
		app.Stop()
		return nil
	}

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
		Y: app.Variables.StartingPosY}

	return app
}

func (app *App) QueueUpdate(f func()) {
	if f == nil {
		return
	}
	select {
	case app.updateQueue <- f:
	default:
		fmt.Println("gui: update queue full, dropping a UI update")
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

func (app *App) ShowConnectPage() {
	app.Variables.Current_view = "Connect"

}

func (app *App) ShowGamePage() {
	app.Variables.Current_view = "Game"
	newPosX := app.Variables.Player.Position.X
	newPosY := app.Variables.Player.Position.Y
	newDirX := app.Variables.Player.Direction.X
	newDirY := app.Variables.Player.Direction.Y
	payload := fmt.Sprintf("%s %f %f %f %f", protocol.CmdNotifyPosition, newPosX, newPosY, newDirX, newDirY)

	app.ActionsChan <- panel.Action{
		Type:    panel.ActionSendServer,
		Payload: payload,
	}
	app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: protocol.CmdGetPositions}

}

func (app *App) ShowCombatPage()                                       {}
func (app *App) ShowPopupPage()                                        {}
func (app *App) ClosePopup()                                           {}
func (app *App) ShowCombatResultPopup(result string, rewards []string) {}
func (app *App) ShowQuestCompletedPopup(questID, reward string)        {}

func (app *App) UpdateNavigation(room *protocol.LookCommandData) {
	app.Variables.PanelsVariables.Room = room
	app.Variables.Current_room = strings.SplitN(room.Id, "room.", 2)[1]
}

func (app *App) UpdateItems(roomItems, inventory []string) {
	app.Variables.PanelsVariables.RoomItems = &roomItems
	app.Variables.PanelsVariables.Inventory = &inventory
}

func (app *App) UpdateDatas(text string) {}
func (app *App) UpdateInteraction(npcs, players []string, npcData map[string]protocol.InspectNPCData, npcDialogues map[string]string, groupMembers []string, quests []protocol.TrackedQuestData, completed_quests []string) {
}

func (app *App) UpdateGroup(groupState state.GroupState) {
	app.Variables.PanelsVariables.GroupState = &groupState
}

func (app *App) UpdateCombat(combatState state.CombatState) {
	app.Variables.PanelsVariables.CombatState = &combatState
}

func (app *App) UpdateQuests(quests []protocol.TrackedQuestData) {
	app.Variables.PanelsVariables.Quests = &quests
}

func (app *App) UpdateInspector(text string)                      {}
func (app *App) AppendChat(scope, user, msg string)               {}
func (app *App) AppendCombatChat(user, msg string)                {}
func (app *App) AppendServerResponse(res protocol.ServerResponse) {}
func (app *App) AppendCliMessage(text string)                     {}
func (app *App) AppendCliResponse(res protocol.ServerResponse)    {}

func (a *App) UpdateRemotePlayerPosition(pseudo string, x, y, dirX, dirY float32) {
	if pseudo == a.Variables.Player.Pseudo {
		return
	}

	if remotePlayer, exists := a.Variables.RemotePlayers[pseudo]; exists {
		remotePlayer.Position.X = x
		remotePlayer.Position.Y = y
		remotePlayer.Direction.X = dirX
		remotePlayer.Direction.Y = dirY
	} else {
		a.Variables.RemotePlayers[pseudo] = &vars.Player{
			Pseudo:    pseudo,
			Position:  &vars.Position{X: x, Y: y},
			Direction: &vars.Direction{X: dirX, Y: dirY},
		}
	}
}

func (a *App) AddRemotePlayer(pseudo string) {
	if a.Variables.RemotePlayers == nil {
		a.Variables.RemotePlayers = make(map[string]*vars.Player)
	}

	if _, exists := a.Variables.RemotePlayers[pseudo]; !exists {
		a.Variables.RemotePlayers[pseudo] = &vars.Player{
			Pseudo:    pseudo,
			Position:  &vars.Position{X: a.Variables.StartingPosX, Y: a.Variables.StartingPosY},
			Direction: &vars.Direction{X: 0, Y: 1},
		}
	}
	a.UpdateRemotePlayerPosition(pseudo, a.Variables.StartingPosX, a.Variables.StartingPosY, 0, 1)
}

func (a *App) RemoveRemotePlayer(pseudo string) {
	if a.Variables.RemotePlayers != nil {
		delete(a.Variables.RemotePlayers, pseudo)
	}
}
func (app *App) GetPseudo() string {
	return app.Variables.Player.Pseudo
}

func (app *App) SetPseudo(pseudo string) {
	app.Variables.Player.Pseudo = pseudo
}

func (app *App) Stop() {
	app.closeOnce.Do(func() {
		if app.Textures != nil {
			app.Textures.UnloadTextures()
		}
		rl.CloseWindow()
	})
}
