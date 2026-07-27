package networktests

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"tap/internal/network"
	"testing"
)

func setupTestServer(t *testing.T) (*network.Server, string) {
	// t.Helper() tells to testing module that it is a helper and not a test!

	t.Helper()
	in := make(chan network.IncomingEvent)
	out := make(chan network.OutgoingEvent)

	server, err := network.NewServer("127.0.0.1:0", in, out)
	if err != nil {
		t.Fatalf("Server failed to start listening: %v", err)
	}

	server.Start()

	return server, server.GetAddress()
}

func sendCommand(t *testing.T, conn net.Conn, cmd string) string {
	t.Helper()

	_, err := fmt.Fprintf(conn, "%s\n", cmd)
	if err != nil {
		t.Fatalf("Erreur d'envoi: %v", err)
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Erreur de lecture: %v", err)
	}

	return strings.TrimSpace(response)
}
