package engine

import (
	"tap/engine/parser"
	pr "tap/protocol"
)

type Engine struct {
	fromServer    chan pr.ServerRequest
	toServer      chan pr.EngineResponse
	updateClients <-chan map[string]*pr.Client
	world         parser.Map
	clients       map[string]*pr.Client
	groups        map[string][]*pr.Client
	dialogues     map[string]map[string]int
	quit          chan struct{}
}

func NewEngine(world parser.Map, fromServerChan chan pr.ServerRequest, toServerChan chan pr.EngineResponse, updateClientsChan <-chan map[string]*pr.Client) *Engine {
	return &Engine{
		fromServer:    fromServerChan,
		toServer:      toServerChan,
		updateClients: updateClientsChan,
		world:         world,
		groups:        make(map[string][]*pr.Client),
		dialogues:     make(map[string]map[string]int),
		quit:          make(chan struct{}),
	}
}

func (e *Engine) Start() {
	e.broadcaster()
}
func (e *Engine) Stop() {
	select {
	case <-e.quit:
	default:
		close(e.quit)
	}
}

func (e *Engine) broadcaster() {
	for {
		select {
		case clients := <-e.updateClients:
			e.clients = clients

		case req := <-e.fromServer:
			activeCli, res, datas, err := e.handleCommands(req)
			e.toServer <- pr.EngineResponse{
				Cli:   activeCli,
				Msg:   res,
				Datas: datas,
				Err:   err,
				Req:   req,
			}
		case <-e.quit:
			return
		}

	}
}
