package main

import (
	"bufio"
	"encoding/json"
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

			// Handle Login/Logout
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

func convertServerResponse(line string) pr.ServerResponse {
	if strings.ContainsRune(line, '{') {
		split_line := strings.SplitN(line, " {", 2)

		msg := split_line[0]
		datas := "{" + split_line[1]

		var data_json map[string]interface{}
		err := json.Unmarshal([]byte(datas), &data_json)
		if err != nil {
			fmt.Println("Erreur de décodage :", err)
			return pr.ServerResponse{}
		}

		return pr.ServerResponse{Msg: msg, Datas: data_json}

	} else {
		return pr.ServerResponse{Msg: line}
	}
}
