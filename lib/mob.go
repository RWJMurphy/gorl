package gorl

import "fmt"

type Mob interface {
	Loc() Coord
	Char() rune
	Flags() Flag
	VisionRadius() int
	LightRadius() int

	Move(Movement)
}

type mob struct {
	loc          Coord
	char         rune
	flags        Flag
	visionRadius int
	lightRadius  int
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

func (m *mob) VisionRadius() int {
	return m.visionRadius
}

func (m *mob) Move(movement Movement) {
	m.loc.x += movement.x
	m.loc.y += movement.y
}

func (m *mob) Flags() Flag {
	return m.flags
}

func (m *mob) LightRadius() int {
	return m.lightRadius
}

func (m *mob) String() string {
	return fmt.Sprintf(
		"<Mob char:%c, loc:%s, flags:%s, visionRadius:%d>",
		m.char,
		m.loc,
		m.flags,
		m.visionRadius,
	)
}
