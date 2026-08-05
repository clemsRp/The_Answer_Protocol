package network

import (
	"bytes"
	"net"
	"tap/tests/utils"
	"testing"
	"time"
)

var server_limit = 1024

// 1024 bytes per line maximum
func TestProtectionAgainstHugePayloads(t *testing.T) {
	s := utils.SetupTestServerEngine(t, "../../world.json")
	addr := s.GetAddress()

	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("Connection failed: %v", err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	conn.Read(buf)

	// Create a payload larger than the default bufio.Scanner buffer (1KB).
	// We send 1KB of garbage without a newline '\n'.
	hugePayload := bytes.Repeat([]byte("A"), server_limit)

	_, err = conn.Write(hugePayload)
	if err != nil {
		t.Fatalf("Failed to send huge payload: %v", err)
	}

	// The server should forcefully close the connection to protect its memory.
	// Therefore, a subsequent read should fail (EOF or connection reset).
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	_, err = conn.Read(buf)

	if err == nil {
		t.Error("Fail: The server kept the connection open after receiving a huge payload. It should have dropped it to prevent memory exhaustion.")
	} else {
		t.Logf("Success: Server dropped the connection upon receiving a huge payload (Error: %v)", err)
	}
}
