package network

import (
	"net"
	"tap/tests/utils"
	"testing"
	"time"
)

func TestServerGracefulShutdown_RefuseNewConnections(t *testing.T) {
	s := utils.SetupTestServerEngine(t, "../../world.json")
	addr := s.GetAddress()

	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		t.Fatalf("Premature failure: the server should accept connections when active: %v", err)
	}
	conn.Close()

	s.Stop()

	lateConn, err := net.DialTimeout("tcp", addr, 1*time.Second)

	if err == nil {
		lateConn.Close()
		t.Fatal("Fail: The server accepted a new connection when it was supposed to be stopped.")
	}

	t.Logf("Success: The server successfully refused the connection (Received error: %v)", err)
}
