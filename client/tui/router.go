package main

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
	PlayersChan     chan pr.ServerResponse
	DialogueChan    chan pr.ServerResponse
	DatasChan       chan pr.ServerResponse
	LastCommand     string
}

func NewRouter(inputs chan string, outputs <-chan pr.ServerResponse) *Router {
	return &Router{
		Inputs:          inputs,
		Outputs:         outputs,
		ChatChan:        make(chan pr.ServerResponse),
		CommandLineChan: make(chan pr.ServerResponse),
		ServerChan:      make(chan pr.ServerResponse),
		NavChan:         make(chan pr.ServerResponse),
		PlayersChan:     make(chan pr.ServerResponse),
		DialogueChan:    make(chan pr.ServerResponse),
		DatasChan:       make(chan pr.ServerResponse),
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
		// split_msg := strings.SplitN(res.Msg, " ", 5)
		// fmt.Sprintf("%s %s %s", split_msg[1], split_msg[3], split_msg[4])
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

			// msg_type := res.Msg
			// msg_text := ""

			// if strings.ContainsRune(res.Msg, ' ') {
			// 	split_msg := strings.SplitN(res.Msg, " ", 2)

			// 	msg_type = split_msg[0]
			// 	msg_text = split_msg[1]
			// }

			r.ServerChan <- res
			// fmt.Sprintf("[%s]%s [white]%s", color, msg_type, msg_text)

			// b, _ := json.Marshal(res.Datas)
			// datas := string(b)

			// if datas == "\"\"" || datas == "null" || datas == "<nil>" {
			// 	datas = ""

			// } else {
			// 	r.DatasChan <- datas
			// }
			r.DatasChan <- res
			r.CommandLineChan <- res
			// fmt.Sprintf("%s %+v", res.Msg, datas)
		}
	}()
}

func (r *Router) handleLastCommandResponse(res pr.ServerResponse) {
	switch r.LastCommand {
	case pr.CmdLook:
		r.NavChan <- res
		r.PlayersChan <- res
		// r.itemsChan <- res
	case pr.CmdMove:
		r.NavChan <- res
		r.Inputs <- pr.CmdLook
	case pr.CmdChat:
		r.ChatChan <- res
	default:

	}
}
