package tui

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"tap/engine/parser"
	pr "tap/protocol"
)

type TuiClient struct {
	app     *MyApp
	inputs  chan string
	outputs chan pr.ServerResponse
	world   parser.Map
}

func NewTuiClient(world parser.Map) *TuiClient {
	// Init Client
	cli := TuiClient{
		inputs:  make(chan string, 10),
		outputs: make(chan pr.ServerResponse, 100),
		world:   world,
	}

	// Connect to server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Connection error:", err)
		return nil
	}

	router := NewRouter(cli.inputs, cli.outputs)
	cli.app = NewMyApp(router)
	router.Start()

	// Handle input
	go func() {
		for input := range cli.inputs {
			// Send command to the server
			fmt.Fprint(conn, input+"\n")
			// Save the last command to handle server returns
			router.LastCommand = strings.ToUpper(strings.Split(input, " ")[0])

			if strings.ToUpper(router.LastCommand) == "GROUP" {
				router.LastCommand = strings.ToUpper(strings.Split(input, " ")[1])
			}
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
			cli.outputs <- res
			if res.Msg == "OK bye" {
				cli.app.Stop()
				conn.Close()
				os.Exit(0)
			}
			if res.Msg == "OK connected" {
				cli.app.ShowGamePage()
				cli.inputs <- "UNGROUPED"
			}
		}
	}()

	return &cli
}

func (c *TuiClient) Start() {
	if err := c.app.Run(); err != nil {
		panic(fmt.Sprintf("Execution error: %v", err))
	}
}
