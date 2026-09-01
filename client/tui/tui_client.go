package tui

import (
	"context"
	"fmt"
	"net"
	"sync"
	"tap/client/controller"
	"tap/client/network"
	"tap/client/state"
	panel "tap/client/tui/panels"
)

type TuiClient struct {
	app         *MyApp
	netCli      *network.Client
	controller  *controller.Controller
	gameState   *state.GameState
	actionsChan chan panel.Action
	wg          sync.WaitGroup
	ctx         context.Context
	cancelFunc  context.CancelFunc
}

func NewTuiClient(conn net.Conn) *TuiClient {
	ctx, cancel := context.WithCancel(context.Background())
	actionsChan := make(chan panel.Action, 100)

	netCli := network.NewClient(ctx, cancel, conn)
	gameState := state.New()

	var wg sync.WaitGroup
	myApp := NewMyApp(ctx, &wg, actionsChan)

	ctrl := controller.New(ctx, cancel, netCli, gameState, myApp, actionsChan)

	return &TuiClient{
		app:         myApp,
		netCli:      netCli,
		controller:  ctrl,
		gameState:   gameState,
		actionsChan: actionsChan,
		ctx:         ctx,
		cancelFunc:  cancel,
	}
}

func (c *TuiClient) Start() {
	c.netCli.Start()

	if err := c.app.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
	}

	c.netCli.Stop()
	fmt.Println(c.netCli.DisconnectMsg)
}

func (c *TuiClient) Stop() {
	c.cancelFunc()
	c.app.Stop()
	c.netCli.Stop()
}
