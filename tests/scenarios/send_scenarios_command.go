package scenarios

import (
	"fmt"
	"net"
	"testing"
)

func SendScenarioCommand(t *testing.T, conn net.Conn, cmd string) {
	t.Helper()

	_, err := fmt.Fprintf(conn, "%s\n", cmd)
	if err != nil {
		t.Fatalf("Send error: %v", err)
	}
}
