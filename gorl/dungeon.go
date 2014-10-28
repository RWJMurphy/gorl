package gorl

type Coord struct {
	x, y int
}

type Dungeon struct {
	Width, Height int
	Origin        Coord
	tiles         []rune
}

func NewDungeon(w, h int) *Dungeon {
	size := w*h
	m := &Dungeon{
		w, h,
		Coord{w/2, h/2},
		make([]rune, size),
	}
	for i, _ := range m.tiles {
		m.tiles[i] = '.'
	}
	m.tiles[m.Origin.x + m.Origin.y * m.Width] = 'O'
	return m
}

func (d Dungeon) Tile(x, y int) rune {
	return d.tiles[y*d.Width+x]
}
