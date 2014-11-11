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

type MobAction struct {
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
		dest := Vec{rand.Intn(dungeon.width), rand.Intn(dungeon.height)}
		for !(dungeon.Tile(dest).Crossable() && dungeon.FeatureGroup(dest).Crossable()) {
			dest = Vec{rand.Intn(dungeon.width), rand.Intn(dungeon.height)}
		}
		mob := NewMob(fmt.Sprintf("orc #%d", i), 'o', game.log, dungeon)
		mob.SetVisionRadius(5)
		mob.SetColor(termbox.ColorGreen)
		mob.SetLoc(dest)

		torch := NewItem("torch", '!', 1)
		torch.SetLightRadius(10)
		mob.AddToInventory(torch)

		dungeon.AddMob(mob)
	}

	ui, err := NewTermboxUI(game)
	if err != nil {
		return nil, err
	}
	game.ui = ui
	game.SetDungeon(dungeon)
	game.updatePlayerFOV()

	game.AddMessage("Welcome to GoRL!")
	game.state = GamePlayerTurn
	return game, nil
}

func (game *Game) updatePlayerFOV() {
	game.currentDungeon.ResetFlag(FlagLit | FlagVisible)
	game.currentDungeon.CalculateLighting()
	game.currentDungeon.OnTilesInLineOfSight(game.player.Loc(), game.player.VisionRadius(), func(t *Tile, loc Vec) {
		if t.Lit() {
			t.flags |= FlagVisible | FlagSeen
		}
	})
}

// MoveOrAct calculates the destination tile based on the movement parameter and
// the Player's location, and then
//   * if there is a mob on the destination, attacks the mob and returns true
//   * if not and destination is Crossable, moves the player there and returns true
//   * if the destination is not Crossable, returns false
func (game *Game) MoveOrAct(mob Mob, movement Vec) bool {
	game.log.Printf("%s MoveOrAct'ing %s", mob, movement)
	destination := mob.Loc().Add(movement)
	if otherMob := game.currentDungeon.MobAt(destination); otherMob != nil {
		if damageDealt, ok := mob.Attack(otherMob); ok {
			game.AddMessage(fmt.Sprintf("%s hit %s for %d damage", mob.Name(), otherMob.Name(), damageDealt))
			if otherMob.Dead() {
				game.AddMessage(fmt.Sprintf("The %s dies!", otherMob.Name()))
			}
			return true
		}
		return false
	} else if moved := game.currentDungeon.MoveMob(mob, movement); !moved {
		return moved
	}

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
	var action MobAction
	nextState := game.state

	game.ui.Paint()
	game.log.Println("Entering main loop")
mainLoop:
	for {
		game.log.Printf("Game state: %s", game.state)
		switch game.state {
		case GameWorldTurn:
			game.WorldTick()
			nextState = GamePlayerTurn
		case GamePlayerTurn:
			action, nextState = game.ui.DoEvent()
			if game.doMobAction(game.player, action) {
				game.state = nextState
				game.ui.PointCameraAt(game.currentDungeon, game.player.Loc())
				game.updatePlayerFOV()
				game.ui.MarkDirty()
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

func (game *Game) doMobAction (mob Mob, action MobAction) bool {
	switch action.action {
	case ActWait:
		return true
	case ActDrop:
		item := action.target.(Item)
		if item == nil {
			game.log.Panicf("%s tried to drop nil!", mob)
			return false
		}
		mob.DropItem(item, game.currentDungeon)
		game.AddMessage(fmt.Sprintf("%s dropped %s", mob.Name(), item.Name()))
		return true
	case ActDropAll:
		for _, item := range mob.Inventory() {
			mob.DropItem(item, game.currentDungeon)
			game.AddMessage(fmt.Sprintf("%s dropped %s", mob.Name(), item.Name()))
		}
		return true
	case ActPickUpAll:
		items := game.currentDungeon.ItemsAt(mob.Loc())
		if len(items) > 0 {
			for _, item := range items {
				game.currentDungeon.DeleteItem(item)
				mob.AddToInventory(item)
				game.AddMessage(fmt.Sprintf("%s picked up %s", mob.Name(), item.Name()))
			}
			return true
		} else {
			game.AddMessage(fmt.Sprintf("Silly %s, there's nothing to pick up.", mob.Name()))
			return false
		}
	case ActMove:
		direction := action.target.(Vec)
		return game.MoveOrAct(mob, direction)
	case ActNone:
		return false
	default:
		log.Panicf("Bad action: %s, %s", mob, action)
		return false
	}
	return false
}

// WorldTick runs a single turn of the game engine
func (game *Game) WorldTick() {
	tickStartTime := time.Now()
	game.turn++
	game.log.Printf("Game tick: %d", game.turn)
	changed := false
	var mobAction MobAction
	for _, mob := range game.currentDungeon.Mobs() {
		mobAction = mob.Tick(game.turn)
		changed = game.doMobAction(mob, mobAction) || changed
	}
	if changed {
		game.ui.MarkDirty()
	}
	game.updatePlayerFOV()

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
