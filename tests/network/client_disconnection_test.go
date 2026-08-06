package network

import (
	"net"
	"tap/tests/utils"
	"testing"
	"time"
)

func TestClientDisconnection(t *testing.T) {
	s, _ := utils.SetupTestServerEngine(t, "../../world.json")
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

func TestHardResetDisconnection(t *testing.T) {
	s, _ := utils.SetupTestServerEngine(t, "../../world.json")
	addr := s.GetAddress()

	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("Connection failed: %v", err)
	}

	// On cast la connexion en *net.TCPConn pour avoir accès aux options bas niveau
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		t.Fatal("Failed to cast to TCPConn")
	}

	// SetLinger(0) forces the OS to throw the remaining buffer
	// and sends  TCP RST immediatly when Close().
	tcpConn.SetLinger(0)
	tcpConn.Close()

	time.Sleep(100 * time.Millisecond)

	conn2, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("Server crashed after client abrupt disconnection (TCP RST): %v", err)
	}
	conn2.Close()
	t.Log("Success : Server survived to abrupt disconnection (TCP RST) (Connection reset by peer).")
}
