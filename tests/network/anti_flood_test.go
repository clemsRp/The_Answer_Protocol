package network

import (
	"bytes"
	"strings"
	pr "tap/protocol"
	"tap/tests/utils"
	"testing"
	"time"
)

func TestAntiFloodProtection(t *testing.T) {
	s, _ := utils.SetupTestServerEngine(t, "../../world.json")
	conn, reader := utils.ConnectAndGreet(t, s.GetAddress())
	defer conn.Close()

	burstPayload := bytes.Repeat([]byte("ACTION\n"), 10)
	if _, err := conn.Write(burstPayload); err != nil {
		t.Fatalf("Failed to send burst commands: %v", err)
	}

	spamErrorFound := false
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if strings.Contains(line, pr.ErrSpam) {
			spamErrorFound = true
			break
		}
	}

	// 3. Bilan du test
	if !spamErrorFound {
		t.Error("Fail: The server did not detect the spam, did not send error, or kept the connection open.")
	} else {
		t.Log("Success: Spam detected, error sent to the client, and connection successfully terminated.")
	}
}
