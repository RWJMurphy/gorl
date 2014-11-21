package gorl

type Attacker interface {
	AttackStrength() uint
	Attack(Defender) (uint, bool)
}

type Defender interface {
	AttackedFor(uint) uint
	Dead() bool
}

type Wielder interface {
	WieldPoints() []string
	Wielding() []Wieldable

	Wield(Wieldable, uint) bool
	Unwield(uint) bool
}

type Wieldable interface {
	Item
	AttackStrength() uint
}
