package server

import (
	"fmt"
	pr "tap/protocol"
	"time"

	"net"
)

type Log struct {
	Timestamp string         `json:"timestamp"`
	Level     string         `json:"level"`
	Message   string         `json:"message"`
	Datas     map[string]any `json:"datas,omitempty"`
}

type Server struct {
	ln      net.Listener
	clients map[string]*pr.Client

	requests chan pr.ClientRequest
	logs     chan Log

	entering chan *pr.Client
	leaving  chan *pr.Client

	toEngine      chan pr.ServerRequest
	fromEngine    chan pr.EngineResponse
	updateClients chan map[string]*pr.Client

	t_start int64
}

func NewServer(listenAddr string, toEngineChan chan pr.ServerRequest, fromEngineChan chan pr.EngineResponse, updateClientsChan chan map[string]*pr.Client) (*Server, error) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("An error occured during the serveur connection:")
		return nil, err
	}
	return &Server{
		ln:            listener,
		clients:       make(map[string]*pr.Client),
		requests:      make(chan pr.ClientRequest, 100),
		logs:          make(chan Log, 500),
		entering:      make(chan *pr.Client, 50),
		leaving:       make(chan *pr.Client, 50),
		toEngine:      toEngineChan,
		fromEngine:    fromEngineChan,
		updateClients: updateClientsChan,
		t_start:       time.Now().Unix(),
	}, nil
}

func (s *Server) Start() {
	s.LogInfo("Connection established on :8080.", nil)
	go s.broadcaster()

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			s.LogError("An error occurred during a client's connection", map[string]any{
				"error": err.Error(),
			})
			continue
		}
		go s.handleClient(conn)
	}
}

func (s *Server) broadcaster() {

	for {
		select {
		// Print server logs
		case log := <-s.logs:
			s.writeLog(log.Level, log.Message, log.Datas)

		// Handle a client command
		case req := <-s.requests:
			s.toEngine <- pr.ServerRequest{Msg: req.Msg, Cli: req.Cli, Req: req}

		// Handle the start of the a client connection
		case cli := <-s.entering:
			// Add client
			s.clients[cli.Ip] = cli
			s.LogInfo("Start of connection", map[string]any{
				"ip":       cli.Ip,
				"duration": s.get_timestamp(),
			})

			snapshot := make(map[string]*pr.Client, len(s.clients))
			for k, v := range s.clients {
				snapshot[k] = v
			}

			// Update engine clients
			s.updateClients <- snapshot

		// Handle the end of the a client connection
		case cli := <-s.leaving:
			// Delete client
			if c, ok := s.clients[cli.Ip]; ok {
				delete(s.clients, cli.Ip)
				close(c.Ch)
			}
			s.LogInfo("End of connection", map[string]any{
				"ip":       cli.Ip,
				"duration": s.get_timestamp(),
				"player":   cli.Name,
			})

			// Update engine clients
			s.updateClients <- s.clients

		// Handle Request Output
		case output := <-s.fromEngine:

			// Handle command errors and status codes
			res := output.Msg
			datas := output.Datas
			if output.Err != nil {
				res, datas = output.Err.Error(), ""
			}

			// Send server response to the client
			output.Cli.Ch <- pr.ServerResponse{Msg: res, Datas: datas, Req: output.Req.Req}
		}
	}
}
