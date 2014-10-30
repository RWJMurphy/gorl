package gorl

import (
	"github.com/nsf/termbox-go"
	"math/rand"
)

const (
	DefaultDungeonWidth  = 256
	DefaultDungeonHeight = 256
)

type Game struct {
	ui             *UI
	messages       []string
	player         *Player
	dungeons       []*Dungeon
	currentDungeon *Dungeon
}

func NewGame() (*Game, error) {
	game := &Game{}
	game.messages = make([]string, 0, 10)

	dungeon := NewDungeon(DefaultDungeonWidth, DefaultDungeonHeight)
	game.dungeons = make([]*Dungeon, 0, 10)
	game.dungeons = append(game.dungeons, dungeon)

	game.player = NewPlayer()
	game.player.loc = dungeon.origin
	dungeon.AddMob(game.player)

	for i := 0; i < 100; i++ {
		x, y := rand.Int()%dungeon.width, rand.Int()%dungeon.height
		for !dungeon.Tile(x, y).Crossable() {
			x, y = rand.Int()%dungeon.width, rand.Int()%dungeon.height
		}
		mob := NewMob("orc", 'o')
		mob.color = termbox.ColorGreen
		mob.loc = Coord{x, y}
		dungeon.AddMob(mob)
	}

	for i := 0; i < 100; i ++ {
		x, y := rand.Int()%dungeon.width, rand.Int()%dungeon.height
		for !dungeon.Tile(x, y).Crossable() {
			x, y = rand.Int()%dungeon.width, rand.Int()%dungeon.height
		}
		feature := NewFeature("torch", '!')
		feature.loc = Coord{x, y}
		feature.color = termbox.ColorRed | termbox.AttrBold
		feature.lightRadius = 20
		dungeon.AddFeature(feature)
	}

	dungeon.ResetFlag(FlagLit | FlagVisible)
	dungeon.CalculateLighting()
	dungeon.FlagByLineOfSight(game.player.loc, game.player.visionRadius, FlagVisible)

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

type Movement struct {
	x, y int
}

func (c Coord) Plus(m Movement) Coord {
	c.x += m.x
	c.y += m.y
	return c
}

func (game *Game) Move(movement Movement) {
	dest := game.player.loc.Plus(movement)
	// ASSUMPTION: only one mob per tile
	_, blocked := game.currentDungeon.mobs[dest]
	if blocked {
		return
	}
	if !game.currentDungeon.Tile(dest.x, dest.y).Crossable() {
		return
	}
	game.currentDungeon.MoveMob(game.player, movement)

	game.currentDungeon.ResetFlag(FlagLit | FlagVisible)
	game.currentDungeon.CalculateLighting()
	game.currentDungeon.FlagByLineOfSight(game.player.loc, game.player.visionRadius, FlagVisible)

	game.ui.PointCameraAt(game.currentDungeon, game.player.loc)
	game.ui.dirty = true
}

func (game *Game) AddMessage(message string) {
	game.messages = append(game.messages, message)
	message_count := game.ui.messageWidget.height - 2
	if message_count > len(game.messages) {
		message_count = len(game.messages)
	}
	game.ui.messages = game.messages[len(game.messages)-message_count:]
	game.ui.dirty = true
}

func (game *Game) Run() {
	game.MainLoop()
}

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

func (game *Game) Close() {
	game.ui.Close()
}

func (game *Game) SetDungeon(d *Dungeon) {
	game.currentDungeon = d
	game.ui.PointCameraAt(d, d.origin)
}
