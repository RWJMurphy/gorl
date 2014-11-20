package gorl

import (
	"log"
	"math/rand"
)

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
	PlayerLightRadius = 2
)

// NewPlayer creates and returns a new Player
func NewPlayer(log *log.Logger, dungeon *Dungeon) Player {
	p := &player{
		*NewMob("Player", '@', log, dungeon).(*mob),
	}
	p.mob.visionRadius = PlayerVisionRadius
	p.mob.lightRadius = PlayerLightRadius
	return p
}

// XXX Should I move input handling here? Or does that couple the Player
// to the UI too closely? Argh.
func (p *player) Tick(turn uint, dice *rand.Rand) MobAction {
	p.lastTicked = turn
	return MobAction{ActNone, nil}
}
