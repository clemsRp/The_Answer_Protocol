package scenarios

import (
	"bufio"
	"net"
	"strings"
	"tap/server"
	"tap/tests/utils"
	"testing"
	"time"
)

func EstablishConnectionForScenario(t *testing.T, s *server.Server, scenario ScenariosCommandTest, connections map[string]net.Conn, readers map[string]*bufio.Reader) {
	t.Helper()
	user := scenario.TestOnConnection

	if _, exists := connections[user]; exists {
		return
	}

	conn, err := net.DialTimeout("tcp", s.GetAddress(), 2*time.Second)
	if err != nil {
		t.Fatalf("%s failed to connect: %v", user, err)
	}

	connections[user] = conn
	readers[user] = bufio.NewReader(conn)

	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	res, err := readers[user].ReadString('\n')
	if err != nil {
		t.Fatalf("Read error during greeting for %s: %v", user, err)
	}

	if !strings.HasPrefix(res, "OK hello proto=") {
		t.Error(utils.FormatMismatch(scenario.Command, "OK hello proto=", res))
	}
}
