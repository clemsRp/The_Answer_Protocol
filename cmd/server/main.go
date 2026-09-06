package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"tap/engine"
	pr "tap/protocol"
	"tap/server"
)

func main() {
	// Define global variables
	var (
		world     *engine.Map
		err       error
		exchanger pr.Exchanger
	)

	// Set logger
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
	}

	// Create logger JSON handler
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	// Define logger as global and unique logger
	slog.SetDefault(logger)

	// Set exchanger
	exchanger = pr.Exchanger{ServerInput: make(chan pr.ServerRequest, 100),
		ServerOutput: make(chan pr.EngineResponse, 100),
		JoinChan:     make(chan string, 10),
		LeaveChan:    make(chan string, 10)}

	// Get the world map
	world, err = engine.Get_map("world.json")
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	// Initialize server
	var s *server.Server
	s, err = server.NewServer("localhost:8080", exchanger)
	if err != nil {
		fmt.Println("Server couldn't start:", err)
		return
	}

	// Start server/engine
	e := engine.NewEngine(world, exchanger)

	errChan := make(chan error, 1)

	// Enveloppe tes lancements pour capturer les erreurs fatales
	go func() {
		// Si la méthode Start() est bloquante et retourne une erreur quand elle plante
		if err := e.Start(); err != nil {
			errChan <- fmt.Errorf("Engine crash: %w", err)
		}
	}()

	go func() {
		if err := s.Start(); err != nil {
			errChan <- fmt.Errorf("Server crash: %w", err)
		}
	}()

	// Handle disconnection
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		fmt.Println("Server and Engine stopped with success.")

	case fatalErr := <-errChan:
		fmt.Printf("\nFatal crash : %v\nStopping system...\n", fatalErr)
	}

	// Stop server/engine
	s.Stop()
	e.Stop()

}
