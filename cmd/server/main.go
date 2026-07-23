package main

import (
	"fmt"
	"os/signal"
	"syscall"
	"tap/internal/game"
	network "tap/internal/network"
)

func main() {
	inChan := make(chan network.IncomingEvent, 100)
	outChan := make(chan network.OutgoingEvent, 100)

	server := network.NewServer(":8080", inChan, outChan)
	engine := game.NewEngine(inChan, outChan)
	// to wait for an error at the launch of the server same way as sigch.

	server.Start()
	engine.Start()
	// The only way to stop the server is by CTRL+C or KILL command
	// then it stops gracefully.
	signal.Notify(server.SigChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-server.SigChan:
		fmt.Printf("\nStop signal received (%v), closing the server...", sig)
	case err := <-server.ErrChan:
		fmt.Printf("\nCritical error from server: %v", err)
	}
	server.Stop()
	engine.Stop()

}
