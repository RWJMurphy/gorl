package gorl

import (
	"fmt"
	"math/rand"
)

// Mob is a Movable OBject.
type Mob interface {
	Feature
	VisionRadius() int
	Move(Movement)
	Tick(*Dungeon) bool
}

type mob struct {
	feature
	visionRadius int
}

// NewMob creates and returns a new Mob
func NewMob(name string, char rune) Mob {
	m := &mob{}
	m.feature = *NewFeature(name, char).(*feature)
	return m
}

func (m *mob) Tick(d *Dungeon) bool {
	dx, dy := rand.Intn(3)-1, rand.Intn(3)-1
	return d.MoveMob(m, Movement{dx, dy})
}

func (m *mob) VisionRadius() int {
	return m.visionRadius
}

func (m *mob) Move(movement Movement) {
	m.loc.x += movement.x
	m.loc.y += movement.y
}

func (m *mob) String() string {
	return fmt.Sprintf(
		"<Mob feature:%s, visionRadius:%d>",
		m.feature,
		m.visionRadius,
	)
}
