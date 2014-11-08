package gorl

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/nsf/termbox-go"
)

// Mob is a Movable OBject.
type Mob interface {
	Feature
	Attacker
	Defender
	VisionRadius() int
	Move(Vec)
	Tick(uint) bool

	Inventory() []Item
	AddToInventory(Item) bool
	DropItem(Item, *Dungeon) bool
	RemoveFromInventory(Item) bool
}

type mob struct {
	feature
	visionRadius int
	inventory    []Item
	// BUG I expect this will come out of sync. Head's up, future Reed.
	dungeon    *Dungeon
	lastTicked uint

	maxHealth  uint
	health     uint
	baseAttack uint

	log log.Logger
}

const MobDefaultHealth = 10
const MobDefaultAttack = 2

// NewMob creates and returns a new Mob
func NewMob(name string, char rune, log log.Logger, dungeon *Dungeon) Mob {
	m := &mob{}
	m.feature = *NewFeature(name, char).(*feature)
	m.inventory = make([]Item, 0)
	m.log = log
	m.maxHealth = MobDefaultHealth
	m.health = m.maxHealth
	m.baseAttack = MobDefaultAttack
	m.dungeon = dungeon
	return m
}

func (m *mob) String() string {
	return fmt.Sprintf(
		"<Mob@%p feature:%s, visionRadius:%d>",
		&m,
		m.feature.String(),
		m.visionRadius,
	)
}

func (m *mob) Tick(turn uint) bool {
	if m.Dead() {
		return false
	}
	if m.lastTicked != turn-1 {
		m.log.Panicf("%s out of sync! Last ticked: %d, ticking: %d", m, m.lastTicked, turn)
	}
	m.lastTicked = turn
	dx, dy := rand.Intn(3)-1, rand.Intn(3)-1
	return m.dungeon.MoveMob(m, Vec{dx, dy})
}

func (m *mob) VisionRadius() int {
	return m.visionRadius
}

func (m *mob) Move(movement Vec) {
	m.loc.x += movement.x
	m.loc.y += movement.y
}

func (m *mob) LightRadius() int {
	max := m.feature.LightRadius()
	for _, i := range m.inventory {
		if i.LightRadius() > max {
			max = i.LightRadius()
		}
	}
	return max
}

func (m *mob) Inventory() []Item {
	inv := make([]Item, len(m.inventory))
	copy(inv, m.inventory)
	return inv
}

func (m *mob) AddToInventory(i Item) bool {
	m.inventory = append(m.inventory, i)
	return true
}

func (m *mob) DropItem(item Item, d *Dungeon) bool {
	if m.RemoveFromInventory(item) {
		item.SetLoc(m.Loc())
		d.AddItem(item)
		return true
	}
	m.log.Panicf("%s tried to drop unheld item %s!", m, item)
	return false
}

func (m *mob) RemoveFromInventory(item Item) bool {
	for i, inventoryItem := range m.inventory {
		if inventoryItem == item {
			copy(m.inventory[i:], m.inventory[i+1:])
			m.inventory[len(m.inventory)-1] = nil
			m.inventory = m.inventory[:len(m.inventory)-1]
			return true
		}
	}
	return false
}

func (m *mob) die() {
	// on death, drop corpse
	corpse := NewItem("corpse", '%', 100)
	corpse.SetColor(termbox.ColorRed)
	m.AddToInventory(corpse)
	m.DropItem(corpse, m.dungeon)
	m.log.Printf("%s dropped %s on death", m, corpse)
}

func (m *mob) Dead() bool {
	return m.health <= 0
}

func (m *mob) AttackStrength() uint {
	return m.baseAttack
}

func (m *mob) AttackedFor(damage uint) uint {
	if damage >= m.health {
		m.health = 0
		m.die()
	} else {
		m.health -= damage
	}
	return damage
}

func (m *mob) Attack(d Defender) (uint, bool) {
	if !d.Dead() {
		damageDealt := d.AttackedFor(m.AttackStrength())
		return damageDealt, true
	}
	return 0, false
}
