package gui

import (
	"context"
	"net"
	"sync"
	"tap/client/controller"
	"tap/client/network"
	"tap/client/state"
	panel "tap/client/tui/panels"
)

type GuiClient struct {
	app         *App
	netCli      *network.Client
	controller  *controller.Controller
	gameState   *state.GameState
	actionsChan chan panel.Action
	wg          sync.WaitGroup
	ctx         context.Context
	cancelFunc  context.CancelFunc
}

func NewGuiClient(conn net.Conn) *GuiClient {
	ctx, cancel := context.WithCancel(context.Background())
	actionsChan := make(chan panel.Action, 100)

	netCli := network.NewClient(ctx, cancel, conn)
	gameState := state.New()

	myApp := NewApp(actionsChan)

	ctrl := controller.New(ctx, cancel, netCli, gameState, myApp, actionsChan)

	return &GuiClient{
		app:         myApp,
		netCli:      netCli,
		controller:  ctrl,
		gameState:   gameState,
		actionsChan: actionsChan,
		ctx:         ctx,
		cancelFunc:  cancel,
	}
}

func (c *GuiClient) Start() {
	c.netCli.Start()
	c.app.Start()
}

func (c *GuiClient) Stop() {
	c.cancelFunc()
	c.app.Stop()
	c.netCli.Stop()
}
