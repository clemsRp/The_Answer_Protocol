package network

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"tap/protocol"
	pr "tap/protocol"
)

type Client struct {
	inputs        chan string
	outputs       chan pr.ServerResponse
	conn          net.Conn
	wg            sync.WaitGroup
	DisconnectMsg string
	stopOnce      sync.Once
	ctx           context.Context
	cancelFunc    context.CancelFunc
}

func NewClient(ctx context.Context, cancel context.CancelFunc, conn net.Conn) *Client {
	return &Client{
		inputs:     make(chan string, 100),
		outputs:    make(chan pr.ServerResponse, 200),
		conn:       conn,
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func (c *Client) Send(cmd string) {
	select {
	case <-c.ctx.Done():
		return
	case c.inputs <- cmd:
	}
}

func (c *Client) Outputs() <-chan pr.ServerResponse {
	return c.outputs
}

func (c *Client) Start() {
	c.wg.Add(2)
	go c.handleInput()
	go c.listenResponses()
}

func (c *Client) Stop() {
	c.stopOnce.Do(func() {
		c.cancelFunc()
		if c.conn != nil {
			c.conn.Close()
		}
	})
	c.wg.Wait()
}

func (c *Client) handleInput() {
	defer c.wg.Done()
	for {
		select {
		case input := <-c.inputs:
			if c.conn != nil {
				fmt.Fprint(c.conn, input+"\n")
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) listenResponses() {
	defer c.wg.Done()
	scanner := bufio.NewScanner(c.conn)
	c.DisconnectMsg = "Closed the app and the connection gracefully."

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		res := convertServerResponse(line)
		if res.Msg == protocol.ErrSpam {
			c.DisconnectMsg = "Disconnected from server due to continuous SPAM."
		}
		select {
		case <-c.ctx.Done():
			return
		case c.outputs <- res:
		}
	}
	if err := scanner.Err(); err != nil {
		c.DisconnectMsg = fmt.Sprintf("Server connection lost with error: %v", err)
	}
	c.cancelFunc()
}

func findJsonStartIndex(line string) int {
	idx := -1
	for i := 0; i < len(line)-1; i++ {
		if line[i] == ' ' && (line[i+1] == '{' || line[i+1] == '[') {
			idx = i
			break
		}
	}
	return idx
}

func convertServerResponse(line string) pr.ServerResponse {
	line = strings.TrimSpace(line)
	jsonIndex := findJsonStartIndex(line)
	isChatEvent := strings.HasPrefix(line, "EVT ") && strings.Contains(line, " CHAT ")

	if jsonIndex != -1 && !isChatEvent {
		msgPart := strings.TrimSpace(line[:jsonIndex])
		jsonPart := strings.TrimSpace(line[jsonIndex:])

		var rawData any
		err := json.Unmarshal([]byte(jsonPart), &rawData)
		if err == nil {
			return pr.ServerResponse{
				Msg:   msgPart,
				Datas: rawData,
			}
		}
	}

	return pr.ServerResponse{
		Msg:   line,
		Datas: nil,
	}
}
