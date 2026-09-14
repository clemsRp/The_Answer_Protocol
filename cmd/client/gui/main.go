package main

import (
	"fmt"
	"net"
	"os"
	gui "tap/client/gui/src"
)

func main() {

	// Connect to server

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Connection error:", err)
		os.Exit(1)
	}
	// Initialize client
	gui_app := gui.NewGuiClient(conn)
	gui_app.Start()
	defer gui_app.Stop()
}
