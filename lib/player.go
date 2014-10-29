package gorl

type Player struct {
	mob
}

func NewPlayer() *Player {
	p := &Player{
		*NewMob('@'),
	}
	return p
}
