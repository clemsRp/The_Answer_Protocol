package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Player struct {
	Name     string
	Location string
	Hp       int
	HpMax    int
}

type Client struct {
	conn       net.Conn
	player     *Player
	isLoggedIn bool
}

type Message struct {
	from    string
	payload []byte
}

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

	muPlayers sync.Mutex
	players   map[string]*Player

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
	}
}

func (s *Server) Start() error {
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
		s.clients[conn] = &Client{conn: conn, player: nil, isLoggedIn: false}
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

		s.messageChan <- Message{
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

func (s *Server) messageLoop() {

	for msg := range s.messageChan {
		payload := strings.TrimSpace(string(msg.payload))
		// parts := strings.Split(payload, " ")
		fmt.Printf("received message from connection (%s): %s\n", msg.from, payload)
		// if strings.HasPrefix(payload, "CONNECT") && len(parts) == 2 {

		// }
	}
}

func (s *Server) safeStart(errch chan<- error) {
	if err := s.Start(); err != nil {
		log.Println("Server error:", err)
		errch <- err
	}
}

func main() {
	server := NewServer(":8080")
	// to wait for CTRL+C signal or KILL signal ASYNCHRONOUSLY with a buffer (not blocking the thread)
	sigch := make(chan os.Signal, 1)
	// to wait for an error at the launch of the server same way as sigch.
	errch := make(chan error, 1)
	go server.messageLoop()

	go server.safeStart(errch)

	// The only way to stop the server is by CTRL+C or KILL command
	// then it stops gracefully.
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-sigch:
		fmt.Printf("\nStop signal received (%v), closing the server...", sig)
	case err := <-errch:
		fmt.Printf("\nCritical error from server: %v", err)
	}

	server.Stop()
	fmt.Println("\nServer stopped graciously.")
}
