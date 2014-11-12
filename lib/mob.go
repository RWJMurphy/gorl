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

	SetVisionRadius(int)
	VisionRadius() int
	Move(Vector)
	Tick(uint) MobAction

	Inventory() []Item
	AddToInventory(Item) bool
	DropItem(Item, *Dungeon) bool
	RemoveFromInventory(Item) bool

	MoveOrAct(Vector) bool
}

type mob struct {
	feature
	visionRadius int
	inventory    []Item
	// XXX I expect this will come out of sync. Head's up, future Reed.
	dungeon    *Dungeon
	lastTicked uint

	maxHealth  uint
	health     uint
	baseAttack uint

	fov        []Vector

	log *log.Logger
}

const MobDefaultHealth = 10
const MobDefaultAttack = 2

// NewMob creates and returns a new Mob
func NewMob(name string, char rune, log *log.Logger, dungeon *Dungeon) Mob {
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
		m,
		m.feature.String(),
		m.visionRadius,
	)
}

func (m *mob) Tick(turn uint) MobAction {
	action := MobAction{ActNone, nil}
	if m.Dead() {
		return action
	}
	if m.lastTicked != turn-1 {
		m.log.Panicf("%s out of sync! Last ticked: %d, ticking: %d", m, m.lastTicked, turn)
		return action
	}
	m.lastTicked = turn
	m.calculateFOV()
	var enemies, items []Feature

	for _, loc := range m.fov {
		fg := m.dungeon.FeatureGroup(loc)
		if fg.mob != nil && fg.mob != m {
			enemies = append(enemies, fg.mob)
		}
		if len(fg.items) > 0 {
			for _, item := range fg.items {
				items = append(items, item)
			}
		}
	}
	var direction Vector
	if len(enemies) > 0 {
		direction = enemies[0].Loc().Sub(m.Loc())
		m.log.Printf("%s@%s targeting %s@%s", m.Name(), m.Loc(), enemies[0].Name(), enemies[0].Loc())
	} else if len(items) > 0 {
		direction = items[0].Loc().Sub(m.Loc())
		m.log.Printf("%s@%s targeting %s@%s", m.Name(), m.Loc(), items[0].Name(), items[0].Loc())
	} else {
		direction = Vector{rand.Intn(3)-1, rand.Intn(3)-1}
		m.log.Printf("%s moving randomly", m.Name())
	}

	if direction.x != 0 {
		direction.x = direction.x / IntAbs(direction.x)
	}
	if direction.y != 0 {
		direction.y = direction.y / IntAbs(direction.y)
	}
	action.action = ActMove
	action.target = direction

	m.log.Printf("%s moving %s", m.Name(), direction)
	return action
}

func (m *mob) calculateFOV() {
	fov := make([]Vector, 0)
	m.dungeon.OnTilesInLineOfSight(m.loc, m.visionRadius, func(t *Tile, loc Vector){
		fov = append(fov, loc)
	})
	m.fov = fov
}

func (m *mob) SetVisionRadius(r int) {
	m.visionRadius = r
}

func (m *mob) VisionRadius() int {
	return m.visionRadius
}

func (m *mob) Move(movement Vector) {
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
	m.log.Printf("%s dropped %s on death", m.Name(), corpse)
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

func (m *mob) MoveOrAct(movement Vector) bool {
	destination := m.Loc().Add(movement)
	if otherMob := m.dungeon.MobAt(destination); otherMob != nil {
		if damageDealt, ok := m.Attack(otherMob); ok {
			m.log.Printf("%s hit %s for %d damage", m.Name(), otherMob.Name(), damageDealt)
			if otherMob.Dead() {
				m.log.Printf("The %s dies!", otherMob.Name())
			}
			return true
		}
		return false
	} else if moved := m.dungeon.MoveMob(m, movement); !moved {
		return false
	}
	return true
}
