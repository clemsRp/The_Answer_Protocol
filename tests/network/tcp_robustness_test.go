package network

import (
	"bufio"
	"net"
	"strings"
	"tap/tests/utils"
	"testing"
	"time"
)

func TestTCPFragmentation(t *testing.T) {
	s := utils.SetupTestServerEngine(t, "../../world.json")
	addr := s.GetAddress()

	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("Connection failed: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	_, _ = reader.ReadString('\n')

	_, err = conn.Write([]byte("CONN"))
	if err != nil {
		t.Fatalf("Failed to write first chunk: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	_, err = conn.Write([]byte("ECT alice\n"))
	if err != nil {
		t.Fatalf("Failed to write second chunk: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	res, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Server failed to respond to fragmented command: %v", err)
	}

	res = strings.TrimSpace(res)
	if res != "OK connected" {
		t.Errorf("Expected 'OK connected', got: '%s'", res)
	} else {
		t.Log("Success: Server correctly buffered and processed the fragmented command.")
	}
}

func TestTCPCoalescing(t *testing.T) {
	s := utils.SetupTestServerEngine(t, "../../world.json")
	addr := s.GetAddress()

	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("Connection failed: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	_, _ = reader.ReadString('\n')

	payload := []byte("CONNECT bob\nLOOK\n")
	_, err = conn.Write(payload)
	if err != nil {
		t.Fatalf("Failed to write coalesced chunk: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(1 * time.Second))

	res1, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read first response: %v", err)
	}
	if !strings.HasPrefix(res1, "OK connected") {
		t.Errorf("First response expected to be 'OK connected', got: '%s'", res1)
	}

	res2, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read second response (LOOK was ignored): %v", err)
	}
	if !strings.HasPrefix(res2, "OK {") && !strings.HasPrefix(res2, "ERR") {
		t.Errorf("Second response invalid for LOOK command, got: '%s'", res2)
	} else {
		t.Log("Success: Server correctly separated and processed coalesced commands.")
	}
}
