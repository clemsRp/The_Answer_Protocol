package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	network "tap/internal/network"
)

func main() {
	server := network.NewServer(":8080")
	// to wait for CTRL+C signal or KILL signal ASYNCHRONOUSLY with a buffer (not blocking the thread)
	sigch := make(chan os.Signal, 1)
	// to wait for an error at the launch of the server same way as sigch.
	errch := make(chan error, 1)
	go server.GameLoop()

	go server.SafeStart(errch)

	// The only way to stop the server is by CTRL+C or KILL command
	// then it stops gracefully.
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
