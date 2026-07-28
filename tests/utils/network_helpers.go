package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"tap/engine"
	"tap/engine/parser"
	pr "tap/protocol"
	"tap/server"
	"testing"
)

func SetupTestServerEngine(t *testing.T) *server.Server {
	// t.Helper() tells to testing module that it is a helper and not a test!

	t.Helper()
	serverInput := make(chan pr.ServerRequest, 100)
	serverOutput := make(chan pr.EngineResponse, 100)
	updateClients := make(chan map[string]*pr.Client, 10)

	var err error
	var world parser.Map
	world, err = parser.Get_map("../world.json")
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

	// Start the serveur
	go s.Start()
	t.Cleanup(func() {
		s.Stop()
		e.Stop()
	})
	log.SetOutput(io.Discard)

	return s
}

func SendCommand(t *testing.T, conn net.Conn, cmd string) string {
	t.Helper()

	_, err := fmt.Fprintf(conn, "%s\n", cmd)
	if err != nil {
		t.Fatalf("Send error: %v", err)
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	return strings.TrimSpace(response)
}

func SendScenarioCommand(t *testing.T, conn net.Conn, cmd string) {
	t.Helper()

	_, err := fmt.Fprintf(conn, "%s\n", cmd)
	if err != nil {
		t.Fatalf("Send error: %v", err)
	}
}
