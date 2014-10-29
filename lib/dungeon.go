package gorl

import (
	"math/rand"
)

type Coord struct {
	x, y int
}

type BitSet int

type Tile struct {
	c rune
	flags BitSet
}

const (
	FlagCrossable BitSet = 1 << iota
	FlagLit
	FlagVisible
)

func NewTile(c rune, flags BitSet) Tile {
	t := Tile{c, flags}
	return t
}

var InvalidTile = Tile{' ', BitSet(0)}

type Dungeon struct {
	width, height int
	origin        Coord
	tiles         [][]Tile
	mobs          []Mob
}

const est_mob_ratio = 0.1

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
		make([]Mob, 0, int(float32(size)*est_mob_ratio)),
	}

	var tile Tile
	for x := 0; x < width; x++ {
		for y := 0; y < width; y++ {
			if rand.Float64() <= 0.1 {
				tile = NewTile('#', BitSet(0)|FlagVisible)
			} else {
				tile = NewTile('.', BitSet(0)|FlagCrossable|FlagVisible)
			}
			m.tiles[y][x] = tile
		}
	}

	return m
}

func (d *Dungeon) AddMob(m Mob) {
	d.mobs = append(d.mobs, m)
}

func (d *Dungeon) CalculateLighting() {
	for x := 0; x < d.width; x++ {
		for y := 0; y < d.width; y++ {
			d.tiles[y][x].flags = d.tiles[y][x].flags & ^FlagLit
		}
	}
	for _, m := range d.mobs {
		for x := m.Loc().x - 5; x < m.Loc().x + 5; x++ {
			for y := m.Loc().y - 5; y < m.Loc().y + 5; y++ {
				d.tiles[y][x].flags = d.tiles[y][x].flags | FlagLit
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
