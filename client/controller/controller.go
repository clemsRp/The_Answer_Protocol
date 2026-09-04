package controller

import (
	"context"
	"strings"
	"sync"
	"tap/client/network"
	"tap/client/state"
	panel "tap/client/tui/panels"
	"tap/protocol"
	pr "tap/protocol"
)

type UIApp interface {
	QueueUpdate(f func())
	ShowConnectPage()
	ShowGamePage()
	ShowCombatPage()
	ShowPopupPage()
	ClosePopup()
	UpdateNavigation(room *protocol.LookCommandData)
	UpdateItems(roomItems, inventory []string)
	UpdateInteraction(npcs, players []string, npcData map[string]protocol.InspectNPCData, npcDialogues map[string]string)
	UpdateQuests(quests []protocol.TrackedQuestData)
	UpdateGroup(groupState state.GroupState)
	UpdateCombat(combatState state.CombatState)
	UpdateDatas(text string)
	UpdateInspector(text string)
	AppendChat(scope, user, msg string)
	AppendCombatChat(user, msg string)
	AppendServerResponse(res protocol.ServerResponse)
	AppendCliMessage(text string)
	AppendCliResponse(res protocol.ServerResponse)
	GetPseudo() string
	SetPseudo(pseudo string)
	Stop()
}

type Controller struct {
	gameState   *state.GameState
	netCli      *network.Client
	ui          UIApp
	actions     chan panel.Action
	ctx         context.Context
	cancel      context.CancelFunc
	lastCommands []string
	cmdMu        sync.Mutex

	npcCache   map[string]protocol.InspectNPCData
	npcCacheMu sync.RWMutex
}

func New(ctx context.Context, cancel context.CancelFunc, netCli *network.Client, gameState *state.GameState, ui UIApp, actions chan panel.Action) *Controller {
	c := &Controller{
		gameState:    gameState,
		netCli:       netCli,
		ui:           ui,
		actions:      actions,
		ctx:          ctx,
		cancel:       cancel,
		lastCommands: make([]string, 0),
	}

	go c.eventLoop()

	return c
}

func (c *Controller) eventLoop() {
	for {
		select {
		case action := <-c.actions:
			c.handleUIAction(action)

		case res, ok := <-c.netCli.Outputs():
			if !ok {
				return
			}
			c.handleServerResponses(res)

		case <-c.ctx.Done():
			c.ui.Stop()
			return
		}
	}
}

func (c *Controller) handleUIAction(action panel.Action) {
	switch action.Type {
	case panel.ActionQuit:
		c.sendToNetwork(pr.CmdQuit)
		c.ui.Stop()
		c.cancel()

	case panel.ActionSendServer:
		if payload, ok := action.Payload.(string); ok {
			payload = strings.TrimSpace(payload)
			if payload != "" {
				if strings.ToUpper(payload) == pr.CmdQuit {
					c.sendToNetwork(pr.CmdQuit)
					c.ui.Stop()
					c.cancel()
					return
				}
				c.sendToNetwork(payload)
			}
		}

	case panel.ActionNavigate:
		if payload, ok := action.Payload.(string); ok {
			c.ui.QueueUpdate(func() {
				switch payload {
				case "Game":
					c.ui.ShowGamePage()
				case "Connexion":
					c.ui.ShowConnectPage()
				case "Combat":
					c.ui.ShowCombatPage()
				}
			})
		}

	case panel.ActionOpenPopUp:
		c.ui.QueueUpdate(func() {
			c.ui.ShowPopupPage()
		})

	case panel.ActionClosePopUp:
		c.ui.QueueUpdate(func() {
			c.ui.ClosePopup()
		})

	case panel.ActionOpenCombat:
		c.ui.QueueUpdate(func() {
			c.ui.ShowCombatPage()
		})

	case panel.ActionCloseCombat:
		c.ui.QueueUpdate(func() {
			c.ui.ShowGamePage()
		})
	}
}

func (c *Controller) sendToNetwork(cmd string) {
	c.setLastCommand(cmd)
	c.netCli.Send(cmd)
}

func (c *Controller) setLastCommand(cmd string) {
	c.cmdMu.Lock()
	defer c.cmdMu.Unlock()
	c.lastCommands = append(c.lastCommands, cmd)
}

func (c *Controller) getLastCommand() string {
	c.cmdMu.Lock()
	defer c.cmdMu.Unlock()
	if len(c.lastCommands) > 0 {
		cmd := c.lastCommands[0]
		c.lastCommands = c.lastCommands[1:]
		return cmd
	}
	return ""
}

func (c *Controller) handleServerResponses(res pr.ServerResponse) {
	if strings.HasPrefix(res.Msg, pr.MsgEvt) {
		c.handleEvents(res)
		return
	}
	
	if res.Msg == "OK hello proto=1" {
		return
	}

	c.handleCommandResponses(res)
}

func (c *Controller) refreshUI() {
	playerSnap := c.gameState.GetPlayerSnapshot()
	var roomCopy protocol.LookCommandData
	c.gameState.Read(func(s *state.GameState) {
		if s.Player != nil && s.Player.Room != nil {
			roomCopy = *s.Player.Room
		}
	})

	npcData := c.getNpcCache()

	c.ui.QueueUpdate(func() {
		c.ui.UpdateNavigation(&roomCopy)
		c.ui.UpdateItems(roomCopy.Items, playerSnap.Inventory)

		filteredPlayers := make([]string, 0)
		for _, p := range roomCopy.Players {
			if p != playerSnap.Name {
				filteredPlayers = append(filteredPlayers, p)
			}
		}
		c.ui.UpdateInteraction(roomCopy.Npcs, filteredPlayers, npcData, playerSnap.NpcDialogues)
	})
}

// setNpcCache replaces the whole cache of inspected npc datas (e.g. after a
// room-wide INSPECT).
func (c *Controller) setNpcCache(data map[string]protocol.InspectNPCData) {
	c.npcCacheMu.Lock()
	defer c.npcCacheMu.Unlock()
	c.npcCache = data
}

// cacheNpc stores/updates a single npc's inspect data (e.g. after
// INSPECT NPC <name>).
func (c *Controller) cacheNpc(n protocol.InspectNPCData) {
	c.npcCacheMu.Lock()
	defer c.npcCacheMu.Unlock()
	if c.npcCache == nil {
		c.npcCache = make(map[string]protocol.InspectNPCData)
	}
	c.npcCache[n.Id] = n
}

func (c *Controller) getNpcCache() map[string]protocol.InspectNPCData {
	c.npcCacheMu.RLock()
	defer c.npcCacheMu.RUnlock()
	cp := make(map[string]protocol.InspectNPCData, len(c.npcCache))
	for k, v := range c.npcCache {
		cp[k] = v
	}
	return cp
}

func (c *Controller) refreshGroupUI() {
	groupSnap := c.gameState.GetGroupSnapshot()
	c.ui.QueueUpdate(func() {
		c.ui.UpdateGroup(groupSnap)
	})
}

func filterUngrouped(ungrouped []string, members []string) []string {
	var filtered []string
	for _, user := range ungrouped {
		isMember := false
		for _, member := range members {
			if user == member {
				isMember = true
				break
			}
		}
		if !isMember {
			filtered = append(filtered, user)
		}
	}
	return filtered
}
