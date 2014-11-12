package gorl

import (
	"fmt"
	"log"
	"strings"

	"github.com/nsf/termbox-go"
)

// A Tile represents a square in a Dungeon.
type Tile struct {
	c     rune
	color termbox.Attribute
	flags Flag
}

func (t Tile) String() string {
	return fmt.Sprintf("<Tile c:%c flags:%s>", t.c, t.flags)
}

// Flag is used to set boolean states on various objects like Tiles,
// Features, etc.
type Flag uint8

const (
	// FlagCrossable is set if the object can be moved over
	FlagCrossable Flag = 1 << iota
	// FlagLit is set if the object is lit
	FlagLit
	// FlagVisible is set if the object is visible by the Player
	FlagVisible
	// FlagSeen is set if the object has been seen by the Player
	FlagSeen
	// FlagBlocksLight is set if the object blocks light
	FlagBlocksLight
)

func (f Flag) String() string {
	var onFlags []string

	if f&FlagCrossable != 0 {
		onFlags = append(onFlags, "Crossable")
	}
	if f&FlagLit != 0 {
		onFlags = append(onFlags, "Lit")
	}
	if f&FlagVisible != 0 {
		onFlags = append(onFlags, "Visible")
	}
	if f&FlagSeen != 0 {
		onFlags = append(onFlags, "Seen")
	}
	if f&FlagBlocksLight != 0 {
		onFlags = append(onFlags, "BlocksLight")
	}

	if len(onFlags) == 0 {
		onFlags = append(onFlags, "None")
	}

	return fmt.Sprintf("<Flag %s>", strings.Join(onFlags, "|"))
}

// NewTile creates and returns a new Tile. The Tile will be rendered as c,
// in the color color, and has its flags set to flags.
func NewTile(c rune, color termbox.Attribute, flags Flag) Tile {
	t := Tile{c, color, flags}
	return t
}

// InvalidTile represents a section of the Dungeon that is out of bounds, or
// otherwise not considered "valid".
var InvalidTile = Tile{' ', termbox.ColorBlack, Flag(0) | FlagBlocksLight}

type FeatureGroup struct {
	mob     Mob
	items   []Item
	feature Feature
}

func (f FeatureGroup) String() string {
	return fmt.Sprintf("<FeatureGroup mob:%s items:%s feature:%s>",
		f.mob,
		f.items,
		f.feature,
	)
}

func (f *FeatureGroup) Each() []Feature {
	var fs []Feature
	if f.mob != nil {
		fs = append(fs, f.mob)
	}
	if f.feature != nil {
		fs = append(fs, f.feature)
	}
	if len(f.items) > 0 {
		for _, f := range f.items {
			fs = append(fs, f)
		}
	}
	return fs
}

func (f *FeatureGroup) Crossable() bool {
	for _, f := range f.Each() {
		if f.Flags()&FlagCrossable == 0 {
			return false
		}
	}
	return true
}

func (f *FeatureGroup) AddItem(i Item) {
	if f.items == nil {
		f.items = make([]Item, 0)
	}
	f.items = append(f.items, i)
}

func (f *FeatureGroup) HasItem(target Item) bool {
	for _, item := range f.items {
		if item == target {
			return true
		}
	}
	return false
}

func (f *FeatureGroup) DeleteItem(target Item) bool {
	for i, item := range f.items {
		if item == target {
			copy(f.items[i:], f.items[i+1:])
			f.items[len(f.items)-1] = nil
			f.items = f.items[:len(f.items)-1]
			return true
		}
	}
	return false
}

func (f *FeatureGroup) LightRadius() int {
	var max, rad int
	for _, feature := range f.Each() {
		rad = feature.LightRadius()
		if rad > max {
			max = rad
		}
	}
	return max
}

// Dungeon represents a level of the game.
type Dungeon struct {
	width, height int
	origin        Vector
	tiles         [][]Tile
	features      map[Vector]*FeatureGroup
	log           *log.Logger
}

// NewDungeon creates and returns a new Dungeon of the specified width and height.
//
// The dungeon's tiles are not populated.
func NewDungeon(width, height int, log *log.Logger) *Dungeon {
	size := width * height
	tiles := make([][]Tile, height)
	tilesRaw := make([]Tile, size)
	for i := range tiles {
		tiles[i], tilesRaw = tilesRaw[:width], tilesRaw[width:]
	}

	d := &Dungeon{
		width, height,
		Vector{width / 2, height / 2},
		tiles,
		make(map[Vector]*FeatureGroup),
		log,
	}
	return d
}

func (d *Dungeon) MobAt(loc Vector) Mob {
	return d.FeatureGroup(loc).mob
}

func (d *Dungeon) FeatureAt(loc Vector) Feature {
	return d.FeatureGroup(loc).feature
}

func (d *Dungeon) ItemsAt(loc Vector) []Item {
	items := d.FeatureGroup(loc).items
	itemsCopy := make([]Item, len(items))
	copy(itemsCopy, items)
	return itemsCopy
}

func (d *Dungeon) FeatureGroup(loc Vector) *FeatureGroup {
	if _, exists := d.features[loc]; !exists {
		d.features[loc] = &FeatureGroup{
			nil,
			make([]Item, 0),
			nil,
		}
	}
	return d.features[loc]
}

// AddItem adds a Item item to the Dungeon.
func (d *Dungeon) AddItem(item Item) {
	loc := item.Loc()
	d.FeatureGroup(loc).AddItem(item)
}

// DeleteItem removes Item item from the Dungeon.
func (d *Dungeon) DeleteItem(item Item) {
	loc := item.Loc()
	f := d.features[loc]
	if f == nil || !f.DeleteItem(item) {
		d.log.Panicf("Tried to delete non-existent item: %s", item)
	}
}

// AddFeature adds a Feature feature to the Dungeon.
func (d *Dungeon) AddFeature(feature Feature) {
	loc := feature.Loc()
	fg := d.FeatureGroup(loc)
	if fg.feature != nil {
		d.log.Panicf(
			"Tried to put two features on same location: %s, %s",
			feature,
			fg.feature,
		)
	}
	fg.feature = feature
}

// DeleteFeature removes Feature feature from the Dungeon.
func (d *Dungeon) DeleteFeature(feature Feature) {
	loc := feature.Loc()
	fg := d.FeatureGroup(loc)
	if fg.feature == feature {
		fg.feature = nil
	} else {
		d.log.Panicf("Tried to delete non-existent feature: %s", feature)
	}
}

// AddMob adds Mob mob to the Dungeon.
func (d *Dungeon) AddMob(mob Mob) {
	loc := mob.Loc()
	fg := d.FeatureGroup(loc)
	if fg.mob != nil {
		d.log.Panicf(
			"Tried to put two mobs on same location: %s, %s",
			mob,
			fg.mob,
		)
	}
	fg.mob = mob
}

func (d *Dungeon) ReapDead() {
	for _, m := range d.Mobs() {
		if m.Dead() {
			d.DeleteMob(m)
		}
	}
}

// DeleteMob removes Mob mob from the Dungeon.
func (d *Dungeon) DeleteMob(mob Mob) {
	loc := mob.Loc()
	fg := d.FeatureGroup(loc)
	if fg.mob == mob {
		fg.mob = nil
	} else {
		d.log.Panicf("Tried to delete non-existent mob: %s", mob)
	}
}

// MoveMob attempts to move mob in the direction move, returning true if
// successful and false otherwise.
func (d *Dungeon) MoveMob(mob Mob, move Vector) bool {
	d.log.Printf("%s moving %s", mob, move)
	dest := mob.Loc().Add(move)
	if !d.FeatureGroup(dest).Crossable() {
		return false
	}
	if !d.Tile(dest).Crossable() {
		return false
	}
	d.DeleteMob(mob)
	mob.Move(move)
	d.AddMob(mob)
	return true
}

func (d *Dungeon) Mobs() []Mob {
	var mobs []Mob
	for _, fg := range d.features {
		if fg.mob != nil {
			mobs = append(mobs, fg.mob)
		}
	}
	return mobs[:len(mobs)]
}

// CalculateLighting ranges over each Mob and Feature in the Dungeon, setting
// FlagLit on any tiles within the Feature's LightRadius that have a clear line
// sight from the Feature
func (d *Dungeon) CalculateLighting() {
	var radius int
	signal := make(chan bool)
	goroutineCount := 0

	for loc, features := range d.features {
		radius = features.LightRadius()
		if radius > 0 {
			goroutineCount++
			go func(loc Vector, radius int) {
				d.FlagByLineOfSight(loc, radius, FlagLit)
				signal <- true
			}(loc, radius)
		}
	}

	for i := 0; i < goroutineCount; i++ {
		<-signal
	}
}

// ResetFlag unsets flag on every Tile in the Dungeon
func (d *Dungeon) ResetFlag(flag Flag) {
	for x := 0; x < d.width; x++ {
		for y := 0; y < d.height; y++ {
			d.tiles[y][x].flags &= ^flag
		}
	}
}

var octantMultiplier = [4][8]int{
	{1, 0, 0, -1, -1, 0, 0, 1},
	{0, 1, -1, 0, 0, -1, 1, 0},
	{0, 1, 1, 0, 0, -1, -1, 0},
	{1, 0, 0, 1, -1, 0, 0, -1},
}

// FlagByLineOfSight uses a recusive shadowcasting algorithm to set flag on any
// Tile in the Dungeon within radius of origin if there is a clear line of sight
// from the origin to the Tile.
//
// Based on http://www.roguebasin.com/index.php?title=Ruby_shadowcasting_implementation
// which in turn is an implementation of Björn Bergström's
// http://www.roguebasin.com/index.php?title=FOV_using_recursive_shadowcasting
//
// TODO: rewrite based on something with a clear FOSS license, e.g.
// https://bitbucket.org/munificent/amaranth/src/2fc3311d903f/Amaranth.Engine/Classes/Fov.cs
func (d *Dungeon) FlagByLineOfSight(origin Vector, radius int, flag Flag) {
	d.OnTilesInLineOfSight(origin, radius, func(t *Tile, loc Vector) {
		t.flags |= flag
	})
}

type tileFunc func(*Tile, Vector)

func (d *Dungeon) OnTilesInLineOfSight(origin Vector, radius int, do tileFunc) {
	if radius == 0 {
		return
	}
	do(&d.tiles[origin.y][origin.x], origin)
	signal := make(chan bool)
	for octant := 0; octant < 8; octant++ {
		go func(origin Vector, radius int, do tileFunc, octant int) {
			d.castFlag(
				origin.x, origin.y, 1,
				1.0, 0.0,
				radius,
				octantMultiplier[0][octant],
				octantMultiplier[1][octant],
				octantMultiplier[2][octant],
				octantMultiplier[3][octant],
				do,
			)
			signal <- true
		}(origin, radius, do, octant)
	}
	for octant := 0; octant < 8; octant++ {
		<-signal
	}
}

func (d *Dungeon) castFlag(
	cx, cy, row int,
	startSlope, endSlope float64,
	radius int,
	xx, xy, yx, yy int,
	do tileFunc,
) {
	var (
		dx, dy, j                            int
		dxFloat, dyFloat                     float64
		leftSlope, rightSlope, newStartSlope float64
	)
	if startSlope < endSlope {
		return
	}
	radiusSquared := radius * radius
	for j = row; j <= radius; j++ {
		dx, dy = -j-1, -j
		blocked := false
		for dx <= 0 {
			dx++
			// Translate the dx, dy coordinates into map coordinates
			mapX, mapY := cx+dx*xx+dy*xy, cy+dx*yx+dy*yy
			if mapX < 0 || mapX >= d.width || mapY < 0 || mapY >= d.height {
				continue
			}
			// leftSlope and rightSlope store the slopes of the left and
			// right extremeties of the square we're considering
			dxFloat, dyFloat = float64(dx), float64(dy)
			leftSlope, rightSlope = (dxFloat-0.5)/(dyFloat+0.5), (dxFloat+0.5)/(dyFloat-0.5)
			if startSlope < rightSlope {
				continue
			} else if endSlope > leftSlope {
				break
			} else {
				t := d.Tile(Vector{mapX, mapY})
				// our light beam is touching this square; flag it:
				if dx*dx+dy*dy < radiusSquared {
					do(t, Vector{mapX, mapY})
				}
				if blocked {
					// we're scanning a row of blocked squares
					if t.BlocksLight() {
						newStartSlope = rightSlope
						continue
					} else {
						blocked = false
						startSlope = newStartSlope
					}
				} else {
					if t.BlocksLight() && j < radius {
						// this is a blocking square, start a child scan:
						blocked = true
						d.castFlag(cx, cy, j+1, startSlope, leftSlope, radius, xx, xy, yx, yy, do)
						newStartSlope = rightSlope
					}
				}
			}
		}
		if blocked {
			break
		}
	}
}

// Tile fetches the Dungeon Tile at (x, y)
func (d *Dungeon) Tile(loc Vector) *Tile {
	if loc.x < 0 || loc.x >= d.width || loc.y < 0 || loc.y >= d.height {
		t := InvalidTile
		return &t
	}
	return &d.tiles[loc.y][loc.x]
}

// Crossable returns true if the Tile can be moved across
func (t *Tile) Crossable() bool {
	return t.flags&FlagCrossable != 0
}

// Lit returns true if the Tile is lit by a light source
func (t *Tile) Lit() bool {
	return t.flags&FlagLit != 0
}

// Visible returns true if the Tile is within the Player's FOV
func (t *Tile) Visible() bool {
	return t.flags&FlagVisible != 0
}

// Seen returns true if the Tile has ever been seen by the Player
func (t *Tile) Seen() bool {
	return t.flags&FlagSeen != 0
}

// BlocksLight returns true if the Tile does not allow light to pass through
func (t *Tile) BlocksLight() bool {
	return t.flags&FlagBlocksLight != 0
}
