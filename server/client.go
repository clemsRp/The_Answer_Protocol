package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	pr "tap/protocol"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	conn        net.Conn               `json:"-"`
	ch          chan pr.ServerResponse `json:"-"`
	id          string
	spamWarning int
}

type ClientRequest struct {
	id  string
	msg string
}

func newClient(conn net.Conn, ch chan pr.ServerResponse) *Client {
	return &Client{
		conn: conn,
		ch:   ch,
		id:   uuid.New().String(),
	}
}

func (s *Server) stopClient(cli *Client, writerDone <-chan struct{}) {
	// SELECT PATTERN: Escape Hatch / Circuit Breaker
	// Attempts to send the client to the leaving queue.
	// If the server is shutting down (s.quit is closed), the send operation is
	// immediately aborted to prevent a deadlock.
	select {
	case s.leaving <- cli:
	case <-s.quit:
	}

	<-writerDone
	cli.conn.Close()
	<-s.playerSlots
}

func (s *Server) handleClient(conn net.Conn) {
	defer s.wg.Done()
	responses := make(chan pr.ServerResponse, 20)
	cli := newClient(conn, responses)
	writerDone := make(chan struct{})

	// when connection interrupts,
	// closes the client channel and removes him from the clients map
	// waits for client writer ends sending its messages.
	defer s.stopClient(cli, writerDone)

	s.wg.Add(1)
	go s.clientWriter(cli, responses, writerDone)

	cli.ch <- pr.ServerResponse{Msg: "OK hello proto=1"}
	s.entering <- cli

	s.readClientInput(cli)
}

func (s *Server) readClientInput(cli *Client) {
	input := bufio.NewScanner(cli.conn)
	input.Buffer(make([]byte, 0, MaxPayloadSize), MaxPayloadSize)
	limiter := NewRateLimiter(MaxTokens, 200*time.Millisecond)

	for {
		cli.conn.SetReadDeadline(time.Now().Add(s.IdleTimeout))

		if !input.Scan() {
			break
		}
		if limiter.Allow() {
			cli.spamWarning = 0
			s.requests <- ClientRequest{id: cli.id, msg: input.Text()}
		} else {
			if s.handleSpam(cli) {
				break
			}
		}
	}
}

func (s *Server) handleSpam(cli *Client) bool {
	cli.spamWarning++
	if cli.spamWarning > 2 {
		cli.ch <- pr.ServerResponse{Msg: pr.ErrSpam}
		return true
	}
	return false
}

func (s *Server) clientWriter(cli *Client, responses <-chan pr.ServerResponse, done chan struct{}) {
	defer s.wg.Done()
	defer close(done)

	for {
		select {
		case <-s.quit:
			return

		case res, ok := <-responses:
			if !ok {
				return
			}
			output := formatResponse(res)
			if _, err := cli.conn.Write([]byte(output)); err != nil {
				return
			}
		}
	}
}
func formatResponse(res pr.ServerResponse) string {
	if res.Datas != nil && res.Datas != "" {
		jsonData, _ := json.Marshal(res.Datas)
		return fmt.Sprintf("%s %s\n", res.Msg, string(jsonData))
	}
	return fmt.Sprintf("%s\n", res.Msg)
}
