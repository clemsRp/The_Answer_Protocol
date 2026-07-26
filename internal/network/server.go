package network

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type IncomingEvent struct {
	ClientID string
	Payload  string
}

type OutgoingEvent struct {
	ClientID string
	Payload  string
}

type Server struct {
	listenAddr     string
	ln             net.Listener
	quitChan       chan struct{}
	timeoutSeconds int
	maxClients     int
	maxPlayers     int
	clientCounter  uint64

	muClients sync.Mutex
	clients   map[*TimeoutConn]*Client

	stopOnce sync.Once
	wg       sync.WaitGroup
	InChan   chan<- IncomingEvent // server sends commands from the clients ->
	// to the game
	OutChan <-chan OutgoingEvent // server receives events from game, ->
	// and sends them to clients
	ErrChan chan error // Channel to listen for listen error of server, to quit gracefully
	SigChan chan os.Signal
}

func NewServer(listenAddr string, in chan<- IncomingEvent, out <-chan OutgoingEvent) (*Server, error) {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, fmt.Errorf("Couldn't start the TCP server: %w", err)
	}
	fmt.Printf("Server started on %s\n", listener.Addr().String())
	return &Server{
		listenAddr: listenAddr,
		ln:         listener,
		InChan:     in,
		OutChan:    out,

		quitChan: make(chan struct{}),
		ErrChan:  make(chan error, 1),
		SigChan:  make(chan os.Signal, 1),

		timeoutSeconds: timeout_seconds,
		maxClients:     max_clients,
		maxPlayers:     max_players,
		clients:        make(map[*TimeoutConn]*Client),
	}, nil
}

func (s *Server) acceptLoop() {
	defer s.wg.Done()
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			if s.handleAcceptError(err) {
				return
			}
			continue
		}
		if s.isMaxClientLimitByMutex() {
			fmt.Println("Server full, rejecting connection:", conn.RemoteAddr())
			fmt.Fprintf(conn, "The server is currenty full. Try again later.\n")
			conn.Close()
			continue
		}
		newClient := s.createNewClient(conn)
		s.addNewClientByMutex(newClient)
		s.wg.Add(1)
		go s.handleConnection(newClient)
	}
}

func (s *Server) handleAcceptError(err error) bool {
	if errors.Is(err, net.ErrClosed) {
		return true
	} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		fmt.Println("acceptLoop temporary error:", err)
		time.Sleep(1 * time.Second)
		return false
	} else {
		// server crash PROBABLY
		s.ErrChan <- fmt.Errorf("fatal error from server: %w", err)
		return true
	}
}

func (s *Server) GetAddress() string {
	return s.ln.Addr().String()
}

func (s *Server) Start() {

	s.wg.Add(1)
	go s.acceptLoop()
}

func (s *Server) disconnectAllClientsByMutex() {
	fmt.Println("\nClosing all connections to clients...")

	s.muClients.Lock()
	for conn := range s.clients {
		conn.Close()
	}
	s.muClients.Unlock()
}

func (s *Server) Stop() {
	s.stopOnce.Do(func() {
		close(s.quitChan)
		if s.ln != nil {
			s.ln.Close()
		}
		s.disconnectAllClientsByMutex()
	})
	s.wg.Wait()
	fmt.Println("\nServer stopped graciously.")

}
