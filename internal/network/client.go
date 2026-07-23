package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync/atomic"
)

type Client struct {
	conn *TimeoutConn
	id   string
}

func (s *Server) handleClientDisconnectMessages(client *Client, err error) {
	if err == nil {
		// the client cut the connection if err == nil from scanner.Scan().
		fmt.Println("\nClient disconnected gracefully:", client.conn.RemoteAddr())
	} else {
		if os.IsTimeout(err) {
			fmt.Printf("Client disconnected due to inactivity (timeout of %v seconds): %s", client.conn.timeout.Seconds(), client.conn.RemoteAddr())
			fmt.Fprintf(client.conn, "\nDisconnected of the game due to inactivity (timeout of %v seconds)\n", client.conn.timeout.Seconds())
		} else {
			fmt.Println("\nServer disconnected.")
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
	newClient := &Client{conn: timeoutConn, id: clientIDStr}
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

func (s *Server) clientReadLoop(client *Client) {
	defer s.wg.Done()
	defer client.conn.Close()
	defer s.removeClientByMutex(client.conn)
	scanner := bufio.NewScanner(client.conn)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		// no spam to the channel if text is empty
		if text == "" {
			continue
		}
		// here maybe serialize deserialize etc. PARSING
		s.InChan <- IncomingEvent{
			ClientID: client.id,
			Payload:  text,
		}
	}

	s.handleClientDisconnectMessages(client, scanner.Err())

}
