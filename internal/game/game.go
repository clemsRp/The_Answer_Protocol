package game

import (
	"fmt"
	"strings"
	"sync"
	"tap/internal/network"
)

type Engine struct {
	InChan       <-chan network.IncomingEvent
	OutChan      chan<- network.OutgoingEvent
	stopGameChan chan struct{}
	players      map[string]*Player
	stopOnce     sync.Once
	wg           sync.WaitGroup
}

func NewEngine(in <-chan network.IncomingEvent, out chan<- network.OutgoingEvent) *Engine {
	return &Engine{
		InChan:       in,
		OutChan:      out,
		stopGameChan: make(chan struct{}),
		players:      make(map[string]*Player),
	}
}
func (e *Engine) Start() {
	e.wg.Add(1)
	go e.runLoop()
}

func (e *Engine) runLoop() {
	defer e.wg.Done()
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
			return
		}
	}
}

// func (e *Engine) processEvent(request string) {

// }

func (e *Engine) Stop() {
	e.stopOnce.Do(func() {
		close(e.stopGameChan)
	})
	e.wg.Wait()
	fmt.Println("[Game] Engine stopped.")
}
