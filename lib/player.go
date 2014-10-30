package gorl

type Player struct {
	mob
}

const (
	PlayerVisionRadius = 100
	PlayerLightRadius  = 10
)

func NewPlayer() *Player {
	p := &Player{
		*NewMob('@'),
	}
	p.visionRadius = PlayerVisionRadius
	p.lightRadius = PlayerLightRadius
	return p
}
