package main

import (
	"encoding/json"
	"fmt"
	"strings"
	pr "tap/protocol"
)

type Router struct {
	Inputs  chan<- string
	Outputs <-chan pr.Response

	ChatChan        chan string
	CommandLineChan chan string
	ServerChan      chan string
	NavChan         chan string
	PlayersChan     chan string
	DialogueChan    chan string
	DatasChan       chan string
	LastCommand     string
}

func NewRouter(inputs chan string, outputs <-chan pr.Response) *Router {
	return &Router{
		Inputs:          inputs,
		Outputs:         outputs,
		ChatChan:        make(chan string),
		CommandLineChan: make(chan string),
		ServerChan:      make(chan string),
		NavChan:         make(chan string),
		PlayersChan:     make(chan string),
		DialogueChan:    make(chan string),
		DatasChan:       make(chan string),
		LastCommand:     "",
	}
}

func (r *Router) HandleEvents(res pr.Response) {
	// Handle CHAT responses
	global := strings.HasPrefix(res.Msg, "EVT GLOBAL CHAT")
	room := strings.HasPrefix(res.Msg, "EVT ROOM CHAT")
	group := strings.HasPrefix(res.Msg, "EVT GROUP CHAT alice ca va ?")
	if global || room || group {
		split_msg := strings.SplitN(res.Msg, " ", 5)
		r.ChatChan <- fmt.Sprintf("%s %s %s", split_msg[1], split_msg[3], split_msg[4])
	}
}

func (r *Router) Start(m *MyApp) {
	go func() {

		for res := range r.Outputs {
			color := "white"
			switch {
			case strings.HasPrefix(res.Msg, "OK"):
				color = "green"

			// TO DO
			case strings.HasPrefix(res.Msg, "ERR"):
				color = "red"

			case strings.HasPrefix(res.Msg, "EVT"):
				color = "orange"
				r.HandleEvents(res)
			}

			msg_type := res.Msg
			msg_text := ""

			if strings.ContainsRune(res.Msg, ' ') {
				split_msg := strings.SplitN(res.Msg, " ", 2)

				msg_type = split_msg[0]
				msg_text = split_msg[1]
			}

			r.ServerChan <- fmt.Sprintf("[%s]%s [white]%s", color, msg_type, msg_text)

			datas_json, _ := json.Marshal(res.Datas)
			datas := string(datas_json)

			if datas == "\"\"" || datas == "null" || datas == "nil" {
				datas = ""
			}
			r.CommandLineChan <- res.Msg + datas
		}
	}()
}
