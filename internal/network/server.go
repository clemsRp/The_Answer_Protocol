package network

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	listenAddr     string
	ln             net.Listener
	quitChan       chan struct{}
	messageChan    chan Message
	timeoutSeconds int
	maxClients     int
	maxPlayers     int

	muClients sync.Mutex
	clients   map[net.Conn]*Client

	stopOnce sync.Once
	wg       sync.WaitGroup
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr:     listenAddr,
		quitChan:       make(chan struct{}),
		messageChan:    make(chan Message),
		timeoutSeconds: timeout_seconds,
		maxClients:     max_clients,
		maxPlayers:     max_players,
		clients:        make(map[net.Conn]*Client),
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
	close(s.messageChan)
	return nil
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
		s.clients[conn] = &Client{conn: conn, isLoggedIn: false}
		s.muClients.Unlock()

		fmt.Println("new connection to the server:", conn.RemoteAddr())

		wrappedConn := &TimeoutConn{
			Conn:    conn,
			timeout: time.Duration(s.timeoutSeconds) * time.Second,
		}
		s.wg.Add(1)
		go s.readLoop(wrappedConn)
	}
}

func (s *Server) SafeStart(errch chan<- error) {
	if err := s.start(); err != nil {
		log.Println("Server error:", err)
		errch <- err
	}
}
