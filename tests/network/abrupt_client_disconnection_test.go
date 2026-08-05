package network

import (
	"net"
	"tap/tests/utils"
	"testing"
	"time"
)

func TestAbruptClientDisconnection(t *testing.T) {
	s := utils.SetupTestServerEngine(t, "../../world.json")
	addr := s.GetAddress()

	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("Connection failed: %v", err)
	}

	conn.Close()

	time.Sleep(100 * time.Millisecond)

	conn2, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("Fail: Server appears to have crashed after a client abruptly disconnected: %v", err)
	}
	defer conn2.Close()

	t.Log("Success: Server survived an abrupt client disconnection.")
}
