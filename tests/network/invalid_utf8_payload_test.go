package network

import (
	"bufio"
	"net"
	"tap/tests/utils"
	"testing"
	"time"
)

func TestInvalidUTF8Payload(t *testing.T) {
	s, _ := utils.SetupTestServerEngine(t, "../../world.json")
	addr := s.GetAddress()

	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("Connection failed: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	_, _ = reader.ReadString('\n')

	// (0xFF et 0xFE are forbidden in UTF8)
	invalidPayload := []byte{0xff, 0xfe, 0xfd, '\n'}

	_, err = conn.Write(invalidPayload)
	if err != nil {
		t.Fatalf("Failed to write invalid utf-8: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	res, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Server crashed instead of managing invalid non-utf8 caracters: %v", err)
	}

	t.Logf("Success : Server survived to non-UTF8 caracter injection, response : %s", res)
}
