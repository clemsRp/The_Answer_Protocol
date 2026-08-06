package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	pr "tap/protocol"
	"time"
)

func NewClient(conn net.Conn, ch chan pr.ServerResponse) *pr.Client {
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

func (s *Server) handleClient(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()
	defer func() { <-s.playerSlots }()

	responses := make(chan pr.ServerResponse, 20)

	s.wg.Add(1)
	go s.clientWriter(conn, responses)

	cli := NewClient(conn, responses)

	cli.Ch <- pr.ServerResponse{Msg: "OK hello proto=1"}
	s.entering <- cli

	input := bufio.NewScanner(conn)
	input.Buffer(make([]byte, 0, MaxPayloadSize), MaxPayloadSize)
	limiter := NewRateLimiter(5, 200*time.Millisecond)

	for input.Scan() {
		if limiter.Allow() {
			cli.Datas.Spam_warning = 0
			s.requests <- pr.ClientRequest{Cli: cli, Msg: input.Text()}
		} else {
			cli.Datas.Spam_warning++

			if cli.Datas.Spam_warning > 2 {
				cli.Ch <- pr.ServerResponse{Msg: pr.ErrSpam}
				break
			}
		}
	}

	s.leaving <- cli
}
func (s *Server) clientWriter(conn net.Conn, responses <-chan pr.ServerResponse) {
	defer s.wg.Done()
	for res := range responses {
		var output string

		if res.Datas != nil && res.Datas != "" {
			jsonData, _ := json.Marshal(res.Datas)
			output = fmt.Sprintf("%s %s\n", res.Msg, string(jsonData))
		} else {
			output = fmt.Sprintf("%s\n", res.Msg)
		}

		if _, err := conn.Write([]byte(output)); err != nil {
			return
		}
	}
}
