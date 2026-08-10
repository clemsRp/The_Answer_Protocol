package engine

import (
	"tap/engine/parser"
	pr "tap/protocol"
)

type Engine struct {
	world parser.Map
	//sessions store IP to pseudos
	sessions map[string]string
	//players store pseudos to Player
	players   map[string]*Player
	groups    map[string][]*Player
	dialogues map[string]map[string]int
	exchanger pr.Exchanger
	quit      chan struct{}
}

func NewEngine(world parser.Map, exchanger pr.Exchanger) *Engine {
	return &Engine{
		exchanger: exchanger,
		world:     world,
		groups:    make(map[string][]*Player),
		dialogues: make(map[string]map[string]int),
		quit:      make(chan struct{}),
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
		case ip := <-e.exchanger.JoinChan:
			// ghost session stored in sessions when a new client arrives
			e.sessions[ip] = ""
		case ip := <-e.exchanger.LeaveChan:
			e.handlePlayerLeave(ip)

		case req := <-e.exchanger.ServerInput:
			res, datas, err := e.handleCommands(req)
			e.exchanger.ServerOutput <- pr.EngineResponse{
				Ip:    req.Ip,
				Msg:   res,
				Datas: datas,
				Err:   err,
			}
		case <-e.quit:
			return
		}

	}
}

func (e *Engine) handlePlayerLeave(ip string) {
	pseudo := e.sessions[ip]

	if pseudo != "" {
		if player, exists := e.players[pseudo]; exists {
			e.playerQuits(player)
			delete(e.players, pseudo)
		}
	}
	delete(e.sessions, ip)
}
