package main

import (
	"fmt"
	"net"
	"os"
	"tap/client/tui"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Connection error:", err)
		os.Exit(1)
	}
	// Initialize client
	cli := tui.NewTuiClient(conn)

	cli.Start()
}
