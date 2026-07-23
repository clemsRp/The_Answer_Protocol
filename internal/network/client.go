package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Client struct {
	conn       net.Conn
	isLoggedIn bool
}

type Message struct {
	from    string
	payload []byte
}

func (s *Server) readLoop(conn *TimeoutConn) {
	defer s.wg.Done()
	defer conn.Close()

	defer func() {
		s.muClients.Lock()
		delete(s.clients, conn)
		s.muClients.Unlock()
	}()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		text := scanner.Text()

		s.InChan <- IncomingEvent{
			from:    conn.RemoteAddr().String(),
			payload: []byte(text),
		}
	}

	if err := scanner.Err(); err != nil {
		if os.IsTimeout(err) {
			fmt.Printf("Client disconnected due to inactivity (timeout of %v seconds): %s\n", conn.timeout.Seconds(), conn.RemoteAddr())
			fmt.Fprintf(conn, "Disconnected of the game due to inactivity (timeout of %v seconds)\n", conn.timeout.Seconds())
		} else {
			fmt.Println("Client connection dropped:", err)
		}
	} else {
		fmt.Println("Client disconnected gracefully:", conn.RemoteAddr())
		fmt.Fprintln(conn, "Disconnected gracefully.")
	}
}
