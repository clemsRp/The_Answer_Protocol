package utils

import (
	"bufio"
	"encoding/json"
	"net"
	"strings"
	sc "tap/tests/scenarios"
	"testing"
	"time"
)

func AssertResponse(t *testing.T, command, wanted, got string, expectsJSON bool) {
	t.Helper()

	if expectsJSON {
		jsonPart := strings.TrimPrefix(got, wanted+" ")
		jsonPart = strings.TrimSpace(jsonPart)

		if !strings.HasPrefix(got, wanted) {
			t.Error(FormatMismatch(command, wanted, got))
		} else if !json.Valid([]byte(jsonPart)) {
			t.Error(FormatInvalidJSON(jsonPart))
		}
	} else {
		if !strings.HasPrefix(got, wanted) {
			t.Error(FormatMismatch(command, wanted, got))
		}
	}
}

func VerifyExpectedReplies(t *testing.T, scenario sc.ScenariosCommandTest, connections map[string]net.Conn, readers map[string]*bufio.Reader) {
	t.Helper()

	for _, reply := range scenario.ExpectedReplies {
		conn, exists := connections[reply.User]
		if !exists {
			t.Fatalf("Test setup error: No connection found for user '%s' to receive message '%s'", reply.User, reply.Msg)
		}

		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		res, err := readers[reply.User].ReadString('\n')
		if err != nil {
			t.Fatalf("Read error (timeout) for %s: %v", reply.User, err)
		}

		AssertResponse(t, scenario.Command, reply.Msg, res, scenario.ExpectsJSON)
	}
}
