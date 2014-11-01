package gorl

import (
	"log"
	"math/rand"

	"github.com/nsf/termbox-go"
)

func GenerateDungeon(width, height int, log log.Logger) *Dungeon {
	d := NewDungeon(width, height, log)
	var tile Tile
	wallChance := 0.02
	for x := 0; x < width; x++ {
		for y := 0; y < width; y++ {
			if rand.Float64() <= wallChance {
				tile = NewTile('#', termbox.ColorYellow, Flag(0)|FlagBlocksLight)
			} else {
				tile = NewTile('.', termbox.ColorWhite, Flag(0)|FlagCrossable)
			}
			d.tiles[y][x] = tile
		}
	}
	return d
}
