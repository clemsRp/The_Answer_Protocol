package network

import (
	"bytes"
	"strings"
	"tap/src/server"
	"tap/src/tests/utils"
	"testing"
	"time"
)

func TestProtectionAgainstHugePayloads(t *testing.T) {
	// 1024 bytes per line maximum
	s, _ := utils.SetupTestServerEngine(t)
	conn, _ := utils.ConnectAndGreet(t, s.GetAddress())
	defer conn.Close()

	overflowSize := server.MaxPayloadSize
	hugePayload := bytes.Repeat([]byte("A"), overflowSize)

	_, err := conn.Write(hugePayload)
	if err != nil {
		t.Fatalf("Failed to send huge payload: %v", err)
	}

	buf := make([]byte, 100)
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("Expected error response before closing connection: %v", err)
	}

	response := string(buf[:n])
	if !strings.HasPrefix(response, "ERR") {
		t.Errorf("Expected ERR response, got: %s", response)
	}

	// Verify the connection is closed by server
	_, err = conn.Read(buf)
	if err == nil {
		t.Error("Fail: The server kept the connection open after receiving a huge payload. It should have closed it.")
	} else {
		t.Logf("Success: Server sent ERR and closed connection (Error: %v)", err)
	}
}
