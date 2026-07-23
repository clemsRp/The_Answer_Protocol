package network

import (
	"bufio"
	"fmt"
	"os"
)

type Client struct {
	conn       *TimeoutConn
	id         string
	isLoggedIn bool
}

func (s *Server) readLoop(client *Client) {
	defer s.wg.Done()
	defer client.conn.Close()

	defer func() {
		s.muClients.Lock()
		delete(s.clients, client.conn)
		s.muClients.Unlock()
	}()

	scanner := bufio.NewScanner(client.conn)

	for scanner.Scan() {
		text := scanner.Text()

		s.InChan <- IncomingEvent{
			ClientID: client.id,
			Payload:  text,
		}
	}

	if err := scanner.Err(); err != nil {
		if os.IsTimeout(err) {
			fmt.Printf("Client disconnected due to inactivity (timeout of %v seconds): %s\n", client.conn.timeout.Seconds(), client.conn.RemoteAddr())
			fmt.Fprintf(client.conn, "Disconnected of the game due to inactivity (timeout of %v seconds)\n", client.conn.timeout.Seconds())
		} else {
			fmt.Println("Client connection dropped:", err)
		}
	} else {
		fmt.Println("Client disconnected gracefully:", client.conn.RemoteAddr())
		fmt.Fprintln(client.conn, "Disconnected gracefully.")
	}
}
