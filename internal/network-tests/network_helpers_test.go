package networktests

import (
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
