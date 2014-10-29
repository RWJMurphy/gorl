package gorl

import (
	"math/rand"
)

type Coord struct {
	x, y int
}

type Dungeon struct {
	width, height int
	origin        Coord
	tiles         [][]rune
	mobs          []*Player
}

const est_mob_ratio = 0.1

func NewDungeon(width, height int) *Dungeon {
	size := width * height
	tiles := make([][]rune, height)
	tiles_raw := make([]rune, size)
	for i := range tiles {
		tiles[i], tiles_raw = tiles_raw[:width], tiles_raw[width:]
	}

	m := &Dungeon{
		width, height,
		Coord{width / 2, height / 2},
		tiles,
		make([]*Player, 0, int(float32(size)*est_mob_ratio)),
	}

	for x := 0; x < width; x++ {
		for y := 0; y < width; y++ {
			if rand.Float64() <= 0.1 {
			    m.tiles[y][x] = '#'
			} else {
			    m.tiles[y][x] = '.'
			}
		}
	}

	return m
}

func (d *Dungeon) AddMob(m *Player) {
	d.mobs = append(d.mobs, m)
}

func (d *Dungeon) Tile(x, y int) rune {
	if x < 0 || x >= d.width || y < 0 || y >= d.height {
		return ' '
	}
	return d.tiles[y][x]
}
