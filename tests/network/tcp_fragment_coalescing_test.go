package network

import (
	"strings"
	"tap/tests/utils"
	"testing"
	"time"
)

func TestTCPFragmentation(t *testing.T) {
	s, _ := utils.SetupTestServerEngine(t, "../../world.json")
	conn, reader := utils.ConnectAndGreet(t, s.GetAddress())
	defer conn.Close()

	if _, err := conn.Write([]byte("CONN")); err != nil {
		t.Fatalf("Failed to write first chunk: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	if _, err := conn.Write([]byte("ECT alice\n")); err != nil {
		t.Fatalf("Failed to write second chunk: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	res, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Server failed to respond to fragmented command: %v", err)
	}

	if strings.TrimSpace(res) != "OK connected" {
		t.Errorf("Expected 'OK connected', got: '%s'", res)
	} else {
		t.Log("Success: Server correctly buffered and processed the fragmented command.")
	}
}

func TestTCPCoalescing(t *testing.T) {
	s, _ := utils.SetupTestServerEngine(t, "../../world.json")
	conn, reader := utils.ConnectAndGreet(t, s.GetAddress())
	defer conn.Close()

	payload := []byte("CONNECT bob\nLOOK\n")
	if _, err := conn.Write(payload); err != nil {
		t.Fatalf("Failed to write coalesced chunk: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(1 * time.Second))

	res1, err := reader.ReadString('\n')
	if err != nil || !strings.HasPrefix(res1, "OK connected") {
		t.Fatalf("Failed on first response. Got: '%s', Err: %v", res1, err)
	}

	res2, err := reader.ReadString('\n')
	if err != nil || (!strings.HasPrefix(res2, "OK {") && !strings.HasPrefix(res2, "ERR")) {
		t.Fatalf("Failed on second response. Got: '%s', Err: %v", res2, err)
	}

	t.Log("Success: Server correctly separated and processed coalesced commands.")
}
