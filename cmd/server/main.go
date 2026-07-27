package main

import (
	"fmt"
	"tap/engine"
	"tap/engine/parser"
	pr "tap/protocol"
	"tap/server"
	"time"
)

var (
	t_start        = time.Now().Unix()
	update_clients = make(chan map[string]*pr.Client, 10)
	world          parser.Map
)

func main() {
	// Get the world
	serverInput := make(chan pr.ServerRequest, 100)
	serverOutput := make(chan pr.EngineResponse, 100)
	var err error
	world, err = parser.Get_map("world.json")
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	// Initialize server
	var s *server.Server
	s, err = server.NewServer("8080", serverInput, serverOutput, update_clients)
	if err != nil {
		fmt.Println("Server couldn't start")
		return
	}

	// Initialize and start engine
	e := engine.NewEngine(world, serverInput, serverOutput, update_clients)
	e.Start()

	// Start the serveur
	s.Start()
}
