package network

import (
	"bytes"
	"net"
	"tap/server"
	"tap/tests/utils"
	"testing"
	"time"
)

// 1024 bytes per line maximum
func TestProtectionAgainstHugePayloads(t *testing.T) {
	s, _ := utils.SetupTestServerEngine(t, "../../world.json")
	conn, _ := utils.ConnectAndGreet(t, s.GetAddress())
	defer conn.Close()

	overflowSize := server.MaxPayloadSize
	hugePayload := bytes.Repeat([]byte("A"), overflowSize)

	_, err := conn.Write(hugePayload)
	if err != nil {
		t.Fatalf("Failed to send huge payload: %v", err)
	}

	buf := make([]byte, 10)
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, err = conn.Read(buf)

	if err == nil {
		t.Error("Fail: The server kept the connection open after receiving a huge payload. It should have dropped it to prevent memory exhaustion.")
	} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		t.Error("Fail: The server kept the connection open (Test timed out waiting for the server to close the socket).")
	} else {
		t.Logf("Success: Server dropped the connection upon receiving a huge payload (Error: %v)", err)
	}
}
