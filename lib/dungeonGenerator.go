package gorl

import (
	"log"
	"math/rand"

	"github.com/nsf/termbox-go"
)

type dungeonRoom struct {
	width, height int
	portals       []Vector
	tiles         [][]Tile
}

func newDungeonRoom(width, height int, dice *rand.Rand) *dungeonRoom {
	tiles := make([][]Tile, height)
	tilesRaw := make([]Tile, width*height)
	for i := range tilesRaw {
		tilesRaw[i] = NewTile('.', termbox.ColorWhite, Flag(0)|FlagCrossable)
	}
	for i := range tiles {
		tiles[i], tilesRaw = tilesRaw[:width], tilesRaw[width:]
	}

	var edgeTiles []Vector
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x == 0 || x == width-1 || y == 0 || y == height-1 {
				edgeTiles = append(edgeTiles, Vector{x, y})
			}
		}
	}

	for _, loc := range edgeTiles {
		tiles[loc.y][loc.x] = NewTile('#', termbox.ColorYellow, Flag(0)|FlagBlocksLight)
	}

	var portals []Vector
	for i := 1; i < 1+dice.Intn(4); i++ {
		portals = append(portals, edgeTiles[dice.Intn(len(edgeTiles))])
	}
	room := dungeonRoom{
		width,
		height,
		portals,
		tiles,
	}
	return &room
}

func (d *Dungeon) paintRoom(room *dungeonRoom, topLeft Vector) []Vector {
	// TODO: support different orientations by changing order tiles iterated
	for x := 0; x < room.width; x++ {
		for y := 0; y < room.height; y++ {
			dungeonLoc := topLeft.Add(Vector{x, y})
			d.tiles[dungeonLoc.y][dungeonLoc.x] = room.tiles[y][x]
		}
	}
	newPortals := make([]Vector, len(room.portals))
	for i, portal := range room.portals {
		portalLoc := portal.Add(topLeft)
		newPortals[i] = portalLoc
	}
	return newPortals
}

func GenerateDungeon(log *log.Logger, dice *rand.Rand) *Dungeon {
	width, height := 100, 100
	d := NewDungeon(width, height, log)
	var tile Tile

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if dice.Float32() < 0.05 {
				tile = NewTile('#', termbox.ColorYellow, Flag(0)|FlagBlocksLight)
			} else {
				tile = NewTile('.', termbox.ColorWhite, Flag(0)|FlagCrossable)
			}
			d.tiles[y][x] = tile
		}
	}

	roomCount := dice.Intn(10) + 10
	var portals []Vector
	for i := 0; i < roomCount; i++ {
		roomWidth, roomHeight := dice.Intn(10)+10, dice.Intn(10)+10
		topLeft := Vector{
			dice.Intn(width - roomWidth),
			dice.Intn(height - roomHeight),
		}

		room := newDungeonRoom(roomWidth, roomHeight, dice)
		roomPortals := d.paintRoom(room, topLeft)
		for _, p := range roomPortals {
			portals = append(portals, p)
		}
	}

	log.Printf("Room portals: %s", portals)

	for _, portalLoc := range portals {
		d.tiles[portalLoc.y][portalLoc.x] = NewTile('+', termbox.ColorWhite, Flag(0)|FlagCrossable|FlagBlocksLight)
	}
	return d
}
