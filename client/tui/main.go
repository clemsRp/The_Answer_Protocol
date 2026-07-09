package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

type Datas struct {
	room          string
	status        string
	inventory     []string
	invitation    []string
	group         string
	hp            int
	max_hp        int
	connected     bool
	last_cmd_time time.Time
	spam_warning  int
}

type Client struct {
	conn  net.Conn
	ch    chan Response
	ip    string
	name  string
	datas Datas
}

type Request struct {
	cli *Client
	msg string
}

type Response struct {
	Msg   string
	Datas any
	Req   Request
}

var (
	inputs  = make(chan string)
	outputs = make(chan string)
)

func main() {
	// Connect to server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	inputs := make(chan string)
	outputs := make(chan string)
	router := NewRouter(inputs, outputs)
	app := NewMyApp(router)

	// Handle input
	go func() {
		for input := range inputs {
			fmt.Fprintf(conn, input+"\n")
		}
	}()

	// Handle output
	go func() {
		decoder := json.NewDecoder(conn)

		for {
			var res Response

			if err := decoder.Decode(&res); err != nil {
				fmt.Print("An error occured during connection:", err)
				app.Stop()
				os.Exit(0)
			}

			outputs <- res.Msg

			if res.Msg == "OK bye" {
				app.Stop()
				conn.Close()
				os.Exit(0)
			}
		}
	}()

	if err := app.Run(); err != nil {
		panic(fmt.Sprintf("Execution error: %v", err))
	}
}
