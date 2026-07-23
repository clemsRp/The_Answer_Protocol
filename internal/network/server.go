package network

import (
	"errors"
	"fmt"
	"log"
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

func NewServer(listenAddr string, in chan<- IncomingEvent, out <-chan OutgoingEvent) *Server {
	return &Server{
		listenAddr: listenAddr,
		InChan:     in,
		OutChan:    out,

		quitChan: make(chan struct{}),
		ErrChan:  make(chan error, 1),
		SigChan:  make(chan os.Signal, 1),

		timeoutSeconds: timeout_seconds,
		maxClients:     max_clients,
		maxPlayers:     max_players,
		clients:        make(map[*TimeoutConn]*Client),
	}
}

func (s *Server) acceptLoop() {
	defer s.wg.Done()
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Println("acceptLoop temporary error:", err)
			time.Sleep(1 * time.Second)
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
		go s.clientReadLoop(newClient)
	}
}

func (s *Server) run() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	s.ln = ln
	fmt.Println("Server started on the port :8080")
	s.wg.Add(1)
	go s.acceptLoop()

	// loops infinitely while no close signal via quitChan
	<-s.quitChan

	return nil
}

func (s *Server) Start() {

	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		if err := s.run(); err != nil {
			log.Println("Server error:", err)
			s.ErrChan <- err
		}
	}()
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
