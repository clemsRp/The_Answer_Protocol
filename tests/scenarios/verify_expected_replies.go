package scenarios

import (
	"bufio"
	"net"
	"tap/tests/utils"
	"testing"
	"time"
)

func VerifyExpectedReplies(t *testing.T, scenario ScenariosCommandTest, connections map[string]net.Conn, readers map[string]*bufio.Reader) {
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

		utils.AssertResponse(t, scenario.Command, reply.Msg, res, scenario.ExpectsJSON)
	}
}
