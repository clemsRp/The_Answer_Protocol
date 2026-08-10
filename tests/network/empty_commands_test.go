package network

import (
	"strings"
	"tap/tests/utils"
	"testing"
	"time"
)

func TestEmptyAndWhitespaceCommands(t *testing.T) {
	s, _ := utils.SetupTestServerEngine(t, "../../world.json")
	conn, reader := utils.ConnectAndGreet(t, s.GetAddress())
	defer conn.Close()

	payloads := []string{"\n", "    \n", "\r\n", "\t\t\n"}

	for _, payload := range payloads {
		if _, err := conn.Write([]byte(payload)); err != nil {
			t.Fatalf("Failed to send empty payload: %v", err)
		}
	}

	if _, err := conn.Write([]byte("CONNECT test_user\n")); err != nil {
		t.Fatalf("Failed to send valid command: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	success := false

	for {
		res, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Server crashed or disconnected while reading responses: %v", err)
		}

		if strings.HasPrefix(strings.TrimSpace(res), "OK connected") {
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
