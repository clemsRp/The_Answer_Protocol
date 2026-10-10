package engine

type Item struct {
	Id          string `json:"-" validate:"-"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Obtainable  bool   `json:"obtainable"`
	Type        string `json:"type" validate:"required,valid_item_type"`
	Tradable    bool   `json:"tradable"`
	Worth       int    `json:"worth" validate:"gte=0"`

	Damage int `json:"damage,omitempty" validate:"required_if=Type weapon,omitempty,gt=0"`

	TargetStat string `json:"target_stat,omitempty" validate:"required_if=Type consumable,omitempty,valid_target_stat"`
	Amount     int    `json:"amount,omitempty" validate:"required_if=Type consumable,omitempty,gt=0"`
	EffectType string `json:"effect_type,omitempty" validate:"required_if=Type consumable,omitempty,valid_effect_type"`
	Duration   int    `json:"duration,omitempty" validate:"required_if=Type consumable,omitempty,gt=0"`
	Charges    int    `json:"charges,omitempty" validate:"required_if=Type consumable,omitempty,gt=0"`
}
type Weapon struct {
	*Item
	Damage int `json:"damage"`
}

type Ressource struct {
	*Item
}

type Consumable struct {
	*Item
	TargetStat string `json:"target_stat"`
	Amount     int    `json:"amount"`
	EffectType string `json:"effect_type"`
	Duration   int    `json:"duration"`
	Charges    int    `json:"charges"`
}

func (item *Item) Clone() *Item {
	// no slices or objects in item so we make a quick copy
	cp := *item
	return &cp
}

func (item *Item) ConvertToWeapon() *Weapon {
	if item.Type != "weapon" {
		return nil
	}
	return &Weapon{
		Item:   item,
		Damage: item.Damage,
	}
}
