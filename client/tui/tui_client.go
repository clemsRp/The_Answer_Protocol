package tui

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	pr "tap/protocol"
)

type TuiClient struct {
	app           *MyApp
	inputs        chan string
	outputs       chan pr.ServerResponse
	conn          net.Conn
	wg            sync.WaitGroup
	disconnectMsg string
	stopOnce      sync.Once
	ctx           context.Context
	cancelFunc    context.CancelFunc
}

func NewTuiClient(conn net.Conn) *TuiClient {
	// Init Client
	ctx, cancel := context.WithCancel(context.Background())
	cli := TuiClient{
		inputs:     make(chan string, 10),
		outputs:    make(chan pr.ServerResponse, 100),
		conn:       conn,
		ctx:        ctx,
		cancelFunc: cancel,
	}

	router := NewRouter(cli.ctx, cli.inputs, cli.outputs)
	cli.app = NewMyApp(cli.ctx, &cli.wg, router)

	return &cli
}

func (tui *TuiClient) startRouter() {
	defer tui.wg.Done()
	tui.app.router.Start()
}
func (tui *TuiClient) handleInput() {
	defer tui.wg.Done()

	router := tui.app.router
	for {
		select {
		case input := <-tui.inputs:
			// Send command to the server
			fmt.Fprint(tui.conn, input+"\n")
			// Save the last command to handle server returns
			router.LastCommand = strings.ToUpper(strings.Split(input, " ")[0])

			if strings.ToUpper(router.LastCommand) == "GROUP" {
				router.LastCommand = strings.ToUpper(strings.Split(input, " ")[1])
			}
		case <-tui.ctx.Done():
			return
		}
	}
}

func (tui *TuiClient) listenResponses() {
	defer tui.wg.Done()
	scanner := bufio.NewScanner(tui.conn)
	tui.disconnectMsg = "Closed the app and the connection gracefully."

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		res := convertServerResponse(line)
		tui.outputs <- res
		if res.Msg == "OK bye" {
			break
		}
		if res.Msg == "OK connected" {
			tui.app.ShowGamePage()
		}
	}
	if err := scanner.Err(); err != nil {
		tui.disconnectMsg = fmt.Sprintf("Server connection lost with error: %v", err)
	}

	tui.Stop()
}

func (tui *TuiClient) Stop() {
	tui.stopOnce.Do(func() {
		tui.app.Stop()
		tui.cancelFunc()
		tui.conn.Close()
	})

}

func (tui *TuiClient) Start() {
	// handle router
	tui.wg.Add(1)
	go tui.startRouter()

	// Handle input
	tui.wg.Add(1)
	go tui.handleInput()

	// handle server responses
	tui.wg.Add(1)
	go tui.listenResponses()

	err := tui.app.Run()

	tui.Stop()

	tui.wg.Wait()
	if err != nil {
		fmt.Printf("TUI encountered fatal error: %v\n", err)
	}
	if tui.disconnectMsg != "" {
		fmt.Println(tui.disconnectMsg)
	}

}
