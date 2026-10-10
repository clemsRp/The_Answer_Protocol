package main

import (
	"fmt"
	"net"
	"os"
	"tap/src/client/tui"
	"tap/src/cmd"
)

func main() {
	// Get args
	ip, port, conn_err := cmd.GetConnection(true)
	if conn_err != nil {
		fmt.Println("Error: ", conn_err.Error())
		return
	}

	// Connect to server
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Connection error:", err)
		os.Exit(1)
	}
	// Initialize client
	cli := tui.NewTuiClient(conn)

	cli.Start()
}
