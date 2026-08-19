package server

import (
	"log/slog"
	"sync"
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
	clients map[string]*Client

	requests chan ClientRequest
	logs     chan Log

	entering    chan *Client
	leaving     chan *Client
	quit        chan struct{}
	playerSlots chan struct{}

	wg sync.WaitGroup

	exchanger pr.Exchanger

	t_start     int64
	IdleTimeout time.Duration
}

func (s *Server) handleNewConnection(conn net.Conn) {
	// SELECT PATTERN: Non-Blocking Send / Load Shedding
	// Attempts to consume a player slot.
	// If a slot is available, the client is connected.
	// Otherwise (default), the server is full, and we instantly reject the connection without blocking.
	select {
	case s.playerSlots <- struct{}{}:
		s.wg.Add(1)
		go s.handleClient(conn)
	default:
		conn.Write([]byte(pr.ErrServerFull + "\n"))
		conn.Close()
	}
}

func (s *Server) broadcaster() {
	defer s.wg.Done()
	// PATTERN: Event Loop / Multiplexer
	// The central dispatcher of the server: guarantees that only one event (connection,
	// disconnection, message) is processed at a time. This protects the s.clients map
	// from Data Races without the need for sync.Mutex locks.
	for {
		select {
		case req := <-s.requests:
			slog.Info("Client request", "command", req.msg)
			s.exchanger.ServerInput <- pr.ServerRequest{Id: req.id, Msg: req.msg}

		case cli := <-s.entering:
			slog.Info("Client connection", "ip", cli.ip)
			s.addClient(cli)
			s.exchanger.JoinChan <- cli.id

		case cli := <-s.leaving:
			slog.Info("Client disconnection", "ip", cli.ip)
			s.removeClient(cli)
			s.exchanger.LeaveChan <- cli.id

		case output := <-s.exchanger.ServerOutput:
			s.sendToClient(output)

		case <-s.quit:
			s.shutdown()
			return
		}
	}
}

func (s *Server) addClient(cli *Client) {
	s.clients[cli.id] = cli
}

func (s *Server) removeClient(cli *Client) {
	if c, ok := s.clients[cli.id]; ok {
		delete(s.clients, cli.id)
		close(c.ch)
	}
}

func (s *Server) sendToClient(output pr.EngineResponse) {

	cli, exists := s.clients[output.Id]
	if !exists {
		return
	}

	res, datas := output.Msg, output.Datas
	if output.Err != nil {
		res, datas = output.Err.Error(), ""
	}

	// SELECT PATTERN: Non-Blocking Send
	// Attempts to send the message to the client. If their channel buffer is full
	// we force a disconnection via the default case
	// rather than paralyzing the entire server broadcast loop.
	select {
	case cli.ch <- pr.ServerResponse{Msg: res, Datas: datas}:
	default:
		cli.conn.Close()
	}
}

func (s *Server) shutdown() {
	for _, c := range s.clients {
		close(c.ch)
		c.conn.Close()
	}
}

func (s *Server) GetAddress() string {
	return s.ln.Addr().String()
}

func NewServer(listenAddr string, exchanger pr.Exchanger) (*Server, error) {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	return &Server{
		ln:          listener,
		clients:     make(map[string]*Client),
		requests:    make(chan ClientRequest, 100),
		logs:        make(chan Log, 500),
		entering:    make(chan *Client, 50),
		leaving:     make(chan *Client, 50),
		playerSlots: make(chan struct{}, MaxPlayerLimit),
		quit:        make(chan struct{}),

		exchanger:   exchanger,
		t_start:     time.Now().Unix(),
		IdleTimeout: time.Duration(MaxClientTimeOutMinutes) * time.Minute,
	}, nil
}

func (s *Server) Start() {

	s.wg.Add(1)
	defer s.wg.Done()

	s.wg.Add(1)
	go s.broadcaster()

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return
		}
		s.handleNewConnection(conn)
	}
}

func (s *Server) Stop() {

	if s.ln != nil {
		s.ln.Close()
	}
	// if the channel already has been closed
	// does nothing. else it closes it
	select {
	case <-s.quit:
	default:
		close(s.quit)
	}
	s.wg.Wait()
}
