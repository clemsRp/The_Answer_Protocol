package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	pr "tap/protocol"
	"time"
)

func newClient(conn net.Conn, ch chan pr.ServerResponse) *pr.Client {
	return &pr.Client{
		Conn: conn,
		Ch:   ch,
		Ip:   conn.RemoteAddr().String(),
		Datas: pr.Datas{
			Room:          "entrance",
			Status:        "healthy",
			Promotion:     false,
			Hp:            100,
			Max_hp:        100,
			Connected:     false,
			Last_cmd_time: time.Now(),
		},
	}
}

func (s *Server) stopClient(cli *pr.Client, conn net.Conn, writerDone <-chan struct{}) {
	select {
	case s.leaving <- cli:
	case <-s.quit:
	}

	<-writerDone
	conn.Close()
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
	defer s.stopClient(cli, conn, writerDone)

	s.wg.Add(1)
	go s.clientWriter(conn, responses, writerDone)

	cli.Ch <- pr.ServerResponse{Msg: "OK hello proto=1"}
	s.entering <- cli

	s.readClientInput(cli, conn)
}

func (s *Server) readClientInput(cli *pr.Client, conn net.Conn) {
	input := bufio.NewScanner(conn)
	input.Buffer(make([]byte, 0, MaxPayloadSize), MaxPayloadSize)
	limiter := NewRateLimiter(5, 200*time.Millisecond)

	for {
		conn.SetReadDeadline(time.Now().Add(s.IdleTimeout))

		if !input.Scan() {
			break
		}
		if limiter.Allow() {
			cli.Datas.Spam_warning = 0
			s.requests <- pr.ClientRequest{Cli: cli, Msg: input.Text()}
		} else {
			if s.handleSpam(cli) {
				break
			}
		}
	}
}

func (s *Server) handleSpam(cli *pr.Client) bool {
	cli.Datas.Spam_warning++
	if cli.Datas.Spam_warning > 2 {
		cli.Ch <- pr.ServerResponse{Msg: pr.ErrSpam}
		return true
	}
	return false
}

func (s *Server) clientWriter(conn net.Conn, responses <-chan pr.ServerResponse, done chan struct{}) {
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
			if _, err := conn.Write([]byte(output)); err != nil {
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
