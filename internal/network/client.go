package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync/atomic"
)

type Client struct {
	conn       *TimeoutConn
	id         string
	isLoggedIn bool
}

func (s *Server) clientReadLoop(client *Client) {
	defer s.wg.Done()
	defer client.conn.Close()
	defer s.removeClientByMutex(client.conn)
	scanner := bufio.NewScanner(client.conn)

	for scanner.Scan() {
		text := scanner.Text()

		s.InChan <- IncomingEvent{
			ClientID: client.id,
			Payload:  text,
		}
	}

	s.handleClientDisconnectMessages(client, scanner.Err())

}

func (s *Server) handleClientDisconnectMessages(client *Client, err error) {
	if err == nil {
		fmt.Println("Client disconnected gracefully:", client.conn.RemoteAddr())
		fmt.Fprintln(client.conn, "Disconnected gracefully.")

	} else {
		if os.IsTimeout(err) {
			fmt.Printf("Client disconnected due to inactivity (timeout of %v seconds): %s\n", client.conn.timeout.Seconds(), client.conn.RemoteAddr())
			fmt.Fprintf(client.conn, "Disconnected of the game due to inactivity (timeout of %v seconds)\n", client.conn.timeout.Seconds())
		} else {
			fmt.Println("Client connection dropped:", err)
		}
	}
}

func (s *Server) createClientID() string {
	newID := atomic.AddUint64(&s.clientCounter, 1)
	clientIDStr := fmt.Sprintf("client-%d", newID)

	return clientIDStr
}

func (s *Server) createNewClient(conn net.Conn) *Client {
	timeoutConn := s.createTimeoutConn(conn)
	clientIDStr := s.createClientID()
	newClient := &Client{conn: timeoutConn, id: clientIDStr, isLoggedIn: false}
	return newClient
}

func (s *Server) addNewClientByMutex(client *Client) {
	s.muClients.Lock()
	s.clients[client.conn] = client
	s.muClients.Unlock()
	fmt.Println("new connection to the server:", client.conn.RemoteAddr())
}

func (s *Server) removeClientByMutex(timeout_conn *TimeoutConn) {
	s.muClients.Lock()
	delete(s.clients, timeout_conn)
	s.muClients.Unlock()
	// inform game engine maybe here?
	// s.InChan <- IncomingEvent{ClientID: client.id, Type: "DISCONNECT"}
}

func (s *Server) isMaxClientLimitByMutex() bool {
	s.muClients.Lock()
	defer s.muClients.Unlock()

	if len(s.clients) >= s.maxClients {
		return true
	}
	return false
}
