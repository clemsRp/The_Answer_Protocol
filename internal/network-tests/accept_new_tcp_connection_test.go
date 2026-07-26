package networktests

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"strings"
	"testing"
	"time"
)

func TestServerSendsProtocolGreeting(t *testing.T) {
	_, realAddress := setupTestServer(t)

	clientConn, err := net.DialTimeout("tcp", realAddress, 2*time.Second)
	if err != nil {
		t.Fatalf("Client failed to connect: %v", err)
	}
	defer clientConn.Close()

	clientConn.SetReadDeadline(time.Now().Add(1 * time.Second))

	reader := bufio.NewReader(clientConn)
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read server response: %v", err)
	}

	if !strings.Contains(response, "OK hello proto=1") {
		t.Errorf("Unexpected response received: %q", response)
	}
}

func TestServerLogsConnectionInternally(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	_, realAddress := setupTestServer(t)

	clientConn, err := net.DialTimeout("tcp", realAddress, 2*time.Second)
	if err != nil {
		t.Fatalf("Client failed to connect: %v", err)
	}
	clientConn.Close()

	time.Sleep(20 * time.Millisecond)

	logsOutput := buf.String()
	if !strings.Contains(logsOutput, "127.0.0.1") {
		t.Errorf("Server did not log the client IP. Current logs: %q", logsOutput)
	}
}
