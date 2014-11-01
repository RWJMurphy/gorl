package gorl

import (
	"fmt"
	"log"

	"github.com/imdario/mergo"
	"github.com/nsf/termbox-go"
)

// Paintable is anything that can be painted
type Paintable interface {
	Paint()
}

// UIState is the state of the UI
type UIState int

const (
	// StateGame is the default UI state -- shows the MapWidget, MessageLogWidget and MenuWidget; waits for input
	StateGame UIState = iota
	// StateClosed is a closed UI. Entering this state is a signal to shut the game down cleanly.
	StateClosed
)

func (state UIState) String() string {
	switch state {
	case StateGame:
		return "StateGame"
	case StateClosed:
		return "StateClosed"
	default:
		return string(int(state))
	}
}

// The UI for the Game
type UI struct {
	Paintables    []Paintable
	cameraWidget  *CameraWidget
	menuWidget    *MenuWidget
	messageWidget *MessageLogWidget
	messages      []string
	state         UIState
	game          *Game
	dirty         bool
	log           log.Logger
}

// A BoxStyle is a collection of runes for drawing boxes with PaintBorder
type BoxStyle struct {
	horizontal rune
	vertical   rune
	corner     rune
}

// The DefaultBoxStyle paints boxes that look like
//   +--+
//   |  |
//   +--+
var DefaultBoxStyle = BoxStyle{'-', '|', '+'}

// Close cleans up after the UI
func (ui *UI) Close() {
	termbox.Close()
}

// NewUI creates a new UI, initiates termbox, and returns the UI.
// Please `defer ui.Close`.
func NewUI(game *Game) (*UI, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}

	width, height := termbox.Size()
	ui := &UI{}
	ui.game = game
	ui.log = game.log
	ui.messages = make([]string, 0, 10)
	ui.messageWidget = &MessageLogWidget{
		Widget{
			0, height - height/4,
			width, height / 4,
			ui,
		},
		nil,
	}
	ui.cameraWidget = &CameraWidget{
		Widget{
			0, 0,
			width - width/4, height - height/4,
			ui,
		},
		nil,
		Coord{0, 0},
	}
	ui.menuWidget = &MenuWidget{
		Widget{
			width - width/4, 0,
			width / 4, height - height/4,
			ui,
		},
	}
	ui.Paintables = []Paintable{
		ui.messageWidget,
		ui.cameraWidget,
		ui.menuWidget,
	}
	ui.state = StateGame
	ui.dirty = true
	return ui, nil
}

// A Widget represents a rectangular box in a fixed position in the UI.
// When painted, it draws a border around itself.
type Widget struct {
	x, y          int
	width, height int
	ui            *UI
}

// A CameraWidget renders part of a Dungeon
type CameraWidget struct {
	Widget
	dungeon *Dungeon
	center  Coord
}

// Paint paints the CameraWidget to the UI
func (camera *CameraWidget) Paint() {
	var (
		tile  *Tile
		loc   Coord
		out   Coord
		x, y  int
		f     Feature
		ok    bool
		color termbox.Attribute
	)
	ne := Coord{camera.center.x - camera.width/2, camera.center.y - camera.height/2}

	for x = 0; x < camera.width; x++ {
		for y = 0; y < camera.height; y++ {
			loc = Coord{ne.x + x, ne.y + y}
			out = Coord{camera.x + x, camera.y + y}
			tile = camera.dungeon.Tile(loc.x, loc.y)
			if tile.Seen() || tile.Visible() {
				color = tile.color
				if tile.Visible() {
					camera.ui.PutRuneColor(out.x, out.y, tile.c, color|termbox.AttrBold, termbox.ColorDefault)
					if f, ok = camera.dungeon.features[loc]; ok {
						camera.ui.PutRuneColor(out.x, out.y, f.Char(), f.Color(), termbox.ColorDefault)
					}
					if f, ok = camera.dungeon.mobs[loc]; ok {
						camera.ui.PutRuneColor(out.x, out.y, f.Char(), f.Color(), termbox.ColorDefault)
					}
				} else {
					camera.ui.PutRuneColor(out.x, out.y, tile.c, color, termbox.ColorDefault)
				}
			}
		}
	}
	camera.Widget.Paint()
}

// A MessageLogWidget renders messages
type MessageLogWidget struct {
	Widget
	messages []string
}

// Paint paints the MessageLogWidget to the UI
func (messageLog *MessageLogWidget) Paint() {
	for i, m := range messageLog.ui.messages {
		messageLog.ui.PrintAt(messageLog.x+1, messageLog.y+1+i, m)
	}
	messageLog.Widget.Paint()
}

// A MenuWidget in theory, displays a menu.
type MenuWidget struct {
	Widget
}

// Paint paints the Widget to the UI
func (w *Widget) Paint() {
	w.ui.PaintBorder(w.x, w.y, w.x+w.width-1, w.y+w.height-1, DefaultBoxStyle)
}

// Paint paints the MenuWidget to the UI
func (menu *MenuWidget) Paint() {
	menu.Widget.Paint()
}

// PutRuneColor paints the rune r at position x, y with foreground color fg and
// background bg.
func (ui *UI) PutRuneColor(x, y int, r rune, fg, bg termbox.Attribute) {
	termbox.SetCell(x, y, r, fg, bg)
}

// PutRune paints the rune r at positions x, y in the default colors.
func (ui *UI) PutRune(x, y int, r rune) {
	ui.PutRuneColor(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
}

// PrintAt paints a string to the UI, left to right starting at x, y
func (ui *UI) PrintAt(x, y int, s string) {
	for i, r := range s {
		ui.PutRune(x+i, y, r)
	}
}

// PaintBox fills a rectangular section with rune r. The rectangle is defined by
// its corners (x1, y1) and (x2, y2).
func (ui *UI) PaintBox(x1, y1, x2, y2 int, r rune) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			ui.PutRune(x, y, r)
		}
	}
}

// PaintBorder paints a border along the rectangle (x1, y2), (x2, y2)
// with the runes defines by style.
func (ui *UI) PaintBorder(x1, y1, x2, y2 int, style BoxStyle) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	if err := mergo.Merge(&style, DefaultBoxStyle); err != nil {
		ui.log.Panic(err)
	}

	ui.PaintBox(x1+1, y1, x2-1, y1, style.horizontal)
	ui.PaintBox(x1+1, y2, x2-1, y2, style.horizontal)

	ui.PaintBox(x1, y1+1, x1, y2-1, style.vertical)
	ui.PaintBox(x2, y1+1, x2, y2-1, style.vertical)

	ui.PutRune(x1, y1, style.corner)
	ui.PutRune(x1, y2, style.corner)
	ui.PutRune(x2, y1, style.corner)
	ui.PutRune(x2, y2, style.corner)
}

// Paint redraws the UI and its Paintables if the UI has been marked as dirty.
func (ui *UI) Paint() {
	if !ui.dirty {
		ui.log.Println("ui.Paint: ui not dirty, nothing to do")
		return
	}
	ui.log.Println("ui.Paint: painting")
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for _, p := range ui.Paintables {
		p.Paint()
	}
	termbox.Flush()
	ui.dirty = false
}

// PointCameraAt sets the dungeon and center for the CameraWidget
func (ui *UI) PointCameraAt(d *Dungeon, c Coord) {
	ui.cameraWidget.dungeon = d
	ui.cameraWidget.center = c
}

// HandleKey handles a termbox.KeyEvent
func (ui *UI) HandleKey(char rune, key termbox.Key) GameState {
	switch char {
	// Quit
	case 'q':
		return GameClosed
	// Move
	case 'h', 'j', 'k', 'l', 'y', 'u', 'b', 'n':
		moved := ui.HandleMovementKey(char, key)
		if moved {
			return GameWorldTurn
		}
	case 0:
		switch key {
		// Quit
		case termbox.KeyCtrlC, termbox.KeyEsc:
			return GameClosed
		// Wait
		case termbox.KeySpace:
			return GameWorldTurn
		// Move
		case termbox.KeyArrowUp, termbox.KeyArrowRight, termbox.KeyArrowDown, termbox.KeyArrowLeft:
			if moved := ui.HandleMovementKey(char, key); moved {
				return GameWorldTurn
			}
		default:
			msg := fmt.Sprintf("Unhandled key: %s", string(key))
			ui.log.Println(msg)
			ui.game.AddMessage(msg)
		}
	default:
		msg := fmt.Sprintf("Unhandled key: %c", char)
		ui.log.Println(msg)
		ui.game.AddMessage(msg)
	}
	return ui.game.state
}

// Single tile Movement constants
var (
	MoveNorth     = Movement{0, -1}
	MoveNorthEast = Movement{1, -1}
	MoveEast      = Movement{1, 0}
	MoveSouthEast = Movement{1, 1}
	MoveSouth     = Movement{0, 1}
	MoveSouthWest = Movement{-1, 1}
	MoveWest      = Movement{-1, 0}
	MoveNorthWest = Movement{-1, -1}
)

// HandleMovementKey maps a key to its respective Movement, and passes it
// to Game.Move. Returns true if the move was successful.
func (ui *UI) HandleMovementKey(char rune, key termbox.Key) bool {
	var movement Movement
	switch char {
	case 'k':
		movement = MoveNorth
	case 'u':
		movement = MoveNorthEast
	case 'l':
		movement = MoveEast
	case 'n':
		movement = MoveSouthEast
	case 'j':
		movement = MoveSouth
	case 'b':
		movement = MoveSouthWest
	case 'h':
		movement = MoveWest
	case 'y':
		movement = MoveNorthWest
	case 0:
		switch key {
		case termbox.KeyArrowUp:
			movement = MoveNorth
		case termbox.KeyArrowRight:
			movement = MoveEast
		case termbox.KeyArrowDown:
			movement = MoveSouth
		case termbox.KeyArrowLeft:
			movement = MoveWest
		default:
			ui.log.Panicf("Not a movement key: %s", string(key))
		}
	default:
		ui.log.Panicf("Not a movement key: %c", char)
	}
	return ui.game.Move(movement)
}

// HandleEvent handles a termbox.Event
func (ui *UI) HandleEvent(e termbox.Event) GameState {
	switch e.Type {
	case termbox.EventKey:
		return ui.HandleKey(e.Ch, e.Key)
	case termbox.EventResize:
		ui.dirty = true
		return ui.game.state
	default:
		fallthrough
	case termbox.EventError:
		ui.log.Panic(e.Err)
		return GameInvalidState
	}
}

// WaitAndHandleInput waits on a termbox.Event, handles it, and repaints the UI if needed.
// Returns a GameState.
func (ui *UI) WaitAndHandleInput() GameState {
	event := termbox.PollEvent()
	nextState := ui.HandleEvent(event)
	return nextState
}
