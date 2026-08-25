package tui

import (
	"context"
	"slices"
	"strings"
	panel "tap/client/tui/panels"
	pr "tap/protocol"
)

type Router struct {
	app            *MyApp
	Inputs         chan<- string
	Outputs        <-chan pr.ServerResponse
	pendingCmds    chan string
	GameUpdateChan chan pr.ServerResponse

	ChatChan        chan pr.ServerResponse
	CommandLineChan chan pr.ServerResponse
	ServerChan      chan pr.ServerResponse
	NavChan         chan pr.ServerResponse
	GroupChan       chan pr.ServerResponse
	ItemsChan       chan pr.ServerResponse
	InteractionChan chan pr.ServerResponse
	DatasChan       chan pr.ServerResponse
	QuitChan        <-chan struct{}
	LastCommand     string
	ctx             context.Context
}

func NewRouter(ctx context.Context, inputs chan string, outputs <-chan pr.ServerResponse) *Router {

	return &Router{
		Inputs:         inputs,
		Outputs:        outputs,
		pendingCmds:    make(chan string, 50),
		GameUpdateChan: make(chan pr.ServerResponse, 100),

		ChatChan:        make(chan pr.ServerResponse, 100),
		CommandLineChan: make(chan pr.ServerResponse, 100),
		ServerChan:      make(chan pr.ServerResponse, 100),
		NavChan:         make(chan pr.ServerResponse, 100),
		GroupChan:       make(chan pr.ServerResponse, 100),
		ItemsChan:       make(chan pr.ServerResponse, 100),
		InteractionChan: make(chan pr.ServerResponse, 100),
		DatasChan:       make(chan pr.ServerResponse, 100),
		LastCommand:     "",
		ctx:             ctx,
	}
}

func (r *Router) Start(app *MyApp) {
	r.app = app
	go func() {
		for res := range r.GameUpdateChan {
			r.app.app.QueueUpdateDraw(func() {
				r.app.NavListenOutputs(res)
				r.app.ItemListenOutputs(res)
				r.app.InteractionListenOutputs(res)

				r.app.app.Sync()
			})
		}
	}()

	for {
		select {
		case res := <-r.Outputs:
			switch {
			case panel.IsOKResponse(res):
				r.handleLastCommandResponse(res)
			case panel.IsErrorResponse(res):
				r.handleLastCommandResponse(res)

			case panel.IsEventResponse(res):
				r.HandleEvents(res)
			}
			r.ServerChan <- res
			r.DatasChan <- res
			r.CommandLineChan <- res
		case <-r.ctx.Done():
			r.Stop()
			return
		}
	}
}

func (r *Router) Stop() {
	close(r.ChatChan)
	close(r.CommandLineChan)
	close(r.ServerChan)
	close(r.NavChan)
	close(r.GroupChan)
	close(r.ItemsChan)
	close(r.InteractionChan)
	close(r.DatasChan)
}

func (m *MyApp) UpdateRouterLayout() {
	if m.grid != nil {
		m.grid.Clear()
	}

	if m.Navigation != nil {
		m.grid.AddItem(m.Navigation.Layout, 0, 0, 1, 1, 0, 0, true)
	}
}

func (r *Router) handleLastCommandResponse(res pr.ServerResponse) {
	// Get last command
	var lastCmd string
	select {
	case lastCmd = <-r.pendingCmds:
	default:
		lastCmd = ""
	}
	r.LastCommand = lastCmd

	// Handle Group commands
	group_commands := []string{
		pr.CmdUnGrouped,
		pr.CmdGrouped,
		pr.PromoteGroup,
		pr.AcceptPromoteGroup,
		pr.DeclinePromoteGroup,
		pr.JoinGroup,
		pr.CreateGroup,
		pr.LeaveGroup,
		pr.KickGroup,
	}
	if slices.Contains(group_commands, lastCmd) {
		r.GroupChan <- res

		if lastCmd == pr.CmdGrouped || lastCmd == pr.CmdUnGrouped {
			r.DatasChan <- res
		}

	} else {
		switch lastCmd {
		case pr.CmdConnect:
			go func() {
				r.Inputs <- pr.CmdUnGrouped
				r.Inputs <- pr.CmdLook
			}()
		case pr.CmdLook:
			r.GameUpdateChan <- res
		case pr.CmdChat:
			r.ChatChan <- res
		default:

		}
	}
}

func (r *Router) HandleEvents(res pr.ServerResponse) {
	// Update INVITE datas
	if strings.HasPrefix(res.Msg, pr.EventPrefixStats) || strings.HasPrefix(res.Msg, pr.EventPrefixRoomPresence) {
		r.Inputs <- pr.CmdUnGrouped
		r.Inputs <- pr.CmdLook
	}

	// Handle CHAT responses
	global := strings.HasPrefix(res.Msg, pr.EventGlobalChat)
	room := strings.HasPrefix(res.Msg, pr.EventRoomChat)
	group := strings.HasPrefix(res.Msg, pr.EventGroupChat)
	if global || room || group {
		r.ChatChan <- res
	}

	// Handle ITEM responses
	if strings.HasPrefix(res.Msg, pr.EventPrefixItem) {
		r.ItemsChan <- res
	}

	// Handle PRESENCE responses
	if strings.HasPrefix(res.Msg, pr.EventPrefixRoomPresence) {
		r.InteractionChan <- res
	}

	// Handle GROUP responses
	if strings.HasPrefix(res.Msg, pr.EventPrefixGroup) || strings.HasPrefix(res.Msg, pr.EventGroupNewLeader) {
		r.GroupChan <- res
	}
}
