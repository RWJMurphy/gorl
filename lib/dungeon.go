package gorl

import (
	"fmt"
	"math/rand"
	"strings"
	"github.com/nsf/termbox-go"
)

type Coord struct {
	x, y int
}

func (c Coord) String() string {
	return fmt.Sprintf("<Coord x:%d, y:%d>", c.x, c.y)
}

type Tile struct {
	c     rune
	color termbox.Attribute
	flags Flag
}

func (t Tile) String() string {
	return fmt.Sprintf("<Tile c:%c flags:%s>", t.c, t.flags)
}

type Flag uint8

const (
	FlagCrossable Flag = 1 << iota
	FlagLit
	FlagVisible
	FlagBlocksLight
)

func (f Flag) String() string {
	on_flags := make([]string, 0)

	if f&FlagCrossable != 0 {
		on_flags = append(on_flags, "Crossable")
	}
	if f&FlagLit != 0 {
		on_flags = append(on_flags, "Lit")
	}
	if f&FlagVisible != 0 {
		on_flags = append(on_flags, "Visible")
	}
	if f&FlagBlocksLight != 0 {
		on_flags = append(on_flags, "BlocksLight")
	}

	if len(on_flags) == 0 {
		on_flags = append(on_flags, "None")
	}

	return fmt.Sprintf("<Flag %s>", strings.Join(on_flags, "|"))
}

func NewTile(c rune, color termbox.Attribute, flags Flag) Tile {
	t := Tile{c, color, flags}
	return t
}

var InvalidTile = Tile{' ', termbox.ColorBlack, Flag(0) | FlagBlocksLight}

type Dungeon struct {
	width, height int
	origin        Coord
	tiles         [][]Tile
	mobs          map[Coord]Mob
	features      map[Coord]Feature
}

func NewDungeon(width, height int) *Dungeon {
	size := width * height
	tiles := make([][]Tile, height)
	tiles_raw := make([]Tile, size)
	for i := range tiles {
		tiles[i], tiles_raw = tiles_raw[:width], tiles_raw[width:]
	}

	m := &Dungeon{
		width, height,
		Coord{width / 2, height / 2},
		tiles,
		make(map[Coord]Mob),
		make(map[Coord]Feature),
	}

	var tile Tile
	wallChance := 0.1
	for x := 0; x < width; x++ {
		for y := 0; y < width; y++ {
			if rand.Float64() <= wallChance {
				tile = NewTile('#', termbox.ColorYellow, Flag(0)|FlagBlocksLight)
			} else {
				tile = NewTile('.', termbox.ColorWhite, Flag(0)|FlagCrossable)
			}
			m.tiles[y][x] = tile
		}
	}

	return m
}

func (d *Dungeon) AddFeature(feature Feature) {
	if other_feature, exists := d.features[feature.Loc()]; exists {
		panic(fmt.Sprintf(
			"Tried to put two features on same location: %s, %s",
			feature,
			other_feature,
		))
	}
	d.features[feature.Loc()] = feature
}

func (d *Dungeon) AddMob(mob Mob) {
	if other_mob, exists := d.mobs[mob.Loc()]; exists {
		panic(fmt.Sprintf(
			"Tried to put two mobs on same location: %s, %s",
			mob,
			other_mob,
		))
	}
	d.mobs[mob.Loc()] = mob
}

func (d *Dungeon) DeleteMob(mob Mob) {
	if _, exists := d.mobs[mob.Loc()]; exists {
		delete(d.mobs, mob.Loc())
	} else {
		panic(fmt.Sprintf("Tried to delete non-existent mob: %s", mob))
	}
}

func (d *Dungeon) MoveMob(mob Mob, move Movement) {
	d.DeleteMob(mob)
	mob.Move(move)
	d.AddMob(mob)
}

func (d *Dungeon) CalculateLighting() {
	for loc, mob := range d.mobs {
		d.FlagByLineOfSight(loc, mob.LightRadius(), FlagLit)
	}
	for loc, feature := range d.features {
		d.FlagByLineOfSight(loc, feature.LightRadius(), FlagLit)
	}
}

var octant_multiplier = [4][8]int{
	{1, 0, 0, -1, -1, 0, 0, 1},
	{0, 1, -1, 0, 0, -1, 1, 0},
	{0, 1, 1, 0, 0, -1, -1, 0},
	{1, 0, 0, 1, -1, 0, 0, -1},
}

func (d *Dungeon) ResetFlag(flag Flag) {
	for x := 0; x < d.width; x++ {
		for y := 0; y < d.width; y++ {
			d.tiles[y][x].flags &= ^flag
		}
	}
}

// Recursive shadow casting
// Based on http://www.roguebasin.com/index.php?title=Ruby_shadowcasting_implementation
// which in turn is an implementation of Björn Bergström's
// http://www.roguebasin.com/index.php?title=FOV_using_recursive_shadowcasting
//
// TODO: rewrite based on something with a clear FOSS license, e.g.
// https://bitbucket.org/munificent/amaranth/src/2fc3311d903f/Amaranth.Engine/Classes/Fov.cs
func (d *Dungeon) FlagByLineOfSight(origin Coord, radius int, flag Flag) {
	if radius == 0 {
		return
	}
	d.tiles[origin.y][origin.x].flags |= flag
	for octant := 0; octant < 8; octant++ {
		d.castFlag(
			origin.x, origin.y, 1,
			1.0, 0.0,
			radius,
			octant_multiplier[0][octant],
			octant_multiplier[1][octant],
			octant_multiplier[2][octant],
			octant_multiplier[3][octant],
			flag,
		)
	}
}

func (d *Dungeon) castFlag(
	cx, cy, row int,
	start_slope, end_slope float64,
	radius int,
	xx, xy, yx, yy int,
	flag Flag,
) {
	var (
		dx, dy, j                         int
		dx_f, dy_f                        float64
		l_slope, r_slope, new_start_slope float64
	)
	if start_slope < end_slope {
		return
	}
	radius_2 := radius * radius
	for j = row; j <= radius; j++ {
		dx, dy = -j-1, -j
		blocked := false
		for dx <= 0 {
			dx += 1
			// Translate the dx, dy coordinates into map coordinates
			mx, my := cx+dx*xx+dy*xy, cy+dx*yx+dy*yy
			if mx < 0 || mx >= d.width || my < 0 || my >= d.height {
				continue
			}
			// l_slope and r_slope store the slopes of the left and
			// right extremeties of the square we're considering
			dx_f, dy_f = float64(dx), float64(dy)
			l_slope, r_slope = (dx_f-0.5)/(dy_f+0.5), (dx_f+0.5)/(dy_f-0.5)
			if start_slope < r_slope {
				continue
			} else if end_slope > l_slope {
				break
			} else {
				t := d.Tile(mx, my)
				// our light beam is touching this square; flag it:
				if dx*dx+dy*dy < radius_2 {
					t.flags |= flag
				}
				if blocked {
					// we're scanning a row of blocked squares
					if t.BlocksLight() {
						new_start_slope = r_slope
						continue
					} else {
						blocked = false
						start_slope = new_start_slope
					}
				} else {
					if t.BlocksLight() && j < radius {
						// this is a blocking square, start a child scan:
						blocked = true
						d.castFlag(cx, cy, j+1, start_slope, l_slope, radius, xx, xy, yx, yy, flag)
						new_start_slope = r_slope
					}
				}
			}
		}
		if blocked {
			break
		}
	}
}

func (d *Dungeon) Tile(x, y int) *Tile {
	if x < 0 || x >= d.width || y < 0 || y >= d.height {
		t := InvalidTile
		return &t
	}
	return &d.tiles[y][x]
}

func (t *Tile) Crossable() bool {
	return t.flags&FlagCrossable != 0
}

func (t *Tile) Lit() bool {
	return t.flags&FlagLit != 0
}

func (t *Tile) Visible() bool {
	return t.flags&FlagVisible != 0
}

func (t *Tile) BlocksLight() bool {
	return t.flags&FlagBlocksLight != 0
}
