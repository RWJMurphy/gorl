package gorl

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

// Default size for new dungeons
const (
	DefaultDungeonWidth  = 512
	DefaultDungeonHeight = 512
)

// GameState represents the state of the Game engine
type GameState int

const (
	// GameInvalidState represents any bad state
	GameInvalidState GameState = iota
	// GamePlayerTurn is when waiting on the player to make an action
	GamePlayerTurn
	// GameWorldTurn is when the AI and world objects get to act
	GameWorldTurn
	// GameClosed is when the game is done, and will shut down
	GameClosed
)

func (s GameState) String() string {
	switch s {
	case GameInvalidState:
		return "GameInvalidState"
	case GamePlayerTurn:
		return "GamePlayerTurn"
	case GameWorldTurn:
		return "GameWorldTurn"
	case GameClosed:
		return "GameClosed"
	default:
		return fmt.Sprintf("GameState(%d)", s)
	}
}

// Game is the entry type to GoRL. Manages the UI, dungeons, player, etc.
type Game struct {
	ui             *UI
	messages       []string
	player         Player
	dungeons       []*Dungeon
	currentDungeon *Dungeon
	state          GameState
	turn           uint
	log            log.Logger
}

// NewGame initializes and returns a new Game. Or an error. You should check that.
// Please `defer game.Close()`.
func NewGame(log log.Logger) (*Game, error) {
	game := &Game{}
	game.log = log
	game.messages = make([]string, 0, 10)
	game.turn = 0

	dungeon := NewDungeon(DefaultDungeonWidth, DefaultDungeonHeight, log)
	game.dungeons = make([]*Dungeon, 0, 10)
	game.dungeons = append(game.dungeons, dungeon)

	game.player = NewPlayer(game.log, dungeon)
	game.player.SetLoc(dungeon.origin)

	torch := NewItem("bright torch", '!', 1)
	torch.SetLightRadius(20)
	game.player.AddToInventory(torch)

	dungeon.AddMob(game.player)

	for i := 0; i < 100; i++ {
		x, y := rand.Intn(dungeon.width), rand.Intn(dungeon.height)
		_, mobExists := dungeon.mobs[Coord{x, y}]
		for !dungeon.Tile(x, y).Crossable() || mobExists {
			x, y = rand.Intn(dungeon.width), rand.Intn(dungeon.height)
			_, mobExists = dungeon.mobs[Coord{x, y}]
		}
		mob := NewMob(fmt.Sprintf("orc #%d", i), 'o', game.log, dungeon)
		mob.SetColor(termbox.ColorGreen)
		mob.SetLoc(Coord{x, y})

		torch := NewItem("torch", '!', 1)
		torch.SetLightRadius(10)
		mob.AddToInventory(torch)

		dungeon.AddMob(mob)
	}

	dungeon.ResetFlag(FlagLit | FlagVisible)
	dungeon.CalculateLighting()
	dungeon.OnTilesInLineOfSight(game.player.Loc(), game.player.VisionRadius(), func(t *Tile) {
		if t.Lit() {
			t.flags |= FlagVisible | FlagSeen
		}
	})

	ui, err := NewUI(game)
	if err != nil {
		return nil, err
	}
	game.ui = ui
	game.SetDungeon(dungeon)

	game.AddMessage("Welcome to GoRL!")
	game.state = GamePlayerTurn
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
	if moved := game.currentDungeon.MoveMob(game.player, movement); !moved {
		return moved
	}

	game.currentDungeon.ResetFlag(FlagLit | FlagVisible)
	game.currentDungeon.CalculateLighting()
	game.currentDungeon.OnTilesInLineOfSight(game.player.Loc(), game.player.VisionRadius(), func(t *Tile) {
		if t.Lit() {
			t.flags |= FlagVisible | FlagSeen
		}
	})

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
	game.log.Println("Entering main loop")
mainLoop:
	for {
		game.log.Printf("game.state = %s", game.state)
		switch state := game.state; state {
		case GameWorldTurn:
			game.WorldTick()
			game.state = GamePlayerTurn
		case GamePlayerTurn:
			game.state = game.ui.WaitAndHandleInput()
		case GameClosed:
			break mainLoop
		case GameInvalidState:
			fallthrough
		default:
			log.Panicf("Bad game state: %s", game.state)
		}
		game.ui.Paint()
	}
}

// WorldTick runs a single turn of the game engine
func (game *Game) WorldTick() {
	tickStartTime := time.Now()
	game.turn++
	game.log.Printf("Game tick: %d", game.turn)
	game.log.Printf("There are %d mobs", len(game.currentDungeon.mobs))
	changed := false
	for _, m := range game.currentDungeon.Mobs() {
		changed = m.Tick(game.turn) || changed
	}
	if changed {
		game.ui.dirty = true
	}
	game.currentDungeon.ResetFlag(FlagLit)
	game.currentDungeon.CalculateLighting()
	tickRunTime := time.Now().Sub(tickStartTime)
	game.log.Printf("Tick took %v to run", tickRunTime)
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
