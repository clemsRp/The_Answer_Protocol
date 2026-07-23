package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	pr "tap/protocol"
)

var (
	inputs  = make(chan string)
	outputs = make(chan pr.Response)
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
	router.Start(app)

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
		decoder := json.NewDecoder(conn)

		for {
			var res pr.Response

			if err := decoder.Decode(&res); err != nil {
				fmt.Print("An error occured during connection:", err)
				app.Stop()
				os.Exit(0)
			}

			outputs <- res

			// Handle Login/Logout
			if res.Msg == "OK bye" {
				app.Stop()
				conn.Close()
				os.Exit(0)
			}
			if res.Msg == "OK connected" {
				app.ShowGamePage()
				app.Draw()
			}
		}
	}()

	if err := app.Run(); err != nil {
		panic(fmt.Sprintf("Execution error: %v", err))
	}
}
