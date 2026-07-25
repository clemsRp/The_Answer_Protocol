package main

import (
	"fmt"
	"strings"
	pr "tap/protocol"
)

type Router struct {
	Inputs  chan<- string
	Outputs <-chan pr.Response

	ChatChan      chan string
	ServerChan    chan string
	NavChan       chan string
	PlayersChan   chan string
	DialogueChan  chan string
	InventoryChan chan string
	QuestChan     chan string
	ItemsChan     chan string
	LastCommand   string
}

func NewRouter(inputs chan string, outputs <-chan pr.Response) *Router {
	return &Router{
		Inputs:        inputs,
		Outputs:       outputs,
		ChatChan:      make(chan string),
		ServerChan:    make(chan string),
		NavChan:       make(chan string),
		PlayersChan:   make(chan string),
		DialogueChan:  make(chan string),
		InventoryChan: make(chan string),
		QuestChan:     make(chan string),
		ItemsChan:     make(chan string),
		LastCommand:   "",
	}
}

func (r *Router) HandleEvents(res pr.Response) {
	// Handle CHAT responses
	global := strings.HasPrefix(res.Msg, "EVT GLOBAL CHAT")
	room := strings.HasPrefix(res.Msg, "EVT ROOM CHAT")
	group := strings.HasPrefix(res.Msg, "EVT GROUP CHAT")
	if global || room || group {
		split_msg := strings.SplitN(res.Msg, " ", 5)
		r.ChatChan <- fmt.Sprintf("%s %s %s", split_msg[1], split_msg[3], split_msg[4])
	}
}

func (r *Router) Start(m *MyApp) {
	go func() {

		for res := range r.Outputs {
			switch {
			// case strings.HasPrefix(res.Msg, "OK"):

			// TO DO
			// case strings.HasPrefix(res.Msg, "ERR"):

			case strings.HasPrefix(res.Msg, "EVT"):
				r.HandleEvents(res)

			default:
				r.ServerChan <- "Unknown format: " + res.Msg
			}
		}
	}()
}
