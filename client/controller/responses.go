package controller

import (
	"encoding/json"
	"fmt"
	"strings"
	"tap/client/state"
	"tap/protocol"
	pr "tap/protocol"
)

func (c *Controller) handleCommandResponses(res pr.ServerResponse) {
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
		c.sendToNetwork(pr.CmdQuests)

	case (lastCmdBase == pr.CmdLook || lastCmdBase == pr.CmdMove) && res.Datas != nil:
		var lookData protocol.LookCommandData
		raw, err := json.Marshal(res.Datas)
		if err == nil && json.Unmarshal(raw, &lookData) == nil {
			c.gameState.UpdateRoomLook(&lookData)
			c.refreshUI()
			c.sendToNetwork(pr.CmdInspect)
			c.sendToNetwork(pr.CmdInventory)
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

	case (lastCmd == pr.CmdCombatStats || strings.EqualFold(strings.TrimSpace(lastCmd), pr.CmdCombatStats)) && res.Datas != nil:
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
			if myData, ok := combatData.Team[c.ui.GetPseudo()]; ok {
				c.gameState.UpdatePlayer(func(p *state.Player) {
					p.Inventory = append([]string{}, myData.Inventory...)
				})
			}
			combatSnap := c.gameState.GetCombatSnapshot()
			c.ui.QueueUpdate(func() {
				c.ui.UpdateCombat(combatSnap)
				c.ui.ShowCombatPage()
			})
		}

	// Regroupement Attack et UseItem pour gérer la fin du combat sur les 2 actions
	case (lastCmdBase == pr.CmdAttack || strings.HasPrefix(lastCmd, pr.CmdAttack) || strings.HasPrefix(lastCmd, pr.CmdChatCombatAttack) || lastCmdBase == pr.CmdUseItem || strings.HasPrefix(lastCmd, pr.CmdUseItem)) && (res.Msg == pr.MsgOK || strings.HasPrefix(res.Msg, pr.MsgOK)):

		if res.Datas != nil {
			var fullTurn struct {
				CombatState string   `json:"combat_state"`
				XpReward    int      `json:"xp_reward,omitempty"`
				ItemsReward []string `json:"items_reward,omitempty"`
			}
			raw, err := json.Marshal(res.Datas)
			if err == nil && json.Unmarshal(raw, &fullTurn) == nil {
				if fullTurn.CombatState == "VICTORY" || fullTurn.CombatState == "DEFEAT" {
					c.gameState.UpdateCombatState(func(cs *state.CombatState) {
						cs.InCombat = false
					})
					// Construire la liste des récompenses à afficher
					rewards := make([]string, 0)
					if fullTurn.XpReward > 0 {
						rewards = append(rewards, fmt.Sprintf("%d XP", fullTurn.XpReward))
					}
					rewards = append(rewards, fullTurn.ItemsReward...)
					combatResult := fullTurn.CombatState
					c.ui.QueueUpdate(func() {
						c.ui.ShowCombatResultPopup(combatResult, rewards)
					})
					c.sendToNetwork(pr.CmdLook)
					c.sendToNetwork(pr.CmdInventory)
					// A defeated npc may fulfil a quest target: refresh
					// progress automatically.
					c.sendToNetwork(pr.CmdQuests)
					break // Empêche de demander les Stats d'un combat terminé
				}
			}
		}
		c.sendToNetwork(pr.CmdCombatStats)

	case lastCmdBase == pr.CmdTalk && res.Msg == pr.MsgOK:
		npcName := ""
		if len(cmdFields) >= 2 {
			npcName = cmdFields[1]
		}
		if npcName != "" && res.Datas != nil {
			dialogue, ok := res.Datas.(string)
			if ok {
				c.gameState.UpdatePlayer(func(p *state.Player) {
					if p.NpcDialogues == nil {
						p.NpcDialogues = make(map[string]string)
					}
					p.NpcDialogues[npcName] = dialogue
				})
				c.refreshUI()
				c.sendToNetwork(pr.CmdLook)
			}
		}

	case (lastCmdBase == pr.CmdFlee || strings.HasPrefix(lastCmd, pr.CmdFlee) || strings.HasPrefix(lastCmd, pr.CmdChatCombatFlee)) && (res.Msg == pr.MsgOK || strings.HasPrefix(res.Msg, pr.MsgOK)):
		c.gameState.UpdateCombatState(func(cs *state.CombatState) {
			cs.InCombat = false
		})
		c.ui.QueueUpdate(func() {
			c.ui.ShowGamePage()
		})
		c.sendToNetwork(pr.CmdLook)
		c.sendToNetwork(pr.CmdInventory)
		c.sendToNetwork(pr.CmdQuests)

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
		c.sendToNetwork(pr.CmdLook)
		// Taking an item can fulfil a quest target: refresh progress
		// automatically instead of waiting for the player to check.
		c.sendToNetwork(pr.CmdQuests)

	case strings.HasPrefix(lastCmd, pr.CmdDrop) && strings.HasPrefix(res.Msg, pr.MsgOK):
		c.sendToNetwork(pr.CmdLook)
		c.sendToNetwork(pr.CmdQuests)

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

	case lastCmdBase == pr.CmdInspect && strings.HasPrefix(res.Msg, pr.MsgOK):
		if res.Datas != nil {
			c.handleInspectResponse(res)
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

		case lastCmdBase == pr.CmdCompleteQuest:
			// Récupérer la récompense si disponible
			reward := ""
			questID := ""
			if len(strings.Fields(lastCmd)) >= 2 {
				questID = strings.Fields(lastCmd)[1]
			}
			if res.Datas != nil {
				var questData protocol.QuestData
				raw, err := json.Marshal(res.Datas)
				if err == nil && json.Unmarshal(raw, &questData) == nil {
					reward = questData.Reward
					if questData.Id != "" {
						questID = questData.Id
					}
				}
			}
			capturedQuestID := questID
			capturedReward := reward
			c.ui.QueueUpdate(func() {
				c.ui.ShowQuestCompletedPopup(capturedQuestID, capturedReward)
			})
			// Fetch updated quests list
			c.sendToNetwork(pr.CmdQuests)
			c.sendToNetwork(pr.CmdLook)
			// Completing a quest can consume the target item from the
			// inventory (server-side), so resync it too.
			c.sendToNetwork(pr.CmdInventory)

		case lastCmdBase == pr.CmdQuest:

		case lastCmdBase == pr.CmdQuests:
			var data []protocol.TrackedQuestData
			raw, err := json.Marshal(res.Datas)
			if err == nil && json.Unmarshal(raw, &data) == nil {
				c.gameState.UpdatePlayer(func(p *state.Player) {
					p.Quests = data
				})
				c.ui.QueueUpdate(func() {
					c.ui.UpdateQuests(data)
				})
			}
		}
	}
}
