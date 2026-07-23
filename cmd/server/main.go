package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tap/internal/game"
	network "tap/internal/network"
)

func main() {
	inChan := make(chan network.IncomingEvent, 100)
	outChan := make(chan network.OutgoingEvent, 100)
	stopGameChan := make(chan struct{})

	server := network.NewServer(":8080", inChan, outChan)
	engine := game.NewEngine(inChan, outChan, stopGameChan)
	// to wait for an error at the launch of the server same way as sigch.

	errch := make(chan error, 1)
	go server.SafeStart(errch)

	// The only way to stop the server is by CTRL+C or KILL command
	// then it stops gracefully.
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-sigch:
		fmt.Printf("\nStop signal received (%v), closing the server...", sig)
	case err := <-errch:
		fmt.Printf("\nCritical error from server: %v", err)
	}

	server.Stop()
	fmt.Println("\nServer stopped graciously.")
}
