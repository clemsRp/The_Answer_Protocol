// engine/inspect.go
package engine

import (
	"errors"
	"slices"
	"strings"
	pr "tap/protocol"
)

func buildInspectPlayerData(p *Player) pr.InspectPlayerData {
	return pr.InspectPlayerData{
		Name:      p.name,
		IsInGroup: p.group != "",
		InCombat:  p.inCombat,
		Hp:        p.stats.Hp,
		MaxHp:     p.stats.HpMax,
		Status:    p.stats.Status,
	}
}

func buildInspectNpcData(n *Npc) pr.InspectNPCData {
	data := pr.InspectNPCData{
		Id:          n.Id,
		Name:        n.Name,
		Description: n.Description,
		Dialogue:    n.Dialogue,
		Role:        n.Role,
		QuestId:     n.QuestId,
		Hostile:     n.Hostile,
		Damage:      n.Damage,
		InCombat:    n.InCombat,
		XpReward:    n.XpReward,
	}

	if n.Stats != nil {
		data.Hp = n.Stats.Hp
		data.HpMax = n.Stats.HpMax
	}

	for _, reward := range n.ItemsReward {
		data.ItemsReward = append(data.ItemsReward, reward.Id)
	}

	return data
}

func buildInspectItemData(it *Item) pr.InspectItemData {
	return pr.InspectItemData{
		Id:          it.Id,
		Name:        it.Name,
		Description: it.Description,
		Obtainable:  it.Obtainable,
		Type:        it.Type,
		Tradable:    it.Tradable,
		Worth:       it.Worth,
		Damage:      it.Damage,
		TargetStat:  it.TargetStat,
		Amount:      it.Amount,
		EffectType:  it.EffectType,
		Duration:    it.Duration,
		Charges:     it.Charges,
	}
}

// inspectRoom gathers inspect datas for every player, npc and item currently
// present in the player's room (defeated npcs are excluded, same as LOOK).
func (e *Engine) inspectRoom(player *Player) pr.InspectRoomData {
	res := pr.InspectRoomData{
		Players: make([]pr.InspectPlayerData, 0),
		Npcs:    make([]pr.InspectNPCData, 0),
		Items:   make([]pr.InspectItemData, 0),
	}

	for pseudo, p := range e.players {
		if p.room == player.room && pseudo != player.name {
			res.Players = append(res.Players, buildInspectPlayerData(p))
		}
	}

	for _, npcName := range player.room.Npcs {
		if slices.Contains(player.DefeatedNpcs, npcName) {
			continue
		}
		if npc, exists := e.world.Npcs[npcName]; exists {
			npcData := buildInspectNpcData(npc)
			for _, cs := range e.activeCombats {
				if cs.RoomId == player.room.Id {
					for _, combatNpc := range cs.Npcs {
						if combatNpc.Id == npcName || combatNpc.Name == npcName {
							npcData.InCombat = true
							if combatNpc.Stats != nil {
								npcData.Hp = combatNpc.Stats.Hp
								npcData.HpMax = combatNpc.Stats.HpMax
							}
						}
					}
				}
			}
			res.Npcs = append(res.Npcs, npcData)
		}
	}

	for _, itemName := range player.room.Items {
		if item, exists := e.world.Items[itemName]; exists {
			res.Items = append(res.Items, buildInspectItemData(item))
		}
	}

	return res
}

func (e *Engine) handleCmdInspect(player *Player, req []string) (string, any, error) {
	// bare INSPECT: return datas about everything in the current room
	if len(req) == 1 {
		return "OK " + pr.EntityTypeRoom, e.inspectRoom(player), nil
	}

	// INSPECT SELF: return the caller's own data
	if len(req) == 2 && strings.ToUpper(req[1]) == pr.EntityTypeSelf {
		return "OK " + pr.EntityTypeSelf, buildInspectPlayerData(player), nil
	}

	if len(req) != 3 {
		return "", nil, errors.New(pr.ErrInvalidCommand)
	}

	entityType := strings.ToUpper(req[1])
	target := req[2]

	switch entityType {
	case pr.EntityTypePlayer:
		for pseudo, p := range e.players {
			if pseudo == target && p.room == player.room {
				return "OK " + pr.EntityTypePlayer, buildInspectPlayerData(p), nil
			}
		}
		return "", nil, errors.New(pr.ErrUnknownUser)

	case pr.EntityTypeNpc:
		if !isNpcInRoom(player.room, target) || slices.Contains(player.DefeatedNpcs, target) {
			return "", nil, errors.New(pr.ErrNpcNotFound)
		}
		npc, exists := e.world.Npcs[target]
		if !exists {
			return "", nil, errors.New(pr.ErrNpcNotFound)
		}
		npcData := buildInspectNpcData(npc)
		for _, cs := range e.activeCombats {
			if cs.RoomId == player.room.Id {
				for _, combatNpc := range cs.Npcs {
					if combatNpc.Id == target || combatNpc.Name == target {
						npcData.InCombat = true
						if combatNpc.Stats != nil {
							npcData.Hp = combatNpc.Stats.Hp
							npcData.HpMax = combatNpc.Stats.HpMax
						}
					}
				}
			}
		}
		return "OK " + pr.EntityTypeNpc, npcData, nil

	case pr.EntityTypeItem:
		if !slices.Contains(player.room.Items, target) {
			return "", nil, errors.New(pr.ErrItemNotFound)
		}
		item, exists := e.world.Items[target]
		if !exists {
			return "", nil, errors.New(pr.ErrItemNotFound)
		}
		return "OK " + pr.EntityTypeItem, buildInspectItemData(item), nil

	case pr.EntityTypeInventoryItem:
		for _, item := range player.inventory {
			if item.Id == target {
				return "OK " + pr.EntityTypeInventoryItem, buildInspectItemData(item), nil
			}
		}
		return "", nil, errors.New(pr.ErrItemNotInInventory)

	default:
		return "", nil, errors.New(pr.ErrInvalidCommand)
	}
}
