package main

import (
	"tap/client/gui/src/parser"
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type App struct {
	textures  *parser.Textures
	variables *vars.Variables
	rooms     map[string]*parser.Map
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
