package main

import "strings"

type Router struct {
	Inputs  chan<- string
	Outputs <-chan string

	ChatChan      chan string
	ServerChan    chan string
	NavChan       chan string
	PlayersChan   chan string
	DialogueChan  chan string
	InventoryChan chan string
	QuestChan     chan string
	ItemsChan     chan string
}

func NewRouter(inputs chan<- string, outputs <-chan string) *Router {
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
	}
}

func (r *Router) Start() {
	go func() {
		for msg := range r.Outputs {
			switch {
			case strings.HasPrefix(msg, "CHAT:"):
				r.ChatChan <- strings.TrimPrefix(msg, "CHAT:")
			case strings.HasPrefix(msg, "SERVER:"):
				r.ServerChan <- strings.TrimPrefix(msg, "SERVER:")
			case strings.HasPrefix(msg, "NAV:"):
				r.NavChan <- strings.TrimPrefix(msg, "NAV:")
			case strings.HasPrefix(msg, "PLAYERS:"):
				r.PlayersChan <- strings.TrimPrefix(msg, "PLAYERS:")
			case strings.HasPrefix(msg, "DIALOGUE:"):
				r.DialogueChan <- strings.TrimPrefix(msg, "DIALOGUE:")
			case strings.HasPrefix(msg, "INVENTORY:"):
				r.InventoryChan <- strings.TrimPrefix(msg, "INVENTORY:")
			case strings.HasPrefix(msg, "QUEST:"):
				r.QuestChan <- strings.TrimPrefix(msg, "QUEST:")
			default:
				r.ServerChan <- "Unknown format: " + msg
			}
		}
	}()
}
