package network

import (
	"bufio"
	"net"
	"strings"
	"tap/tests/utils"
	"testing"
	"time"
)

func TestServerSendsGreeting(t *testing.T) {
	s, _ := utils.SetupTestServerEngine(t, "../../world.json")

	clientConn, err := net.DialTimeout("tcp", s.GetAddress(), 2*time.Second)
	if err != nil {
		t.Fatalf("Client failed to connect: %v", err)
	}
	defer clientConn.Close()

	reader := bufio.NewReader(clientConn)
	res, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read server response: %v", err)
	}

	res = strings.TrimSpace(res)

	if res != "OK hello proto=1" {
		t.Error(utils.FormatMismatch("", "OK hello proto=1", res))
	}
}
