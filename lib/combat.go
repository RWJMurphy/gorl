package gorl

type Attacker interface {
	AttackStrength() uint
	Attack(Defender) (uint, bool)
}

type Defender interface {
	AttackedFor(uint) uint
	Dead() bool
}
