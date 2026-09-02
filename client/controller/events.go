package controller

import (
	"encoding/json"
	"strings"
	"tap/client/state"
	"tap/protocol"
	pr "tap/protocol"
)

func (c *Controller) handleEvents(res pr.ServerResponse) {
	if !strings.HasPrefix(res.Msg, pr.MsgEvt) {
		return
	}

	trimmed := strings.TrimSpace(strings.TrimPrefix(res.Msg, pr.MsgEvt))
	fields := strings.Fields(trimmed)
	if len(fields) == 0 {
		return
	}

	category := fields[0]

	switch {
	case strings.HasPrefix(trimmed, pr.CategoryRoom+" "+pr.TypePresenceEnter):
		target := strings.TrimSpace(strings.TrimPrefix(trimmed, pr.CategoryRoom+" "+pr.TypePresenceEnter))
		if target != "" {
			c.gameState.UpdateRoom(func(r *protocol.LookCommandData) {
				for _, p := range r.Players {
					if p == target {
						return
					}
				}
				r.Players = append(r.Players, target)
			})
			c.refreshUI()
		}

	case strings.HasPrefix(trimmed, pr.CategoryRoom+" "+pr.TypePresenceLeave):
		target := strings.TrimSpace(strings.TrimPrefix(trimmed, pr.CategoryRoom+" "+pr.TypePresenceLeave))
		if target != "" {
			c.gameState.UpdateRoom(func(r *protocol.LookCommandData) {
				for i, p := range r.Players {
					if p == target && i < len(r.Players) {
						r.Players = append(r.Players[:i], r.Players[i+1:]...)
						break
					}
				}
			})
			c.refreshUI()
		}

	case strings.HasPrefix(trimmed, pr.TypeItemDropped):
		target := strings.SplitN(trimmed, pr.TypeItemDropped+" ", 2)[1]
		if target != "" {
			c.gameState.UpdateRoom(func(r *protocol.LookCommandData) {
				r.Items = append(r.Items, target)
			})
			c.refreshUI()
		}

	case strings.HasPrefix(trimmed, pr.TypeItemTook):
		target := strings.SplitN(trimmed, pr.TypeItemTook+" ", 2)[1]
		if target != "" {
			c.gameState.UpdateRoom(func(r *protocol.LookCommandData) {
				for i, it := range r.Items {
					if it == target {
						r.Items = append(r.Items[:i], r.Items[i+1:]...)
						break
					}
				}
			})
			c.refreshUI()
		}

	case strings.HasPrefix(trimmed, pr.CategoryGroup+" "+pr.TypeInvite):
		groupName := strings.TrimSpace(strings.TrimPrefix(trimmed, pr.CategoryGroup+" "+pr.TypeInvite))
		if groupName != "" {
			c.gameState.UpdateGroupState(func(gs *state.GroupState) {
				for _, inv := range gs.Invitations {
					if inv == groupName {
						return
					}
				}
				gs.Invitations = append(gs.Invitations, groupName)
			})
			c.refreshGroupUI()
		}

	case strings.HasPrefix(trimmed, pr.CategoryGroup+" "+pr.TypeJoin):
		c.sendToNetwork(pr.CmdGrouped)

	case strings.HasPrefix(trimmed, pr.CategoryGroup+" "+pr.TypeKick):
		playerSnap := c.gameState.GetPlayerSnapshot()
		kickedUser := strings.TrimSpace(strings.TrimPrefix(trimmed, pr.CategoryGroup+" "+pr.TypeKick))
		if kickedUser == playerSnap.Name {
			c.gameState.UpdateGroupState(func(gs *state.GroupState) {
				gs.Group = ""
				gs.Leader = false
				gs.Grouped = make([]string, 0)
			})
		} else {
			c.sendToNetwork(pr.CmdGrouped)
		}
		c.refreshGroupUI()

	case strings.HasPrefix(trimmed, pr.CategoryGroup+" "+pr.TypeLeave):
		playerSnap := c.gameState.GetPlayerSnapshot()
		leftUser := strings.TrimSpace(strings.TrimPrefix(trimmed, pr.CategoryGroup+" "+pr.TypeLeave))
		if leftUser == playerSnap.Name {
			c.gameState.UpdateGroupState(func(gs *state.GroupState) {
				gs.Group = ""
				gs.Leader = false
				gs.Grouped = make([]string, 0)
			})
		} else {
			c.sendToNetwork(pr.CmdGrouped)
		}
		c.refreshGroupUI()

	case strings.HasPrefix(trimmed, pr.CategoryGroup+" "+pr.TypeGroupPromoteAccepted):
		c.sendToNetwork(pr.CmdGrouped)

	case strings.HasPrefix(trimmed, pr.CategoryGroup+" "+pr.TypeGroupPromoteDeclined):
		c.gameState.UpdateGroupState(func(gs *state.GroupState) {
			gs.SendPromotion = false
		})
		c.refreshGroupUI()

	case strings.HasPrefix(trimmed, pr.CategoryGroup+" "+pr.TypeGroupPromote):
		c.gameState.UpdateGroupState(func(gs *state.GroupState) {
			gs.Promotion = true
		})
		c.refreshGroupUI()

	case strings.HasPrefix(res.Msg, pr.PrefixEvtNewLeader):
		parts := strings.SplitN(res.Msg, "new_leader=", 2)
		if len(parts) == 2 {
			newLeader := parts[1]
			playerSnap := c.gameState.GetPlayerSnapshot()
			c.gameState.UpdateGroupState(func(gs *state.GroupState) {
				targetCopy := ""
				gs.LastKick = &targetCopy
				gs.Leader = (newLeader == playerSnap.Name)
				gs.Promotion = false
				gs.SendPromotion = false
			})
			c.refreshGroupUI()
		}

	case category == pr.CategoryCombat:
		if strings.HasPrefix(trimmed, pr.CategoryCombat+" "+pr.TypeChat) {
			rest := strings.TrimSpace(strings.TrimPrefix(trimmed, pr.CategoryCombat+" "+pr.TypeChat))
			parts := strings.SplitN(rest, " ", 2)
			if len(parts) >= 2 {
				user := parts[0]
				msg := parts[1]
				c.gameState.UpdateCombatState(func(cs *state.CombatState) {
					cs.Chats = append(cs.Chats, state.CombatChat{Pseudo: user, Msg: msg})
				})
				c.ui.QueueUpdate(func() {
					c.ui.AppendCombatChat(user, msg)
				})
			}
		} else if strings.HasPrefix(trimmed, pr.CategoryCombat+" "+pr.TypeAllyTurn) {
			playerSnap := c.gameState.GetPlayerSnapshot()
			target := strings.TrimSpace(strings.TrimPrefix(trimmed, pr.CategoryCombat+" "+pr.TypeAllyTurn))
			c.gameState.UpdateCombatState(func(cs *state.CombatState) {
				cs.CurrentTurn = target
			})
			combatSnap := c.gameState.GetCombatSnapshot()
			c.ui.QueueUpdate(func() {
				c.ui.UpdateCombat(combatSnap)
			})
			if target == playerSnap.Name {
				c.sendToNetwork(pr.CmdCombatStats)
			}
		} else if strings.HasPrefix(trimmed, pr.CategoryCombat+" VICTORY") || strings.HasPrefix(trimmed, pr.CategoryCombat+" DEFEAT") {
			c.gameState.UpdateCombatState(func(cs *state.CombatState) {
				cs.InCombat = false
			})
			c.ui.QueueUpdate(func() {
				c.ui.ShowGamePage()
			})
			c.sendToNetwork(pr.CmdLook)
		} else if strings.HasPrefix(trimmed, pr.CategoryCombat+" ALLY_LEAVE_COMBAT") {
			leftUser := strings.TrimSpace(strings.TrimPrefix(trimmed, pr.CategoryCombat+" ALLY_LEAVE_COMBAT"))
			playerSnap := c.gameState.GetPlayerSnapshot()
			if leftUser == playerSnap.Name {
				c.gameState.UpdateCombatState(func(cs *state.CombatState) {
					cs.InCombat = false
				})
				c.ui.QueueUpdate(func() {
					c.ui.ShowGamePage()
				})
				c.sendToNetwork(pr.CmdLook)
			} else {
				c.sendToNetwork(pr.CmdCombatStats)
			}
		} else if strings.HasPrefix(trimmed, pr.CategoryCombat+" FIGHT_STARTED") || strings.HasPrefix(trimmed, pr.CategoryCombat+" UPDATE") {
			c.sendToNetwork(pr.CmdCombatStats)
		} else if res.Datas != nil {
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
		}

	case category == pr.CategoryGlobal || category == pr.CategoryRoom || category == pr.CategoryGroup:
		scope := category
		rest := strings.TrimSpace(strings.TrimPrefix(trimmed, scope))
		if strings.HasPrefix(rest, pr.TypeChat) {
			rest = strings.TrimSpace(strings.TrimPrefix(rest, pr.TypeChat))
		}
		parts := strings.SplitN(rest, " ", 2)
		if len(parts) >= 2 {
			user := parts[0]
			msg := parts[1]
			c.ui.QueueUpdate(func() {
				c.ui.AppendChat(scope, user, msg)
			})
		}

	case strings.HasPrefix(trimmed, pr.CategoryStats):
		parts := strings.SplitN(trimmed, " ", 2)
		if len(parts) == 2 {
			c.ui.QueueUpdate(func() {
				c.ui.UpdateDatas(parts[1])
			})
		}
	}
}
