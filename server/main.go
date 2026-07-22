package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr:     listenAddr,
		quitChan:       make(chan struct{}),
		messageChan:    make(chan Message),
		timeoutSeconds: 5,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln
	fmt.Println("Server started on the port :8080")
	go s.acceptLoop()

	<-s.quitChan
	return nil
}

func (s *Server) Stop() {
	close(s.quitChan)
	if s.ln != nil {
		s.ln.Close()
	}
}

func (s *Server) acceptLoop() {
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
		fmt.Println("new connection to the server:", conn.RemoteAddr())

		wrappedConn := &TimeoutConn{
			Conn:    conn,
			timeout: time.Duration(s.timeoutSeconds) * time.Second,
		}

		go s.readLoop(wrappedConn)
	}
}

func (s *Server) readLoop(conn *TimeoutConn) {
	defer conn.Close()

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
		} else {
			fmt.Println("Client connection dropped:", err)
		}
	} else {
		fmt.Println("Client disconnected gracefully:", conn.RemoteAddr())
	}
}

func (s *Server) messageLoop() {

	for msg := range s.messageChan {
		fmt.Printf("received message from connection (%s): %s\n", msg.from, string(msg.payload))
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
