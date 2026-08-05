package network

import (
	"bufio"
	"net"
	"strings"
	"tap/tests/utils"
	"testing"
	"time"
)

func TestUnknownCommandBehavior(t *testing.T) {
	s := utils.SetupTestServerEngine(t, "../../world.json")
	addr := s.GetAddress()

	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("Connection failed: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	_, _ = reader.ReadString('\n')

	_, err = conn.Write([]byte("DO_A_BARREL_ROLL\n"))
	if err != nil {
		t.Fatalf("Failed to send unknown command: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	res, err := reader.ReadString('\n')

	if err != nil {
		t.Fatalf("Server dropped connection instead of returning an error message: %v", err)
	}

	if !strings.HasPrefix(res, "ERR") {
		t.Errorf("Expected response to start with 'ERR', got: '%s'", res)
	} else {
		t.Logf("Success: Server handled unknown command gracefully with response: %s", strings.TrimSpace(res))
	}
}
