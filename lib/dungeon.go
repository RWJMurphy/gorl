package gorl

import (
	"strings"
	"fmt"
	"math/rand"
)

type Coord struct {
	x, y int
}

func (c Coord) String() string {
	return fmt.Sprintf("<Coord x:%d, y:%d>", c.x, c.y)
}

type Tile struct {
	c     rune
	flags Flag
}

type Flag uint8

const (
	FlagCrossable Flag = 1 << iota
	FlagLit
	FlagVisible
)

func (f Flag) String() string {
	on_flags := make([]string, 0)

	if f&FlagCrossable != 0 { on_flags = append(on_flags, "Crossable") }
	if f&FlagLit != 0 { on_flags = append(on_flags, "Lit") }
	if f&FlagVisible != 0 { on_flags = append(on_flags, "Visible") }

	if len(on_flags) == 0 {
		on_flags = append(on_flags, "None")
	}

	return fmt.Sprintf("<Flag %s>", strings.Join(on_flags, "|"))
}

func NewTile(c rune, flags Flag) Tile {
	t := Tile{c, flags}
	return t
}

var InvalidTile = Tile{' ', Flag(0)}

type Dungeon struct {
	width, height int
	origin        Coord
	tiles         [][]Tile
	mobs          map[Coord]Mob
}

func NewDungeon(width, height int) *Dungeon {
	size := width * height
	tiles := make([][]Tile, height)
	tiles_raw := make([]Tile, size)
	for i := range tiles {
		tiles[i], tiles_raw = tiles_raw[:width], tiles_raw[width:]
	}

	m := &Dungeon{
		width, height,
		Coord{width / 2, height / 2},
		tiles,
		make(map[Coord]Mob),
	}

	var tile Tile
	for x := 0; x < width; x++ {
		for y := 0; y < width; y++ {
			if rand.Float64() <= 0.1 {
				tile = NewTile('#', Flag(0)|FlagVisible)
			} else {
				tile = NewTile('.', Flag(0)|FlagCrossable|FlagVisible)
			}
			m.tiles[y][x] = tile
		}
	}

	return m
}

func (d *Dungeon) AddMob(mob Mob) {
	if other_mob, exists := d.mobs[mob.Loc()]; exists {
		panic(fmt.Sprintf(
			"Tried to put two mobs on same location: %s, %s",
			mob,
			other_mob))
	}
	d.mobs[mob.Loc()] = mob
}

func (d *Dungeon) DeleteMob(mob Mob) {
	if _, exists := d.mobs[mob.Loc()]; exists {
		delete(d.mobs, mob.Loc())
	} else {
		panic(fmt.Sprintf("Tried to delete non-existent mob: %s", mob))
	}
}

func (d *Dungeon) MoveMob(mob Mob, move Movement) {
	d.DeleteMob(mob)
	mob.Move(move)
	d.AddMob(mob)
}

func (d *Dungeon) CalculateLighting() {
	for x := 0; x < d.width; x++ {
		for y := 0; y < d.width; y++ {
			d.tiles[y][x].flags = d.tiles[y][x].flags & ^FlagLit
		}
	}
	for loc, _ := range d.mobs {
		for x := loc.x - 5; x < loc.x+5; x++ {
			for y := loc.y - 5; y < loc.y+5; y++ {
				if x >= 0 && x < d.width && y >= 0 && y < d.height {
					d.tiles[y][x].flags = d.tiles[y][x].flags | FlagLit
				}
			}
		}
	}
}

func (d *Dungeon) Tile(x, y int) Tile {
	if x < 0 || x >= d.width || y < 0 || y >= d.height {
		return InvalidTile
	}
	return d.tiles[y][x]
}
