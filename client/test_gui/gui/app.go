package main

import (
	"test_gui/gui/parser"
	vars "test_gui/gui/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type App struct {
	textures  *parser.Textures
	variables *vars.Variables
	rooms     map[string]*parser.Room
	convertor *parser.ConvertTextures
}

func NewApp() *App {
	return &App{
		textures:  parser.LoadTextures(),
		variables: vars.GetVariables(),
	}
}

func (app *App) Update() {
	// Update player coordinates
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
	// Draw floor
	cur_tilesets := app.rooms[app.variables.Current_room].Tilesets
	for y := range len(cur_tilesets) {
		for x := range len(cur_tilesets[y]) {
			// Get tileset datas
			texture_name, indX, indY, err := app.convertor.GetTilesetDatas(cur_tilesets[y][x])

			// Calculate tileset coordinates
			coor_x := x * app.variables.Tileset_size
			coor_y := y * app.variables.Tileset_size

			if err == nil {
				parser.DrawImage(
					(*app.textures)[texture_name],
					float32(coor_x), float32(coor_y),
					float32(indX), float32(indY),
					1, 1, float32(app.variables.Tileset_size/vars.FRAME_WIDTH),
				)
			}
		}
	}
}
