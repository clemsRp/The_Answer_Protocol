package network

import (
	"bufio"
	"net"
	"strings"
	"tap/tests/utils"
	"testing"
	"time"
)

func TestEmptyAndWhitespaceCommands(t *testing.T) {
	s := utils.SetupTestServerEngine(t, "../../world.json")
	addr := s.GetAddress()

	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("Connection failed: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	_, _ = reader.ReadString('\n')

	payloads := []string{
		"\n",
		"    \n",
		"\r\n",
		"\t\t\n",
	}

	for _, payload := range payloads {
		_, err = conn.Write([]byte(payload))
		if err != nil {
			t.Fatalf("Failed to send empty payload: %v", err)
		}
	}

	_, err = conn.Write([]byte("CONNECT test_user\n"))
	if err != nil {
		t.Fatalf("Failed to send valid command: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	success := false
	for {
		res, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Server crashed or disconnected while reading responses: %v", err)
		}

		res = strings.TrimSpace(res)

		if strings.HasPrefix(res, "OK connected") {
			success = true
			break
		}
	}

	if !success {
		t.Errorf("Never received the 'OK connected' response.")
	} else {
		t.Log("Success: Server survived empty lines and processed the subsequent command.")
	}
}
