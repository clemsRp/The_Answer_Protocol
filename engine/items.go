package engine

type Item struct {
	Id          string `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Obtainable  bool   `json:"obtainable"`
	Type        string `json:"type"`
	Tradable    bool   `json:"tradable"`
	Worth       int    `json:"worth"`

	Damage int `json:"damage,omitempty"`

	TargetStat string `json:"target_stat,omitempty"`
	Amount     int    `json:"amount,omitempty"`
	EffectType string `json:"effect_type,omitempty"`
	Duration   int    `json:"duration,omitempty"`
	Charges    int    `json:"charges,omitempty"`
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
