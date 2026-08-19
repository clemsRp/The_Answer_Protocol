package main

import (
	"fmt"
	"test_gui/gui/parser"
	vars "test_gui/gui/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.SetTraceLogLevel(rl.LogNone)
	rl.SetTargetFPS(60)

	rl.InitWindow(0, 0, "")
	monitor := rl.GetCurrentMonitor()
	screenWidth := rl.GetMonitorWidth(monitor)
	screenHeight := rl.GetMonitorHeight(monitor)
	rl.CloseWindow()

	rl.InitWindow(int32(screenWidth), int32(screenHeight), "TAP: OAP")
	defer rl.CloseWindow()

	// Init app
	var err error

	app := NewApp()
	defer app.textures.UnloadTextures()

	// Get rooms
	app.rooms, err = parser.ParseRooms([]string{"true"})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get convertor
	app.convertor, err = parser.ParseConvertTextures(vars.CONVERT_TEXTURES_PATH, app.textures)
	if err != nil {
		fmt.Println(err)
		return
	}

	app.variables.Tileset_size = screenWidth / 30

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
