package gorl

import (
	"math/rand"

	"github.com/nsf/termbox-go"
)

// Default size for new dungeons
const (
	DefaultDungeonWidth  = 256
	DefaultDungeonHeight = 256
)

// Game is the entry type to GoRL. Manages the UI, dungeons, player, etc.
type Game struct {
	ui             *UI
	messages       []string
	player         Player
	dungeons       []*Dungeon
	currentDungeon *Dungeon
}

// NewGame initializes and returns a new Game. Or an error. You should check that.
// Please `defer game.Close()`.
func NewGame() (*Game, error) {
	game := &Game{}
	game.messages = make([]string, 0, 10)

	dungeon := NewDungeon(DefaultDungeonWidth, DefaultDungeonHeight)
	game.dungeons = make([]*Dungeon, 0, 10)
	game.dungeons = append(game.dungeons, dungeon)

	game.player = NewPlayer()
	game.player.SetLoc(dungeon.origin)
	dungeon.AddMob(game.player)

	for i := 0; i < 100; i++ {
		x, y := rand.Int()%dungeon.width, rand.Int()%dungeon.height
		for !dungeon.Tile(x, y).Crossable() {
			x, y = rand.Int()%dungeon.width, rand.Int()%dungeon.height
		}
		mob := NewMob("orc", 'o')
		mob.SetColor(termbox.ColorGreen)
		mob.SetLoc(Coord{x, y})
		dungeon.AddMob(mob)
	}

	for i := 0; i < 100; i++ {
		x, y := rand.Int()%dungeon.width, rand.Int()%dungeon.height
		for !dungeon.Tile(x, y).Crossable() {
			x, y = rand.Int()%dungeon.width, rand.Int()%dungeon.height
		}
		feature := NewFeature("torch", '!')
		feature.SetLoc(Coord{x, y})
		feature.SetColor(termbox.ColorRed | termbox.AttrBold)
		feature.SetLightRadius(20)
		dungeon.AddFeature(feature)
	}

	dungeon.ResetFlag(FlagLit | FlagVisible)
	dungeon.CalculateLighting()
	dungeon.FlagByLineOfSight(game.player.Loc(), game.player.VisionRadius(), FlagVisible)

	ui, err := NewUI()
	if err != nil {
		return nil, err
	}
	game.ui = ui
	ui.game = game
	game.SetDungeon(dungeon)

	game.AddMessage("Welcome to GoRL!")
	return game, nil
}

// Movement represents a change in location.
type Movement struct {
	x, y int
}

// Plus adds a Movement to a Coord, and returns a new Coord
func (c Coord) Plus(m Movement) Coord {
	c.x += m.x
	c.y += m.y
	return c
}

// Move tries to move the player in the direction `movement`, and returns a bool
// indicating whether it was successful.
func (game *Game) Move(movement Movement) bool {
	dest := game.player.Loc().Plus(movement)
	// ASSUMPTION: only one mob per tile
	_, blocked := game.currentDungeon.mobs[dest]
	if blocked {
		return false
	}
	if !game.currentDungeon.Tile(dest.x, dest.y).Crossable() {
		return false
	}
	game.currentDungeon.MoveMob(game.player, movement)

	game.currentDungeon.ResetFlag(FlagLit | FlagVisible)
	game.currentDungeon.CalculateLighting()
	game.currentDungeon.FlagByLineOfSight(game.player.Loc(), game.player.VisionRadius(), FlagVisible)

	game.ui.PointCameraAt(game.currentDungeon, game.player.Loc())
	game.ui.dirty = true
	return true
}

// AddMessage adds a message to the UI's message buffer for display in the MessageLogWidget
func (game *Game) AddMessage(message string) {
	game.messages = append(game.messages, message)
	messageCount := game.ui.messageWidget.height - 2
	if messageCount > len(game.messages) {
		messageCount = len(game.messages)
	}
	game.ui.messages = game.messages[len(game.messages)-messageCount:]
	game.ui.dirty = true
}

// Run runs the Game.
func (game *Game) Run() {
	game.MainLoop()
}

// MainLoop is the Game's main loop. Ticks the UI until it closes.
func (game *Game) MainLoop() {
	game.ui.Paint()
mainLoop:
	for {
		game.ui.Tick()
		if game.ui.State == StateClosed {
			break mainLoop
		}
	}
}

// Close cleans up after a Game.
func (game *Game) Close() {
	game.ui.Close()
}

// SetDungeon sets the current dungeon to d, and points the UI's camera at its
// origin.
func (game *Game) SetDungeon(d *Dungeon) {
	game.currentDungeon = d
	game.ui.PointCameraAt(d, d.origin)
}
