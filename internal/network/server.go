package network

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
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
}

func NewServer(listenAddr string, in chan<- IncomingEvent, out <-chan OutgoingEvent) *Server {
	return &Server{
		listenAddr:     listenAddr,
		quitChan:       make(chan struct{}),
		InChan:         in,
		OutChan:        out,
		timeoutSeconds: timeout_seconds,
		maxClients:     max_clients,
		maxPlayers:     max_players,
		clients:        make(map[*TimeoutConn]*Client),
	}
}

func (s *Server) Stop() {
	s.stopOnce.Do(func() {
		close(s.quitChan)
		if s.ln != nil {
			s.ln.Close()
		}
	})
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

		s.muClients.Lock()
		if len(s.clients) >= s.maxClients {
			s.muClients.Unlock()

			fmt.Println("Server full, rejecting connection:", conn.RemoteAddr())
			fmt.Fprintf(conn, "The server is currenty full. Try again later.\n")
			conn.Close()
			continue

		}
		wrappedConn := &TimeoutConn{
			Conn:    conn,
			timeout: time.Duration(s.timeoutSeconds) * time.Second,
		}
		newID := atomic.AddUint64(&s.clientCounter, 1)
		clientIDStr := fmt.Sprintf("client-%d", newID)
		newClient := &Client{conn: wrappedConn, id: clientIDStr, isLoggedIn: false}
		s.clients[wrappedConn] = newClient
		s.muClients.Unlock()

		fmt.Println("new connection to the server:", conn.RemoteAddr())

		s.wg.Add(1)
		go s.readLoop(newClient)
	}
}

func (s *Server) start() error {
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

	fmt.Println("Waiting for all connections to close...")
	s.wg.Wait()
	return nil
}

func (s *Server) SafeStart(errch chan<- error) {
	if err := s.start(); err != nil {
		log.Println("Server error:", err)
		errch <- err
	}
}
