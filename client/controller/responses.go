package controller

import (
	"encoding/json"
	"strings"
	"tap/client/state"
	"tap/protocol"
	pr "tap/protocol"
)

func (c *Controller) handleCommandResponses(res pr.ServerResponse) {
	// Update CLI/Server panel
	c.ui.QueueUpdate(func() {
		c.ui.AppendServerResponse(res)
		c.ui.AppendCliResponse(res)
	})

	lastCmd := c.getLastCommand()
	lastCmdBase := ""
	cmdFields := strings.Fields(lastCmd)
	if len(cmdFields) > 0 {
		lastCmdBase = strings.ToUpper(cmdFields[0])
	}

	if strings.HasPrefix(res.Msg, pr.MsgErr) {
		return
	}

	switch {
	case lastCmdBase == pr.CmdConnect && (res.Msg == pr.MsgOK || strings.HasPrefix(res.Msg, pr.MsgOK)):
		if len(cmdFields) >= 2 {
			pseudo := cmdFields[1]
			c.gameState.UpdatePlayer(func(p *state.Player) {
				p.Name = pseudo
			})
			c.ui.SetPseudo(pseudo)
		}
		c.ui.QueueUpdate(func() {
			c.ui.ShowGamePage()
		})
		c.sendToNetwork(pr.CmdLook)

	case lastCmdBase == pr.CmdLook && res.Datas != nil:
		var lookData protocol.LookCommandData
		raw, err := json.Marshal(res.Datas)
		if err == nil && json.Unmarshal(raw, &lookData) == nil {
			c.gameState.UpdateRoomLook(&lookData)
			c.refreshUI()
		}

	case lastCmdBase == pr.CmdInventory && res.Datas != nil:
		var inventoryData []string
		raw, err := json.Marshal(res.Datas)
		if err == nil && json.Unmarshal(raw, &inventoryData) == nil {
			c.gameState.UpdatePlayer(func(p *state.Player) {
				p.Inventory = append([]string{}, inventoryData...)
			})
			c.refreshUI()
		}

	case lastCmd == pr.CmdCombatStats && res.Datas != nil:
		var combatData protocol.CombatStatsCommandData
		raw, err := json.Marshal(res.Datas)
		if err == nil && json.Unmarshal(raw, &combatData) == nil {
			c.gameState.UpdateCombatState(func(cs *state.CombatState) {
				cs.InCombat = true
				cs.CurrentTurn = combatData.CurrentTurn
				cs.Leader = combatData.Leader
				cs.Team = combatData.Team
				cs.Opponents = combatData.Opponents
			})
			combatSnap := c.gameState.GetCombatSnapshot()
			c.ui.QueueUpdate(func() {
				c.ui.UpdateCombat(combatSnap)
				c.ui.ShowCombatPage()
			})
		}

	case (lastCmdBase == pr.CmdAttack || strings.HasPrefix(lastCmd, pr.CmdAttack) || strings.HasPrefix(lastCmd, pr.CmdChatCombatAttack)) && (res.Msg == pr.MsgOK || strings.HasPrefix(res.Msg, pr.MsgOK)):
		if res.Datas != nil {
			var fullTurn struct {
				CombatState string `json:"combat_state"`
			}
			raw, err := json.Marshal(res.Datas)
			if err == nil && json.Unmarshal(raw, &fullTurn) == nil {
				if fullTurn.CombatState == "VICTORY" || fullTurn.CombatState == "DEFEAT" {
					c.gameState.UpdateCombatState(func(cs *state.CombatState) {
						cs.InCombat = false
					})
					c.ui.QueueUpdate(func() {
						c.ui.ShowGamePage()
					})
					c.sendToNetwork(pr.CmdLook)
					break
				}
			}
		}
		c.sendToNetwork(pr.CmdCombatStats)

	case (lastCmdBase == pr.CmdFlee || strings.HasPrefix(lastCmd, pr.CmdFlee) || strings.HasPrefix(lastCmd, pr.CmdChatCombatFlee)) && (res.Msg == pr.MsgOK || strings.HasPrefix(res.Msg, pr.MsgOK)):
		c.gameState.UpdateCombatState(func(cs *state.CombatState) {
			cs.InCombat = false
		})
		c.ui.QueueUpdate(func() {
			c.ui.ShowGamePage()
		})
		c.sendToNetwork(pr.CmdLook)

	case strings.HasPrefix(lastCmd, pr.CmdChatCombat) && res.Msg == pr.MsgOK:
		chatMsg := strings.TrimSpace(strings.TrimPrefix(lastCmd, pr.CmdChatCombat))
		pseudo := c.ui.GetPseudo()
		c.gameState.UpdateCombatState(func(cs *state.CombatState) {
			cs.LastCombatChat = chatMsg
			cs.Chats = append(cs.Chats, state.CombatChat{Pseudo: pseudo, Msg: chatMsg})
		})
		c.ui.QueueUpdate(func() {
			c.ui.AppendCombatChat(pseudo, chatMsg)
		})

	case strings.HasPrefix(lastCmd, pr.CmdChat+" ") && res.Msg == pr.MsgOK:
		parts := strings.SplitN(lastCmd, " ", 3)
		if len(parts) >= 2 {
			scope := parts[1]
			chatMsg := ""
			if len(parts) >= 3 {
				chatMsg = parts[2]
			}
			pseudo := c.ui.GetPseudo()
			c.ui.QueueUpdate(func() {
				c.ui.AppendChat(scope, pseudo, chatMsg)
			})
		}

	case strings.HasPrefix(lastCmd, pr.CmdTake) && strings.HasPrefix(res.Msg, pr.MsgOK):
		item := strings.SplitN(res.Msg, "taken=", 2)[1]
		c.gameState.UpdatePlayer(func(p *state.Player) {
			p.Inventory = append(p.Inventory, item)
		})
		c.gameState.UpdateRoom(func(r *protocol.LookCommandData) {
			for i, it := range r.Items {
				if it == item {
					r.Items = append(r.Items[:i], r.Items[i+1:]...)
					break
				}
			}
		})

		c.refreshUI()

	case strings.HasPrefix(lastCmd, pr.CmdDrop) && strings.HasPrefix(res.Msg, pr.MsgOK):
		item := strings.SplitN(res.Msg, "dropped=", 2)[1]
		c.gameState.UpdatePlayer(func(p *state.Player) {
			for i, it := range p.Inventory {
				if it == item {
					p.Inventory = append(p.Inventory[:i], p.Inventory[i+1:]...)
					break
				}
			}
		})
		c.gameState.UpdateRoom(func(r *protocol.LookCommandData) {
			r.Items = append(r.Items, item)
		})

		c.refreshUI()

	case strings.HasPrefix(res.Msg, pr.PrefixOKGroup):
		parts := strings.SplitN(res.Msg, "group=", 2)
		if len(parts) == 2 {
			groupName := parts[1]
			c.gameState.UpdateGroupState(func(gs *state.GroupState) {
				gs.Group = groupName
				gs.Invitations = make([]string, 0)
				if lastCmd == pr.CreateGroup || strings.HasPrefix(lastCmd, pr.CmdCreateGroup) {
					gs.Leader = true
				}
			})
			c.sendToNetwork(pr.CmdGrouped)
		}

	case strings.HasPrefix(res.Msg, pr.PrefixOKPendingLeader):
		parts := strings.SplitN(res.Msg, "pending_leader=", 2)
		if len(parts) == 2 {
			pendingLeader := parts[1]
			playerSnap := c.gameState.GetPlayerSnapshot()
			c.gameState.UpdateGroupState(func(gs *state.GroupState) {
				gs.Leader = (pendingLeader == playerSnap.Name)
				if lastCmd == pr.PromoteGroup || strings.HasPrefix(lastCmd, pr.CmdPromoteGroup) {
					gs.SendPromotion = true
				}
			})
			c.sendToNetwork(pr.CmdGrouped)
		}

	case strings.HasPrefix(res.Msg, pr.PrefixOKNewLeader):
		c.gameState.UpdateGroupState(func(gs *state.GroupState) {
			targetCopy := ""
			gs.LastKick = &targetCopy
			gs.Leader = true
			gs.Promotion = false
			gs.SendPromotion = false
		})
		c.sendToNetwork(pr.CmdGrouped)

	case res.Msg == pr.MsgOK:
		switch {
		case lastCmd == pr.CmdLeaveGroup || lastCmd == pr.LeaveGroup:
			c.gameState.UpdateGroupState(func(gs *state.GroupState) {
				gs.Group = ""
				gs.Leader = false
				gs.Grouped = make([]string, 0)
				gs.Invitations = make([]string, 0)
			})
			c.refreshGroupUI()

		case lastCmd == pr.CmdDeclinePromoteGroup || lastCmd == pr.DeclinePromoteGroup:
			c.gameState.UpdateGroupState(func(gs *state.GroupState) {
				gs.Promotion = false
			})
			c.refreshGroupUI()

		case lastCmd == pr.CmdKickGroup || strings.HasPrefix(lastCmd, pr.CmdKickGroup):
			c.sendToNetwork(pr.CmdGrouped)

		case lastCmd == pr.CmdUnGrouped:
			var data []string
			raw, err := json.Marshal(res.Datas)
			if err == nil && json.Unmarshal(raw, &data) == nil {
				groupSnap := c.gameState.GetGroupSnapshot()
				filtered := filterUngrouped(data, groupSnap.Grouped)
				c.gameState.UpdateGroupState(func(gs *state.GroupState) {
					gs.UnGrouped = filtered
				})
				c.refreshGroupUI()
			}

		case lastCmd == pr.CmdGrouped:
			var data []string
			raw, err := json.Marshal(res.Datas)
			if err == nil && json.Unmarshal(raw, &data) == nil {
				c.gameState.UpdateGroupState(func(gs *state.GroupState) {
					gs.Grouped = data
				})
				c.refreshGroupUI()
				c.sendToNetwork(pr.CmdUnGrouped)
			}

		case lastCmd == pr.CmdStatus:
			if res.Datas != nil {
				b, _ := json.Marshal(res.Datas)
				c.ui.QueueUpdate(func() {
					c.ui.UpdateDatas(string(b))
				})
			}

		case strings.HasPrefix(lastCmdBase, pr.CmdInspect):
			if res.Datas != nil {
				c.handleInspectResponse(res)
			}
		}
	}
}
