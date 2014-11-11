package gorl

import "fmt"

// Paintable is anything that can be painted
type Paintable interface {
	Paint()
}

// State is the state of the UI
type State int

const (
	// StateInvalid represents any invalid UI State
	StateInvalid State = iota
	// StateGame is the default UI state -- shows the MapWidget, MessageLogWidget and MenuWidget; waits for input
	StateGame
	// StateInventory displays the inventory view
	StateInventory
	// StateClosed is a closed UI. Entering this state is a signal to shut the game down cleanly.
	StateClosed
)

func (state State) String() string {
	switch state {
	case StateInvalid:
		return "StateInvalid"
	case StateGame:
		return "StateGame"
	case StateClosed:
		return "StateClosed"
	case StateInventory:
		return "StateInventory"
	default:
		return fmt.Sprintf("State(%d)", state)
	}
}

type UI interface {
	Close()
	Paintables() []Paintable
	State() State
	MarkDirty()
	IsDirty() bool
	Paint()
	DoEvent() (MobAction, GameState)

	PointCameraAt(*Dungeon, Vec)

	MessagesWanted() int
	SetMessages([]string)
}

// A Widget represents a rectangular box in a fixed position in the UI.
type Widget interface {
	RectangleI
	SetRectangle(RectangleI)
	Paint()
}

// A CameraWidget renders part of a Dungeon
type CameraWidget interface {
	Widget
	Dungeon() *Dungeon
	SetDungeon(*Dungeon)
	Center() Vec
	SetCenter(Vec)
}

// A LogWidget renders messages
type LogWidget interface {
	Widget
}

// A MenuWidget in theory, displays a menu.
type MenuWidget interface {
	Widget
}

type InventoryWidget interface {
	Widget
	SetOwner(Mob)
}

// Single tile Movement constants
var (
	MoveNorth     = Vec{0, -1}
	MoveNorthEast = Vec{1, -1}
	MoveEast      = Vec{1, 0}
	MoveSouthEast = Vec{1, 1}
	MoveSouth     = Vec{0, 1}
	MoveSouthWest = Vec{-1, 1}
	MoveWest      = Vec{-1, 0}
	MoveNorthWest = Vec{-1, -1}
)
