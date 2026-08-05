package network

import (
	"bufio"
	"net"
	"tap/tests/utils"
	"testing"
	"time"
)

func TestMaxPlayersLimit(t *testing.T) {
	s, _ := utils.SetupTestServerEngine(t, "../../world.json")
	addr := s.GetAddress()

	const maxPlayers = 50
	var connections []net.Conn

	defer func() {
		for _, c := range connections {
			c.Close()
		}
	}()

	for i := 0; i < maxPlayers; i++ {
		conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
		if err != nil {
			t.Fatalf("Failed to connect client %d: %v", i, err)
		}
		connections = append(connections, conn)

		bufio.NewReader(conn).ReadString('\n')
	}

	intruderConn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		t.Log("SUCCESS : Server refused the connection.")
		return
	}
	defer intruderConn.Close()

	intruderConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	reader := bufio.NewReader(intruderConn)
	res, err := reader.ReadString('\n')

	if err == nil {
		if res == "ERR 901 SERVER_FULL\n" {
			t.Logf("SUCCESS : Server politely refused the connection with the correct error message: %s.", res)
		} else {
			t.Fatalf("FAIL: Server accepted the %dth player and answered with an unexpected message : %s", maxPlayers+1, res)
		}
	} else {
		t.Logf("SUCCESS : Server closed the connection. Error read: %v", err)
	}
}
