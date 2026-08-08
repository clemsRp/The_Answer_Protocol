package server

import (
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
	clients map[string]*pr.Client

	requests chan pr.ClientRequest
	logs     chan Log

	entering    chan *pr.Client
	leaving     chan *pr.Client
	quit        chan struct{}
	playerSlots chan struct{}

	wg sync.WaitGroup

	toEngine      chan pr.ServerRequest
	fromEngine    chan pr.EngineResponse
	updateClients chan map[string]*pr.Client

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
			s.toEngine <- pr.ServerRequest{Msg: req.Msg, Cli: req.Cli, Req: req}
		case cli := <-s.entering:
			s.addClient(cli)
		case cli := <-s.leaving:
			s.removeClient(cli)
		case output := <-s.fromEngine:
			s.sendToClient(output)
		case <-s.quit:
			s.shutdown()
			return
		}
	}
}

func (s *Server) addClient(cli *pr.Client) {
	s.clients[cli.Ip] = cli
	s.sendClientsUpdate()
}

func (s *Server) removeClient(cli *pr.Client) {
	if c, ok := s.clients[cli.Ip]; ok {
		delete(s.clients, cli.Ip)
		close(c.Ch)
	}
	s.sendClientsUpdate()
}

func (s *Server) sendClientsUpdate() {
	snapshot := make(map[string]*pr.Client, len(s.clients))
	for k, v := range s.clients {
		snapshot[k] = v
	}
	s.updateClients <- snapshot
}

func (s *Server) sendToClient(output pr.EngineResponse) {

	if _, exists := s.clients[output.Cli.Ip]; !exists {
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
	case output.Cli.Ch <- pr.ServerResponse{Msg: res, Datas: datas}:
	default:
		output.Cli.Conn.Close()
	}
}

func (s *Server) shutdown() {
	for _, c := range s.clients {
		close(c.Ch)
		c.Conn.Close()
	}
}

func (s *Server) GetAddress() string {
	return s.ln.Addr().String()
}

func NewServer(listenAddr string, toEngineChan chan pr.ServerRequest, fromEngineChan chan pr.EngineResponse, updateClientsChan chan map[string]*pr.Client) (*Server, error) {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	return &Server{
		ln:            listener,
		clients:       make(map[string]*pr.Client),
		requests:      make(chan pr.ClientRequest, 100),
		logs:          make(chan Log, 500),
		entering:      make(chan *pr.Client, 50),
		leaving:       make(chan *pr.Client, 50),
		playerSlots:   make(chan struct{}, MaxPlayerLimit),
		quit:          make(chan struct{}),
		toEngine:      toEngineChan,
		fromEngine:    fromEngineChan,
		updateClients: updateClientsChan,
		t_start:       time.Now().Unix(),
		IdleTimeout:   time.Duration(MaxClientTimeOutMinutes) * time.Minute,
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
