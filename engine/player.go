package engine

import (
	"errors"
	pr "tap/protocol"
)

type Player struct {
	ip             string
	name           string
	room           *Room
	group          string
	inventory      []*Item
	quests         []*Quest
	invitations    []string
	promotion      bool
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

func (p *Player) takeDamage(amount int) {
	p.stats.Hp -= amount
}

func (p *Player) playCombatTurn(target Fighter) *CombatTurnResult {
	target.takeDamage(p.equippedWeapon.Damage)
	turn_result := &CombatTurnResult{AttackerHp: p.getHp(), TargetHp: target.getHp(), Damage: p.equippedWeapon.Damage, Status: "combat"}
	return turn_result
}
func (p *Player) getName() string {
	return p.name
}
func (p *Player) getInitiative() int {
	return p.stats.Initiative
}
func (e *Engine) createNewPlayerInstance(pseudo string, ip string) (*Player, error) {
	base_item, exists := e.world.Items["sword"]
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
		ip: ip,
	}, nil
}
