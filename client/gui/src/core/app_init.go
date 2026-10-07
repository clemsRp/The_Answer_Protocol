package core

import (
	vars "tap/client/gui/src/variables"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (app *App) AddMissingVariables(screenWidth int) {
	app.initBaseVariables(screenWidth)
	app.initPlayerVariables()
	app.initPanelsVariables()
	app.initRemoteVariables()
}

func (app *App) initBaseVariables(screenWidth int) {
	app.Variables.Tileset_size = float32(screenWidth / 32)
	app.Variables.FontSize = 0.4 * app.Variables.Tileset_size
	app.Variables.Zoom = float32(app.Variables.Tileset_size) / float32(vars.FRAME_WIDTH)
	app.Variables.StartTime = time.Now()
}

func (app *App) initPlayerVariables() {
	app.Variables.Player.Speed = int(0.12 * app.Variables.Tileset_size)
	app.Variables.Player.LastTimeTyped = time.Now()
	app.Variables.StartingPosX = float32(15 * app.Variables.Tileset_size)
	app.Variables.StartingPosY = float32(9 * app.Variables.Tileset_size)
	app.Variables.Player.Position = &vars.Position{
		X: app.Variables.StartingPosX,
		Y: app.Variables.StartingPosY,
	}
	app.Variables.Player.Zoom = app.Variables.Zoom
}

func (app *App) initPanelsVariables() {
	tileSize := app.Variables.Tileset_size

	// Chat Panel
	app.Variables.PanelsVariables.Chat.Rect = rl.NewRectangle(
		vars.CHAT_START_X*tileSize,
		vars.CHAT_START_Y*tileSize,
		vars.CHAT_WIDTH*tileSize,
		vars.CHAT_HEIGHT*tileSize,
	)
	app.Variables.PanelsVariables.Chat.ScrollRect = rl.NewRectangle(
		(vars.CHAT_START_X+vars.CHAT_WIDTH-0.6)*tileSize,
		(vars.CHAT_START_Y+0.25)*tileSize,
		0.25*tileSize,
		(vars.CHAT_HEIGHT-0.5)*tileSize,
	)

	// Group Panel
	app.Variables.PanelsVariables.Group.Rect = rl.NewRectangle(
		vars.GROUP_START_X*tileSize,
		vars.GROUP_START_Y*tileSize,
		vars.GROUP_WIDTH*tileSize,
		vars.GROUP_HEIGHT*tileSize,
	)
	app.Variables.PanelsVariables.Group.ScrollRect = rl.NewRectangle(
		(vars.GROUP_START_X+vars.GROUP_WIDTH-0.6)*tileSize,
		(vars.GROUP_START_Y+0.25)*tileSize,
		0.25*tileSize,
		(vars.GROUP_HEIGHT-0.5)*tileSize,
	)

	// Talk Panel
	app.Variables.PanelsVariables.Talk.Rect = rl.NewRectangle(
		float32(vars.TALK_START_X)*tileSize,
		float32(vars.TALK_START_Y)*tileSize,
		float32(vars.TALK_WIDTH)*tileSize,
		float32(vars.TALK_HEIGHT)*tileSize,
	)
}

func (app *App) initRemoteVariables() {
	remote_players := make(map[string]*vars.Player)
	app.Variables.RemotePlayers = &remote_players

	remote_player_emotes := make(map[string]*vars.PlayerEmote)
	app.Variables.RemotePlayerEmotes = &remote_player_emotes

	app.Manager.PlayerPos = app.Variables.Player.Position
	app.Manager.Tile = app.Variables.Tileset_size
}

func (app *App) AddItemPositions() {
	positions := make(map[string]*vars.Position)

	// Set positions
	for it_name, it_datas := range vars.ItemConvertor {
		positions[it_name] = &vars.Position{X: it_datas.Pos.X, Y: it_datas.Pos.Y}
	}

	// Scale positions to map size
	for _, pos := range positions {
		(*pos).X *= app.Variables.Tileset_size
		(*pos).Y *= app.Variables.Tileset_size
	}

	app.Variables.ItemPositions = &positions
}
