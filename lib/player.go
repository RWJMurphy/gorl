package gorl

type Player struct {
	mob
}

const (
	PlayerVisionRadius = 100
	PlayerLightRadius  = 2
)

func NewPlayer() *Player {
	p := &Player{
		*NewMob("Player", '@'),
	}
	p.visionRadius = PlayerVisionRadius
	p.lightRadius = PlayerLightRadius
	return p
}
