package gorl

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
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

type playerAction uint

const (
	ActNone playerAction = iota
	ActWait
	ActMove // target is a Vec to move
	ActDrop // target is an Item to drop
	ActDropAll
	ActPickUpAll
)

type PlayerAction struct {
	action playerAction
	target interface{}
}

func (a playerAction) String() string {
	switch a {
	case ActNone:
		return "ActNone"
	case ActWait:
		return "ActWait"
	case ActDropAll:
		return "ActDropAll"
	case ActPickUpAll:
		return "ActPickUpAll"
	default:
		return fmt.Sprintf("playerAction(%d)", a)
	}
}

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
	ui             UI
	messages       []string
	player         Player
	dungeons       []*Dungeon
	currentDungeon *Dungeon
	state          GameState
	turn           uint
	log            *log.Logger
}

// NewGame initializes and returns a new Game. Or an error. You should check that.
// Please `defer game.Close()`.
func NewGame(log *log.Logger) (*Game, error) {
	game := &Game{}
	game.log = log
	game.messages = make([]string, 0, 10)
	game.turn = 0

	dungeon := GenerateDungeon(log)
	game.dungeons = make([]*Dungeon, 0, 10)
	game.dungeons = append(game.dungeons, dungeon)

	game.player = NewPlayer(game.log, dungeon)
	game.player.SetLoc(dungeon.origin)

	torch := NewItem("bright torch", '!', 1)
	torch.SetLightRadius(20)
	game.player.AddToInventory(torch)

	dungeon.AddMob(game.player)

	for i := 0; i < 10; i++ {
		x, y := rand.Intn(dungeon.width), rand.Intn(dungeon.height)
		for !(dungeon.Tile(x, y).Crossable() && dungeon.FeatureGroup(Vec{x, y}).Crossable()) {
			x, y = rand.Intn(dungeon.width), rand.Intn(dungeon.height)
		}
		mob := NewMob(fmt.Sprintf("orc #%d", i), 'o', game.log, dungeon)
		mob.SetColor(termbox.ColorGreen)
		mob.SetLoc(Vec{x, y})

		torch := NewItem("torch", '!', 1)
		torch.SetLightRadius(10)
		mob.AddToInventory(torch)

		dungeon.AddMob(mob)
	}

	dungeon.ResetFlag(FlagLit | FlagVisible)
	dungeon.CalculateLighting()
	dungeon.OnTilesInLineOfSight(game.player.Loc(), game.player.VisionRadius(), func(t *Tile, loc Vec) {
		if t.Lit() {
			t.flags |= FlagVisible | FlagSeen
		}
	})

	ui, err := NewTermboxUI(game)
	if err != nil {
		return nil, err
	}
	game.ui = ui
	game.SetDungeon(dungeon)

	game.AddMessage("Welcome to GoRL!")
	game.state = GamePlayerTurn
	return game, nil
}

// MoveOrAct calculates the destination tile based on the movement parameter and
// the Player's location, and then
//   * if there is a mob on the destination, attacks the mob and returns true
//   * if not and destination is Crossable, moves the player there and returns true
//   * if the destination is not Crossable, returns false
func (game *Game) MoveOrAct(movement Vec) bool {
	destination := game.player.Loc().Add(movement)
	if mob := game.currentDungeon.MobAt(destination); mob != nil {
		if damageDealt, ok := game.player.Attack(mob); ok {
			game.AddMessage(fmt.Sprintf("Hit %s for %d damage", mob.Name(), damageDealt))
			if mob.Dead() {
				game.AddMessage(fmt.Sprintf("The %s dies!", mob.Name()))
			}
			return true
		}
		return false
	} else if moved := game.currentDungeon.MoveMob(game.player, movement); !moved {
		return moved
	}

	game.currentDungeon.ResetFlag(FlagLit | FlagVisible)
	game.currentDungeon.CalculateLighting()
	game.currentDungeon.OnTilesInLineOfSight(game.player.Loc(), game.player.VisionRadius(), func(t *Tile, loc Vec) {
		if t.Lit() {
			t.flags |= FlagVisible | FlagSeen
		}
	})

	game.ui.PointCameraAt(game.currentDungeon, game.player.Loc())
	game.ui.MarkDirty()
	return true
}

// AddMessage adds a message to the UI's message buffer for display in the MessageLogWidget
func (game *Game) AddMessage(message string) {
	message = fmt.Sprintf("%d: %s", game.turn, message)
	game.log.Println(message)
	game.messages = append(game.messages, message)
	messageCount := game.ui.MessagesWanted()
	if messageCount > len(game.messages) {
		messageCount = len(game.messages)
	}
	game.ui.SetMessages(game.messages[len(game.messages)-messageCount:])
	game.ui.MarkDirty()
}

// Run runs the Game.
func (game *Game) Run() {
	game.MainLoop()
}

// MainLoop is the Game's main loop. Ticks the UI until it closes.
func (game *Game) MainLoop() {
	var action PlayerAction
	nextState := game.state

	game.ui.Paint()
	game.log.Println("Entering main loop")
mainLoop:
	for {
		switch game.state {
		case GameWorldTurn:
			game.WorldTick()
			nextState = GamePlayerTurn
		case GamePlayerTurn:
			action, nextState = game.ui.DoEvent()

			switch action.action {
			case ActWait:
				nextState = GameWorldTurn
			case ActDrop:
				item := action.target.(Item)
				game.player.DropItem(item, game.currentDungeon)
				game.AddMessage(fmt.Sprintf("Dropped %s", item.Name()))
			case ActDropAll:
				for _, item := range game.player.Inventory() {
					game.player.DropItem(item, game.currentDungeon)
					game.AddMessage(fmt.Sprintf("Dropped %s", item.Name()))
				}
			case ActPickUpAll:
				items := game.currentDungeon.ItemsAt(game.player.Loc())
				game.log.Printf("Attempting to pick up items: %s", items)
				if len(items) > 0 {
					for _, item := range items {
						game.log.Printf("Player picking up %s", item)
						game.currentDungeon.DeleteItem(item)
						game.player.AddToInventory(item)
						game.AddMessage(fmt.Sprintf("Picked up %s", item.Name()))
					}
				} else {
					game.AddMessage("Nothing to pick up.")
					nextState = game.state
				}
			case ActMove:
				direction := action.target.(Vec)
				if !game.MoveOrAct(direction) {
					nextState = GamePlayerTurn
				}
			case ActNone:
			default:
				log.Panicf("Bad action: %s", action)
			}
		case GameClosed:
			break mainLoop
		case GameInvalidState:
			fallthrough
		default:
			log.Panicf("Bad game state: %s", game.state)
		}
		game.currentDungeon.ReapDead()
		game.ui.Paint()
		game.log.Printf("Game state change: %s -> %s", game.state, nextState)
		game.state = nextState
	}
}

// WorldTick runs a single turn of the game engine
func (game *Game) WorldTick() {
	tickStartTime := time.Now()
	game.turn++
	game.log.Printf("Game tick: %d", game.turn)
	changed := false
	for _, m := range game.currentDungeon.Mobs() {
		changed = m.Tick(game.turn) || changed
	}
	if changed {
		game.ui.MarkDirty()
	}
	game.currentDungeon.ResetFlag(FlagLit | FlagVisible)
	game.currentDungeon.CalculateLighting()
	game.currentDungeon.OnTilesInLineOfSight(game.player.Loc(), game.player.VisionRadius(), func(t *Tile, loc Vec) {
		if t.Lit() {
			t.flags |= FlagVisible | FlagSeen
		}
	})

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
