package tui

import (
	"bufio"
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
	quitChan      chan struct{}
	conn          net.Conn
	wg            sync.WaitGroup
	disconnectMsg string
	stopOnce      sync.Once
}

func NewTuiClient(conn net.Conn) *TuiClient {
	// Init Client
	cli := TuiClient{
		inputs:   make(chan string, 10),
		outputs:  make(chan pr.ServerResponse, 100),
		quitChan: make(chan struct{}),
		conn:     conn,
	}

	router := NewRouter(cli.inputs, cli.outputs, cli.quitChan)
	cli.app = NewMyApp(router)

	return &cli
}

func (tui *TuiClient) closeQuitChannelSafely() {
	select {
	case <-tui.quitChan:
	default:
		close(tui.quitChan)
	}
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
		case <-tui.quitChan:
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
	tui.closeQuitChannelSafely()
	tui.conn.Close()
	tui.wg.Wait()

	if err != nil {
		fmt.Printf("TUI encountered fatal error: %v\n", err)
	}
	if tui.disconnectMsg != "" {
		fmt.Println(tui.disconnectMsg)
	}

}
