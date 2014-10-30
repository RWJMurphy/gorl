package gorl

import (
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
	game := new(Game)
	game.messages = make([]string, 0, 10)

	dungeon := NewDungeon(DefaultDungeonWidth, DefaultDungeonHeight)
	game.dungeons = make([]*Dungeon, 0, 10)
	game.dungeons = append(game.dungeons, dungeon)

	game.player = NewPlayer()
	game.player.loc = dungeon.origin
	dungeon.AddMob(game.player)

	for i := 0; i < 100; i++ {
		x, y := rand.Int() % dungeon.width, rand.Int() % dungeon.height
		for dungeon.Tile(x, y).flags&FlagCrossable == 0 {
			x, y = rand.Int() % dungeon.width, rand.Int() % dungeon.height
		}
		mob := NewMob('o')
		mob.loc.x = x
		mob.loc.y = y
		dungeon.AddMob(mob)
	}

	dungeon.CalculateLighting()

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
	other_mob, blocked := game.currentDungeon.mobs[dest]
	if blocked && other_mob.Flags()&FlagCrossable == 0 {
		return
	}
	dest_tile := game.currentDungeon.Tile(dest.x, dest.y)
	if dest_tile.flags&FlagCrossable == 0 {
		return
	}
	game.currentDungeon.MoveMob(game.player, movement)
	game.currentDungeon.CalculateLighting()
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
