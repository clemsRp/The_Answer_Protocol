package utils

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func ConnectAndGreet(t *testing.T, addr string) (net.Conn, *bufio.Reader) {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("Helper failed to connect: %v", err)
	}

	reader := bufio.NewReader(conn)

	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	_, err = reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Helper failed to read greeting: %v", err)
	}

	conn.SetReadDeadline(time.Time{})

	return conn, reader
}
