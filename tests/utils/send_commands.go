package utils

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
)

func SendCommand(t *testing.T, conn net.Conn, cmd string) string {
	t.Helper()

	_, err := fmt.Fprintf(conn, "%s\n", cmd)
	if err != nil {
		t.Fatalf("Send error: %v", err)
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	return strings.TrimSpace(response)
}


