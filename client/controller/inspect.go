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
	case pr.EntityTypeNpc:
		var npcData protocol.InspectNPCData
		raw, err := json.Marshal(res.Datas)
		if err == nil && json.Unmarshal(raw, &npcData) == nil {
			formatted := formatNPC(npcData)
			c.ui.QueueUpdate(func() {
				c.ui.UpdateDatas(formatted)
			})
		}
	case pr.EntityTypeItem, pr.EntityTypeInventoryItem:
		var itemData protocol.InspectItemData
		raw, err := json.Marshal(res.Datas)
		if err == nil && json.Unmarshal(raw, &itemData) == nil {
			formatted := formatItem(itemData)
			c.ui.QueueUpdate(func() {
				c.ui.UpdateDatas(formatted)
			})
		}
	case pr.EntityTypePlayer:
		var playerData protocol.InspectPlayerData
		raw, err := json.Marshal(res.Datas)
		if err == nil && json.Unmarshal(raw, &playerData) == nil {
			formatted := formatPlayer(playerData)
			c.ui.QueueUpdate(func() {
				c.ui.UpdateDatas(formatted)
			})
		}
	}
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
