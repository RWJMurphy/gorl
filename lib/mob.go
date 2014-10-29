package gorl

type Mob interface {
	Move(Movement)
	Loc() Coord
	Char() rune
}

type mob struct {
	loc  Coord
	char rune
}

func NewMob(char rune) *mob {
	m := &mob{}
	m.char = char
	return m
}

func (m *mob) Loc() Coord {
	return m.loc
}

func (m *mob) Char() rune {
	return m.char
}

func (m *mob) Move(movement Movement) {
	m.loc.x += movement.x
	m.loc.y += movement.y
}
