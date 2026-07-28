package networktests

import (
	"net"
	"testing"
	"time"
)

func TestClientNormalDeconnection(t *testing.T) {
	s := setupTestServerEngine(t)

	clientConn1, err := net.DialTimeout("tcp", s.GetAddress(), 2*time.Second)
	if err != nil {
		t.Fatalf("Client 1 failed to connect: %v", err)
	}

	clientConn1.Close()

	time.Sleep(50 * time.Millisecond)

	clientConn2, err := net.DialTimeout("tcp", s.GetAddress(), 2*time.Second)
	if err != nil {
		t.Fatalf("Server crashed or stopped accepting connections after EOF: %v", err)
	}
	defer clientConn2.Close()

	if clientConn2 == nil {
		t.Fatal("Client 2 connection attempt returned nil")
	}
}
