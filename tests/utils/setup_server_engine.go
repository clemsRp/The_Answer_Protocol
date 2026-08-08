package utils

import (
	"io"
	"log"
	"tap/engine"
	"tap/engine/parser"
	pr "tap/protocol"
	"tap/server"
	"testing"
	"time"
)

func SetupTestServerEngine(t *testing.T, world_path string) (*server.Server, *engine.Engine) {
	t.Helper()
	serverInput := make(chan pr.ServerRequest, 100)
	serverOutput := make(chan pr.EngineResponse, 100)
	updateClients := make(chan map[string]*pr.Client, 10)

	var err error
	var world parser.Map
	world, err = parser.Get_map(world_path)
	if err != nil {
		t.Fatalf("ERROR parsing: %v", err.Error())
	}

	// Initialize server
	var s *server.Server
	s, err = server.NewServer("localhost:0", serverInput, serverOutput, updateClients)
	if err != nil {
		t.Fatalf("Server couldn't start %v", err)
	}

	// Initialize and start engine
	e := engine.NewEngine(world, serverInput, serverOutput, updateClients)
	go e.Start()
	time.Sleep(10 * time.Millisecond)
	// Start the serveur
	go s.Start()
	time.Sleep(10 * time.Millisecond)

	t.Cleanup(func() {
		s.Stop()
		e.Stop()
	})
	log.SetOutput(io.Discard)

	return s, e
}
