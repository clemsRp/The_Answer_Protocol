package engine

import (
	"fmt"
	pr "tap/protocol"
)

type Engine struct {
	world *Map
	//sessions store id to pseudos
	sessions map[string]string
	//players store pseudos to Player
	players       map[string]*Player
	groups        map[string]*Group
	dialogues     map[string]map[string]int
	activeCombats map[string]*CombatSession
	exchanger     pr.Exchanger
	quit          chan struct{}
}

func NewEngine(world *Map, exchanger pr.Exchanger) *Engine {
	return &Engine{
		world:         world,
		sessions:      make(map[string]string),
		players:       make(map[string]*Player),
		groups:        make(map[string]*Group),
		dialogues:     make(map[string]map[string]int),
		activeCombats: make(map[string]*CombatSession),
		quit:          make(chan struct{}),
		exchanger:     exchanger,
	}
}

func (e *Engine) Start() error {
	err := e.broadcaster()
	return err
}
func (e *Engine) Stop() {
	select {
	case <-e.quit:
	default:
		close(e.quit)
	}
}

func (e *Engine) broadcaster() (errPanic error) {
	defer func() {
		if r := recover(); r != nil {
			errPanic = fmt.Errorf("Engine panic: %v", r)
		}
	}()
	for {
		select {
		case id := <-e.exchanger.JoinChan:
			// ghost session stored in sessions when a new client arrives
			e.sessions[id] = ""
		case id := <-e.exchanger.LeaveChan:
			e.handlePlayerLeave(id)

		case req := <-e.exchanger.ServerInput:
			res, datas, err := e.handleCommands(req)
			e.exchanger.ServerOutput <- pr.EngineResponse{
				Id:    req.Id,
				Msg:   res,
				Datas: datas,
				Err:   err,
			}
		case <-e.quit:
			return nil
		}

	}
}

func (e *Engine) handlePlayerLeave(id string) {
	pseudo := e.sessions[id]

	if pseudo != "" {
		if player, exists := e.players[pseudo]; exists {
			e.playerQuits(player)
			delete(e.players, pseudo)
		}
	}
	delete(e.sessions, id)
}
