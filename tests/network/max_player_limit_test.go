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
		// Le système a carrément refusé la connexion (TCP RST)
		t.Log("SUCCESS : Server refused the connection.")
		return
	}
	defer intruderConn.Close()

	intruderConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	reader := bufio.NewReader(intruderConn)
	res, err := reader.ReadString('\n')

	if err == nil {
		t.Fatalf("FAIL: Server accepted the %dth player and answered : %s", maxPlayers+1, res)
	} else {
		t.Log("SUCCESS : Server closed the connection of the exceeded connection.")
	}
}
