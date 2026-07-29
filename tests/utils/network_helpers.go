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
	sc "tap/tests/scenarios"
	"testing"
	"time"
)

func SetupTestServerEngine(t *testing.T) *server.Server {
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

func EstablishConnection(t *testing.T, s *server.Server, scenario sc.ScenariosCommandTest, connections map[string]net.Conn, readers map[string]*bufio.Reader) {
	t.Helper()
	user := scenario.TestOnConnection

	if _, exists := connections[user]; exists {
		return
	}

	conn, err := net.DialTimeout("tcp", s.GetAddress(), 2*time.Second)
	if err != nil {
		t.Fatalf("%s failed to connect: %v", user, err)
	}

	connections[user] = conn
	readers[user] = bufio.NewReader(conn)

	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	res, err := readers[user].ReadString('\n')
	if err != nil {
		t.Fatalf("Read error during greeting for %s: %v", user, err)
	}

	if !strings.HasPrefix(res, "OK hello proto=") {
		t.Error(FormatMismatch(scenario.Command, "OK hello proto=", res))
	}
}
