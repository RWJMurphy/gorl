package gorl

import (
	"math/rand"
)

type Coord struct {
	x, y int
}

type BitSet int64

const (
	FlagCrossable = iota << 1
)

type Tile struct {
	c rune
	flags BitSet
}

func NewTile(c rune, flags BitSet) Tile {
	t := Tile{c, flags}
	return t
}

var InvalidTile = Tile{' ', BitSet(0)}

type Dungeon struct {
	width, height int
	origin        Coord
	tiles         [][]Tile
	mobs          []*Player
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
		make([]*Player, 0, int(float32(size)*est_mob_ratio)),
	}

	var tile Tile
	for x := 0; x < width; x++ {
		for y := 0; y < width; y++ {
			if rand.Float64() <= 0.1 {
				tile = NewTile('#', BitSet(0))
			} else {
				tile = NewTile('.', BitSet(0 & FlagCrossable))
			}
			m.tiles[y][x] = tile
		}
	}

	return m
}

func (d *Dungeon) AddMob(m *Player) {
	d.mobs = append(d.mobs, m)
}

func (d *Dungeon) Tile(x, y int) Tile {
	if x < 0 || x >= d.width || y < 0 || y >= d.height {
		return InvalidTile
	}
	return d.tiles[y][x]
}
