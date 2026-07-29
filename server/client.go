package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	pr "tap/protocol"
	"time"
)

func (s *Server) handleClient(conn net.Conn) {
	// Send server responses to client terminal
	responses := make(chan pr.ServerResponse)
	go s.clientWriter(conn, responses)

	// Init client
	who := conn.RemoteAddr().String()

	cli := &pr.Client{
		Conn: conn,
		Ch:   responses,
		Ip:   who,
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

	// Log player connection
	cli.Ch <- pr.ServerResponse{Msg: "OK hello proto=1"}
	s.entering <- cli

	// Handle player commands user buffers
	input := bufio.NewScanner(conn)
	buf := make([]byte, 0, 1024)
	input.Buffer(buf, 1024)

	for input.Scan() {
		// Handle valid commands
		now := time.Now()
		// if now.Sub(cli.Datas.Last_cmd_time) >= 200*time.Millisecond {
		cli.Datas.Spam_warning = 0
		cli.Datas.Last_cmd_time = now
		s.requests <- pr.ClientRequest{Cli: cli, Msg: input.Text()}

		// Handle command spams
		// }
		//  else {
		// 	cli.Datas.Spam_warning++
		// 	cli.Datas.Last_cmd_time = now

		// 	s.LogWarn("Abuse detected", map[string]any{"warnings": cli.Datas.Spam_warning, "ip": cli.Ip, "player": cli.Name})

		// 	if cli.Datas.Spam_warning > 3 {
		// 		cli.Ch <- pr.ServerResponse{Msg: "ERR 900 CONNECTION_CLOSED_DUE_TO_SPAM", Req: pr.ClientRequest{Cli: cli, Msg: input.Text()}}
		// 		break
		// 	}
		// }
	}

	s.leaving <- cli
}

func (s *Server) clientWriter(conn net.Conn, responses <-chan pr.ServerResponse) {
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
