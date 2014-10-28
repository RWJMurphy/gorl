package gorl

type Player struct{
	loc Coord
	c rune
}

func NewPlayer() *Player {
	p := &Player{}
	p.c = '@'
	return p
}

func (p *Player) Move(m Movement) {
	p.loc.x += m.x
	p.loc.y += m.y
}
