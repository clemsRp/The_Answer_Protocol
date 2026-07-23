package game

import (
	"fmt"
	"strings"
	"tap/internal/network"
)

type Engine struct {
	InChan       <-chan network.IncomingEvent
	OutChan      chan<- network.OutgoingEvent
	stopGameChan <-chan struct{}
	players      map[string]*Player
}

func NewEngine(in <-chan network.IncomingEvent, out chan<- network.OutgoingEvent, stop_game_chan <-chan struct{}) *Engine {
	return &Engine{
		InChan:       in,
		OutChan:      out,
		stopGameChan: stop_game_chan,
		players:      make(map[string]*Player),
	}
}

func (e *Engine) Run() {
	fmt.Println("[Game] Engine started.")
	for {
		select {
		case event := <-e.InChan:
			cmd := strings.TrimSpace(event.Payload)
			fmt.Printf("[Game] command received from %s : %s\n", event.ClientID, cmd)

			reponse := fmt.Sprintf("TEST FROM ASSAULT PARTY : %s\n", cmd)

			e.OutChan <- network.OutgoingEvent{
				ClientID: event.ClientID,
				Payload:  reponse,
			}

		case <-e.stopGameChan:
			fmt.Println("[Game] Engine stopped.")
			return
		}
	}
}
