package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tap/engine"
	pr "tap/protocol"
	"tap/server"
)

func main() {
	var (
		world     *engine.Map
		err       error
		exchanger pr.Exchanger
	)
	exchanger = pr.Exchanger{ServerInput: make(chan pr.ServerRequest, 100),
		ServerOutput: make(chan pr.EngineResponse, 100),
		JoinChan:     make(chan string, 10),
		LeaveChan:    make(chan string, 10)}

	world, err = engine.Get_map("world.json")
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	var s *server.Server
	s, err = server.NewServer("localhost:8080", exchanger)
	if err != nil {
		fmt.Println("Server couldn't start:", err)
		return
	}

	e := engine.NewEngine(world, exchanger)
	go e.Start()

	go s.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("\n Stopping the server due to Signal...")

	s.Stop()
	e.Stop()

	fmt.Println("Server and Engine stopped with success.")
}
