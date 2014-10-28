package gorl

type Coord struct {
	x, y int
}

type Dungeon struct {
	width, height int
	origin        Coord
	tiles         []rune
	mobs	      []*Player
}
const est_mob_ratio = 0.1

func NewDungeon(width, height int) *Dungeon {
	size := width*height
	m := &Dungeon{
		width, height,
		Coord{width/2, height/2},
		make([]rune, size),
		make([]*Player, 0, int(float32(size) * est_mob_ratio)),
	}
	for i, _ := range m.tiles {
		m.tiles[i] = '.'
	}
	for x := m.origin.x - 10; x < m.origin.x + 10; x++ {
		for y := m.origin.y - 10; y < m.origin.y + 10; y++ {
			m.tiles[x + y * m.width] = '#'
		}
	}
	return m
}

func (d *Dungeon) AddMob(m *Player) {
	d.mobs = append(d.mobs, m)
}

func (d *Dungeon) Tile(x, y int) rune {
	return d.tiles[y*d.width+x]
}
