package tui

import (
	"strings"
	panel "tap/client/tui/panels"
	pr "tap/protocol"
)

type Router struct {
	Inputs  chan<- string
	Outputs <-chan pr.ServerResponse

	ChatChan        chan pr.ServerResponse
	CommandLineChan chan pr.ServerResponse
	ServerChan      chan pr.ServerResponse
	NavChan         chan pr.ServerResponse
	GroupChan       chan pr.ServerResponse
	GroupLeaveChan  chan pr.ServerResponse
	UsersChan       chan pr.ServerResponse
	ItemsChan       chan pr.ServerResponse
	InteractionChan chan pr.ServerResponse
	DatasChan       chan pr.ServerResponse
	QuitChan        <-chan struct{}

	LastCommand string
}

func NewRouter(inputs chan string, outputs <-chan pr.ServerResponse, quitChan <-chan struct{}) *Router {

	return &Router{
		Inputs:          inputs,
		Outputs:         outputs,
		ChatChan:        make(chan pr.ServerResponse, 100),
		CommandLineChan: make(chan pr.ServerResponse, 100),
		ServerChan:      make(chan pr.ServerResponse, 100),
		NavChan:         make(chan pr.ServerResponse, 100),
		GroupChan:       make(chan pr.ServerResponse, 100),
		GroupLeaveChan:  make(chan pr.ServerResponse, 100),
		UsersChan:       make(chan pr.ServerResponse, 100),
		ItemsChan:       make(chan pr.ServerResponse, 100),
		InteractionChan: make(chan pr.ServerResponse, 100),
		DatasChan:       make(chan pr.ServerResponse, 100),
		QuitChan:        quitChan,

		LastCommand: "",
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

	// Handle ITEM responses
	if strings.HasPrefix(res.Msg, "EVT ITEM") {
		r.ItemsChan <- res
	}

	// Handle PRESENCE responses
	if strings.HasPrefix(res.Msg, "EVT ROOM PRESENCE") {
		r.InteractionChan <- res
	}

	// Handle GROUP responses
	if strings.HasPrefix(res.Msg, "EVT GROUP") {
		r.GroupChan <- res
	}
}

func (r *Router) Start() {
	for {
		select {
		case res := <-r.Outputs:
			switch {
			case panel.IsOKResponse(res):
				r.handleLastCommandResponse(res)
			case panel.IsErrorResponse(res):
				r.handleLastCommandResponse(res)

			case panel.IsEventResponse(res):
				r.HandleEvents(res)
			}
			r.ServerChan <- res
			r.DatasChan <- res
			r.CommandLineChan <- res
		case <-r.QuitChan:
			r.stop()
			return
		}
	}
}

func (r *Router) stop() {
	close(r.ChatChan)
	close(r.CommandLineChan)
	close(r.ServerChan)
	close(r.NavChan)
	close(r.GroupChan)
	close(r.GroupLeaveChan)
	close(r.UsersChan)
	close(r.ItemsChan)
	close(r.InteractionChan)
	close(r.DatasChan)
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
	case pr.CmdUsers:
		r.UsersChan <- res
	case pr.JoinGroup:
		r.GroupChan <- res
	case pr.CreateGroup:
		r.GroupChan <- res
	case pr.LeaveGroup:
		r.GroupLeaveChan <- res

	default:

	}
}
