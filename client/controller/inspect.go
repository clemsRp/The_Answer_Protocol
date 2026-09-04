package controller

import (
	"encoding/json"
	"fmt"
	"strings"
	"tap/protocol"
	pr "tap/protocol"
)

func (c *Controller) handleInspectResponse(res pr.ServerResponse) {
	if res.Datas == nil {
		return
	}
	fields := strings.Fields(res.Msg)
	if len(fields) < 2 {
		return
	}
	entityType := strings.ToUpper(fields[1])

	switch entityType {
	case pr.EntityTypeRoom:
		var roomData protocol.InspectRoomData
		raw, err := json.Marshal(res.Datas)
		if err != nil || json.Unmarshal(raw, &roomData) != nil {
			return
		}

		npcs := make(map[string]protocol.InspectNPCData, len(roomData.Npcs))
		for _, n := range roomData.Npcs {
			npcs[n.Id] = n
		}
		c.setNpcCache(npcs)
		// re-render the interaction panel now that we know who is hostile / has a quest
		c.refreshUI()

		formatted := formatRoom(roomData)
		c.ui.QueueUpdate(func() {
			c.ui.UpdateInspector(formatted)
		})

	case pr.EntityTypeSelf:
		var playerData protocol.InspectPlayerData
		raw, err := json.Marshal(res.Datas)
		if err == nil && json.Unmarshal(raw, &playerData) == nil {
			formatted := formatPlayer(playerData)
			c.ui.QueueUpdate(func() {
				c.ui.UpdateInspector(formatted)
			})
		}

	case pr.EntityTypeNpc:
		var npcData protocol.InspectNPCData
		raw, err := json.Marshal(res.Datas)
		if err == nil && json.Unmarshal(raw, &npcData) == nil {
			c.cacheNpc(npcData)
			c.refreshUI()

			formatted := formatNPC(npcData)
			c.ui.QueueUpdate(func() {
				c.ui.UpdateInspector(formatted)
			})
		}
	case pr.EntityTypeItem, pr.EntityTypeInventoryItem:
		var itemData protocol.InspectItemData
		raw, err := json.Marshal(res.Datas)
		if err == nil && json.Unmarshal(raw, &itemData) == nil {
			formatted := formatItem(itemData)
			c.ui.QueueUpdate(func() {
				c.ui.UpdateInspector(formatted)
			})
		}
	case pr.EntityTypePlayer:
		var playerData protocol.InspectPlayerData
		raw, err := json.Marshal(res.Datas)
		if err == nil && json.Unmarshal(raw, &playerData) == nil {
			formatted := formatPlayer(playerData)
			c.ui.QueueUpdate(func() {
				c.ui.UpdateInspector(formatted)
			})
		}
	}
}

// formatRoom renders a compact overview of everyone/everything inspectable
// in the current room, flagging hostile npcs and npcs with an available quest.
func formatRoom(data protocol.InspectRoomData) string {
	sb := &StatBuilder{}
	sb.WriteString("[blue]ROOM OVERVIEW[-]\n\n")

	sb.WriteString("[yellow]PLAYERS:[-]\n")
	if len(data.Players) == 0 {
		sb.WriteString("  [white]-[-]\n")
	} else {
		for _, p := range data.Players {
			status := "ok"
			if p.InCombat {
				status = "in combat"
			}
			sb.WriteString(fmt.Sprintf("  [white]- %s (%d/%d hp, %s)[-]\n", p.Name, p.Hp, p.MaxHp, status))
		}
	}

	sb.WriteString("\n[yellow]NPCS:[-]\n")
	if len(data.Npcs) == 0 {
		sb.WriteString("  [white]-[-]\n")
	} else {
		for _, n := range data.Npcs {
			tags := ""
			if n.Hostile {
				tags += " [red](hostile)[-]"
			}
			if n.QuestId != "" {
				tags += " [green](quest)[-]"
			}
			sb.WriteString(fmt.Sprintf("  [white]- %s[-]%s\n", n.Name, tags))
		}
	}

	sb.WriteString("\n[yellow]ITEMS:[-]\n")
	if len(data.Items) == 0 {
		sb.WriteString("  [white]-[-]\n")
	} else {
		for _, it := range data.Items {
			sb.WriteString(fmt.Sprintf("  [white]- %s[-]\n", it.Name))
		}
	}

	return strings.TrimSpace(sb.String())
}

type StatBuilder struct {
	strings.Builder
}

func (sb *StatBuilder) AddString(label, val string) {
	if val != "" {
		sb.WriteString(fmt.Sprintf("[blue]%s:[-] [white]%s[-]\n", strings.ToUpper(label), val))
	}
}

func (sb *StatBuilder) AddInt(label string, val int) {
	if val != 0 {
		sb.WriteString(fmt.Sprintf("[blue]%s:[-] [white]%d[-]\n", strings.ToUpper(label), val))
	}
}

func (sb *StatBuilder) AddBool(label string, val bool) {
	if val {
		sb.WriteString(fmt.Sprintf("[blue]%s:[-] [green]YES[-]\n", strings.ToUpper(label)))
	} else {
		sb.WriteString(fmt.Sprintf("[blue]%s:[-] [white]NO[-]\n", strings.ToUpper(label)))
	}
}

func formatNPC(data protocol.InspectNPCData) string {
	sb := &StatBuilder{}
	name := data.Name
	if name == "" && data.Id != "" {
		name = data.Id
	}
	sb.AddString("NAME", name)

	if data.Description != "" {
		sb.WriteString(fmt.Sprintf("\n[white]%s[-]\n\n", data.Description))
	}

	sb.AddString("ROLE", data.Role)
	sb.AddString("QUEST ID", data.QuestId)

	if data.HpMax != 0 {
		sb.WriteString(fmt.Sprintf("[blue]HP:[-] [white]%d / %d[-]\n", data.Hp, data.HpMax))
	}

	sb.AddInt("DAMAGE", data.Damage)
	sb.AddInt("XP REWARD", data.XpReward)

	if len(data.ItemsReward) > 0 {
		sb.AddString("DROPS", strings.Join(data.ItemsReward, ", "))
	}

	sb.AddBool("HOSTILE", data.Hostile)
	sb.AddBool("COMBAT", data.InCombat)

	return strings.TrimSpace(sb.String())
}

func formatItem(data protocol.InspectItemData) string {
	sb := &StatBuilder{}
	sb.AddString("NAME", data.Name)

	if data.Description != "" {
		sb.WriteString(fmt.Sprintf("\n[white]%s[-]\n\n", data.Description))
	}

	sb.AddString("TYPE", data.Type)
	sb.AddInt("WORTH", data.Worth)
	sb.AddInt("DAMAGE", data.Damage)

	sb.AddString("STAT", data.TargetStat)
	sb.AddInt("AMOUNT", data.Amount)
	sb.AddString("EFFECT", data.EffectType)
	sb.AddInt("DURATION", data.Duration)
	sb.AddInt("CHARGES", data.Charges)

	sb.AddBool("TRADABLE", data.Tradable)
	sb.AddBool("OBTAINABLE", data.Obtainable)

	return strings.TrimSpace(sb.String())
}

func formatPlayer(data protocol.InspectPlayerData) string {
	sb := &StatBuilder{}
	sb.AddString("NAME", data.Name)

	if data.MaxHp != 0 {
		sb.WriteString(fmt.Sprintf("\n[blue]HP:[-] [white]%d / %d[-]\n", data.Hp, data.MaxHp))
	}

	sb.AddString("STATUS", data.Status)
	sb.AddBool("COMBAT", data.InCombat)
	sb.AddBool("GROUP", data.IsInGroup)

	return strings.TrimSpace(sb.String())
}
