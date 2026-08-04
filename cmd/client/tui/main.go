package main

import (
	"fmt"
	"tap/client/tui"
	"tap/engine/parser"
)

var (
	world parser.Map
)

func main() {
	// Get the world
	var err error
	world, err = parser.Get_map("world.json")
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	// Initialize client
	cli := tui.NewTuiClient(world)

	cli.Start()
}
