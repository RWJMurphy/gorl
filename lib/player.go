package gorl

import "log"

// Player represents the player -- it's basically a special form of Mob
type Player interface {
	Mob
}

type player struct {
	mob
}

const (
	// PlayerVisionRadius sets how far the player can see
	PlayerVisionRadius = 100
	// PlayerLightRadius sets how far the player emits light
	PlayerLightRadius = 0
)

// NewPlayer creates and returns a new Player
func NewPlayer(log log.Logger, dungeon *Dungeon) Player {
	p := &player{
		*NewMob("Player", '@', log, dungeon).(*mob),
	}
	p.mob.visionRadius = PlayerVisionRadius
	p.mob.lightRadius = PlayerLightRadius
	return p
}

func (p *player) Tick(turn uint) bool {
	p.lastTicked = turn
	return false
}
