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
	exchanger := pr.Exchanger{ServerInput: make(chan pr.ServerRequest, 100),
		ServerOutput: make(chan pr.EngineResponse, 100),
		JoinChan:     make(chan string, 10),
		LeaveChan:    make(chan string, 10)}

	var err error
	var world parser.Map
	world, err = parser.Get_map(world_path)
	if err != nil {
		t.Fatalf("ERROR parsing: %v", err.Error())
	}

	// Initialize server
	var s *server.Server
	s, err = server.NewServer("localhost:0", exchanger)
	if err != nil {
		t.Fatalf("Server couldn't start %v", err)
	}

	// Initialize and start engine
	e := engine.NewEngine(world, exchanger)
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
