package gorl

type Weapon interface {
	Wieldable
}

type weapon struct {
	item
	attackStrength uint
}

func NewWeapon(name string, char rune, weight int, attackStrength uint) Weapon {
	w := weapon{
		item{
			*NewFeature(name, char).(*feature),
			weight,
		},
		attackStrength,
	}
	w.flags |= FlagCrossable
	return &w
}

func (w *weapon) AttackStrength() uint {
	return w.attackStrength
}
