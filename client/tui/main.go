package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	pr "tap/protocol"
)

var (
	inputs  = make(chan string, 10)
	outputs = make(chan pr.ServerResponse)
)

func main() {
	// Connect to server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	router := NewRouter(inputs, outputs)
	app := NewMyApp(router)
	router.Start()

	// Handle input
	go func() {
		for input := range inputs {
			// Send command to the server
			fmt.Fprint(conn, input+"\n")
			// Save the last command to handle server returns
			router.LastCommand = strings.Split(input, " ")[0]
		}
	}()

	// Handle output
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			res := convertServerResponse(line)
			outputs <- res
			if res.Msg == "OK bye" {
				app.Stop()
				conn.Close()
				os.Exit(0)
			}
			if res.Msg == "OK connected" {
				app.ShowGamePage()
			}
		}
	}()

	if err := app.Run(); err != nil {
		panic(fmt.Sprintf("Execution error: %v", err))
	}
}
