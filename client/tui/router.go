package tui

import (
	"strings"
	pr "tap/protocol"
)

type Router struct {
	Inputs  chan<- string
	Outputs <-chan pr.ServerResponse

	ChatChan        chan pr.ServerResponse
	CommandLineChan chan pr.ServerResponse
	ServerChan      chan pr.ServerResponse
	NavChan         chan pr.ServerResponse
	ItemsChan       chan pr.ServerResponse
	InteractionChan chan pr.ServerResponse
	DatasChan       chan pr.ServerResponse
	LastCommand     string
}

func NewRouter(inputs chan string, outputs <-chan pr.ServerResponse) *Router {
	return &Router{
		Inputs:          inputs,
		Outputs:         outputs,
		ChatChan:        make(chan pr.ServerResponse, 100),
		CommandLineChan: make(chan pr.ServerResponse, 100),
		ServerChan:      make(chan pr.ServerResponse, 100),
		NavChan:         make(chan pr.ServerResponse, 100),
		ItemsChan:       make(chan pr.ServerResponse, 100),
		InteractionChan: make(chan pr.ServerResponse, 100),
		DatasChan:       make(chan pr.ServerResponse, 100),
		LastCommand:     "",
	}
}

func (r *Router) HandleEvents(res pr.ServerResponse) {
	// Handle CHAT responses
	global := strings.HasPrefix(res.Msg, "EVT GLOBAL CHAT")
	room := strings.HasPrefix(res.Msg, "EVT ROOM CHAT")
	group := strings.HasPrefix(res.Msg, "EVT GROUP CHAT")
	if global || room || group {
		r.ChatChan <- res
	}

	// Handle PRESENCE responses
	if strings.HasPrefix(res.Msg, "EVT ROOM PRESENCE") {
		r.InteractionChan <- res
	}
}

func (r *Router) Start() {
	go func() {

		for res := range r.Outputs {
			switch {
			case strings.HasPrefix(res.Msg, "OK"):
				r.handleLastCommandResponse(res)
			// TO DO
			case strings.HasPrefix(res.Msg, "ERR"):
				r.handleLastCommandResponse(res)

			case strings.HasPrefix(res.Msg, "EVT"):
				r.HandleEvents(res)
			}

			r.ServerChan <- res
			r.DatasChan <- res
			r.CommandLineChan <- res
		}
	}()
}

func (m *MyApp) UpdateRouterLayout() {
	if m.grid != nil {
		m.grid.Clear()
	}

	if m.Navigation != nil {
		m.grid.AddItem(m.Navigation.Layout, 0, 0, 1, 1, 0, 0, true)
	}
}

func (r *Router) handleLastCommandResponse(res pr.ServerResponse) {
	switch r.LastCommand {
	case pr.CmdLook:
		r.NavChan <- res
		r.ItemsChan <- res
		r.InteractionChan <- res
	case pr.CmdChat:
		r.ChatChan <- res
	default:

	}
}
