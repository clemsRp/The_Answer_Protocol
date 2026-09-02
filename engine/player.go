package engine

import (
	"errors"
	pr "tap/protocol"
)

type Player struct {
	id             string
	name           string
	room           *Room
	group          string
	inventory      []*Item
	quests         []*Quest
	invitations    []string
	DefeatedNpcs   []string
	promotion      bool
	send_promotion bool
	inCombat       bool
	equippedWeapon *Weapon
	stats          *CombatStats
	Fighter
}

type Group struct {
	id      string
	leader  *Player
	players []*Player
}

func (p *Player) isDead() bool {
	if p.stats.Hp <= 0 {
		return true
	}
	return false
}

func (p *Player) getHp() int {
	return p.stats.Hp
}

func (p *Player) takeDamage(amount int) int {
	p.stats.Hp -= amount
	return amount
}

func (p *Player) getDamage() int {
	return p.equippedWeapon.Damage
}

func (p *Player) getName() string {
	return p.name
}

func (p *Player) getInitiative() int {
	return p.stats.Initiative
}
func (e *Engine) createNewPlayerInstance(pseudo string, id string) (*Player, error) {
	base_item, exists := e.world.Items["plastic_hanger"]
	if !exists {
		return nil, errors.New(pr.ErrInternalServer)
	}

	start_item_copy := base_item.Clone()

	weapon_start := start_item_copy.ConvertToWeapon()
	if weapon_start == nil {
		return nil, errors.New(pr.ErrInternalServer)
	}

	return &Player{
		name:           pseudo,
		room:           e.world.Rooms[RoomEntrance],
		equippedWeapon: weapon_start,
		stats: &CombatStats{
			Hp:         100,
			HpMax:      100,
			Initiative: 100,
			Status:     StatusNormal,
		},
		id: id,
	}, nil
}
