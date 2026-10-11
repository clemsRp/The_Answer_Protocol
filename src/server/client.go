package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"
	pr "tap/src/protocol"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	conn        net.Conn               `json:"-"`
	ch          chan pr.ServerResponse `json:"-"`
	id          string
	ip          string
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
		ip:   conn.RemoteAddr().String(),
	}
}

func (s *Server) stopClient(cli *Client, writerDone <-chan struct{}) {
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
	limiter := NewRateLimiter(MaxTokens, 35*time.Millisecond)

	for {
		cli.conn.SetReadDeadline(time.Now().Add(s.IdleTimeout))

		if !input.Scan() {
			if errors.Is(input.Err(), bufio.ErrTooLong) {
				select {
				case cli.ch <- pr.ServerResponse{Msg: pr.ErrInvalidCommand}:
				case <-s.quit:
				default:
				}
			}
			break
		}

		if strings.Contains(input.Text(), "\\n") {
			msgs := strings.Split(input.Text(), "\\n")
			for _, msg := range msgs {
				if s.HandleInput(cli, limiter, msg) {
					break
				}
			}

		} else {
			if s.HandleInput(cli, limiter, input.Text()) {
				break
			}
		}
	}
}

func (s *Server) HandleInput(cli *Client, limiter *RateLimiter, msg string) bool {
	if limiter.Allow(msg) {
		cli.spamWarning = 0
		select {
		case s.requests <- ClientRequest{id: cli.id, msg: msg}:
		case <-s.quit:
			return true
		}
	} else if s.handleSpam(cli) {
		return true
	}

	return false
}

func (s *Server) handleSpam(cli *Client) bool {
	cli.spamWarning++
	msg := pr.WarnSpam
	if cli.spamWarning >= 5 {
		msg = pr.ErrSpam
	}

	select {
	case cli.ch <- pr.ServerResponse{Msg: msg}:
	case <-s.quit:
		return true
	default:
	}

	if cli.spamWarning < 5 {
		slog.Warn("Command spam", "client_id", cli.id, "ip", cli.ip, "result", pr.WarnSpam)
		return false
	}
	slog.Error("Command spam", "client_id", cli.id, "ip", cli.ip, "result", pr.ErrSpam)
	return true
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
			if res.Msg == "OK bye" {
				cli.conn.Close()
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
