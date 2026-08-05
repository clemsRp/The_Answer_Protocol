package network

import (
	"bufio"
	"bytes"
	"log"
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

func TestServerLogsNewConnection(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	s, _ := utils.SetupTestServerEngine(t, "../../world.json")

	clientConn, err := net.DialTimeout("tcp", s.GetAddress(), 2*time.Second)
	if err != nil {
		t.Fatalf("Client failed to connect: %v", err)
	}
	clientConn.Close()

	time.Sleep(20 * time.Millisecond)

	logsOutput := buf.String()
	if !strings.Contains(logsOutput, "127.0.0.1") {
		t.Error(utils.FormatMismatch("", "CLIENT IP", logsOutput))
	}
}
