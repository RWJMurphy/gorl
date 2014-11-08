package gorl

import (
	"fmt"
	"log"

	"github.com/imdario/mergo"
	"github.com/nsf/termbox-go"
)

type TermboxUI interface {
	UI

	Messages() []string
	PaintBorder(RectangleI, boxStyle)
	PutRuneColor(Vec, rune, termbox.Attribute, termbox.Attribute)
	PrintAt(Vec, string)
	PutRune(Vec, rune)
}

type termboxUI struct {
	paintables   []Paintable
	cameraWidget *cameraWidget
	menuWidget   *menuWidget
	logWidget    *logWidget
	messages     []string
	state        State
	game         *Game
	dirty        bool
	log          log.Logger
}

func NewTermboxUI(game *Game) (TermboxUI, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}

	width, height := termbox.Size()
	ui := &termboxUI{}
	ui.game = game
	ui.log = game.log
	ui.messages = make([]string, 0, 10)
	ui.logWidget = &logWidget{
		widget{
			Rectangle{
				Vec{0, height - height/4},
				Vec{width, height / 4},
			},
			ui,
		},
		nil,
	}
	ui.cameraWidget = &cameraWidget{
		widget{
			Rectangle{
				Vec{0, 0},
				Vec{width - width/4, height - height/4},
			},
			ui,
		},
		nil,
		Vec{0, 0},
	}
	ui.menuWidget = &menuWidget{
		widget{
			Rectangle{
				Vec{width - width/4, 0},
				Vec{width / 4, height - height/4},
			},
			ui,
		},
	}
	ui.paintables = []Paintable{
		ui.logWidget,
		ui.cameraWidget,
		ui.menuWidget,
	}
	ui.state = StateGame
	ui.dirty = true
	return ui, nil
}

// UI interface implementation

func (ui *termboxUI) Close() {
	termbox.Close()
}

func (ui *termboxUI) MarkDirty() {
	ui.dirty = true
}

func (ui *termboxUI) IsDirty() bool {
	return ui.dirty
}

func (ui *termboxUI) Paintables() []Paintable {
	return ui.paintables
}

func (ui *termboxUI) State() State {
	return ui.state
}

func (ui *termboxUI) Tick() (PlayerAction, GameState) {
	event := termbox.PollEvent()
	action, nextState := ui.HandleEvent(event)
	return action, nextState
}

// PointCameraAt sets the dungeon and center for the CameraWidget
func (ui *termboxUI) PointCameraAt(d *Dungeon, c Vec) {
	ui.cameraWidget.dungeon = d
	ui.cameraWidget.center = c
}

func (ui *termboxUI) MessagesWanted() int {
	return ui.logWidget.Height() - 2
}

func (ui *termboxUI) SetMessages(messages []string) {
	ui.messages = messages
}

// Paint redraws the UI and its Paintables if the UI has been marked as dirty.
func (ui *termboxUI) Paint() {
	if !ui.dirty {
		ui.log.Println("ui.Paint: ui not dirty, nothing to do")
		return
	}
	ui.log.Println("ui.Paint: painting")
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for _, p := range ui.paintables {
		p.Paint()
	}
	termbox.Flush()
	ui.dirty = false
}

// TermboxUI implementation

func (ui *termboxUI) Messages() []string {
	return ui.messages
}

// PutRuneColor paints the rune r at position x, y with foreground color fg and
// background bg.
func (ui *termboxUI) PutRuneColor(loc Vec, r rune, fg, bg termbox.Attribute) {
	termbox.SetCell(loc.x, loc.y, r, fg, bg)
}

// PutRune paints the rune r at positions x, y in the default colors.
func (ui *termboxUI) PutRune(loc Vec, r rune) {
	ui.PutRuneColor(loc, r, termbox.ColorDefault, termbox.ColorDefault)
}

// PrintAt paints a string to the UI, left to right starting at x, y
func (ui *termboxUI) PrintAt(loc Vec, s string) {
	for i, r := range s {
		ui.PutRune(loc.Plus(Vec{i, 0}), r)
	}
}

// PaintBox fills a rectangular section with rune r. The rectangle is defined by
// its corners (x1, y1) and (x2, y2).
func (ui *termboxUI) PaintBox(rect RectangleI, r rune) {
	ui.log.Printf("Painting box %s", rect)
	x1, y1 := rect.TopLeft().x, rect.TopLeft().y
	x2, y2 := rect.BottomRight().x, rect.BottomRight().y

	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			ui.PutRune(Vec{x, y}, r)
		}
	}
}

// PaintBorder paints a border along the rectangle (x1, y2), (x2, y2)
// with the runes defines by style.
func (ui *termboxUI) PaintBorder(rect RectangleI, style boxStyle) {
	ui.log.Printf("Painting border around %s", rect)
	x1, y1 := rect.TopLeft().x, rect.TopLeft().y
	x2, y2 := rect.BottomRight().x, rect.BottomRight().y

	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	if err := mergo.Merge(&style, defaultBoxStyle); err != nil {
		ui.log.Panic(err)
	}

	ui.PaintBox(Rectangle{Vec{x1 + 1, y1}, Vec{rect.Width() - 2, 1}}, style.horizontal)
	ui.PaintBox(Rectangle{Vec{x1 + 1, y2}, Vec{rect.Width() - 2, 1}}, style.horizontal)

	ui.PaintBox(Rectangle{Vec{x1, y1 + 1}, Vec{1, rect.Height() - 2}}, style.vertical)
	ui.PaintBox(Rectangle{Vec{x2, y1 + 1}, Vec{1, rect.Height() - 2}}, style.vertical)

	ui.PutRune(rect.TopLeft(), style.corner)
	ui.PutRune(rect.TopRight(), style.corner)
	ui.PutRune(rect.BottomRight(), style.corner)
	ui.PutRune(rect.BottomLeft(), style.corner)
}

// HandleKey handles a termbox.KeyEvent
func (ui *termboxUI) HandleKey(char rune, key termbox.Key) (PlayerAction, GameState) {
	switch char {
	// Quit
	case 'q':
		return ActNone, GameClosed
	// Move
	case 'h', 'j', 'k', 'l', 'y', 'u', 'b', 'n':
		// TODO: move movement handling to Game
		moved := ui.HandleMovementKey(char, key)
		if moved {
			return ActNone, GameWorldTurn
		}
	// Drop all
	case 'D':
		return ActDropAll, GameWorldTurn
	case ',', 'g':
		return ActPickUp, GameWorldTurn
	case 0:
		switch key {
		// Quit
		case termbox.KeyCtrlC, termbox.KeyEsc:
			return ActNone, GameClosed
		// Wait
		case termbox.KeySpace:
			return ActWait, GameWorldTurn
		// Move
		case termbox.KeyArrowUp, termbox.KeyArrowRight, termbox.KeyArrowDown, termbox.KeyArrowLeft:
			// TODO: move movement handling to Game
			if moved := ui.HandleMovementKey(char, key); moved {
				return ActNone, GameWorldTurn
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
	return ActNone, ui.game.state
}

// HandleMovementKey maps a key to its respective Vec, and passes it
// to Game.Move. Returns true if the move was successful.
func (ui *termboxUI) HandleMovementKey(char rune, key termbox.Key) bool {
	var movement Vec
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
	return ui.game.MoveOrAct(movement)
}

// HandleEvent handles a termbox.Event
func (ui *termboxUI) HandleEvent(e termbox.Event) (PlayerAction, GameState) {
	switch e.Type {
	case termbox.EventKey:
		return ui.HandleKey(e.Ch, e.Key)
	case termbox.EventResize:
		ui.dirty = true
		return ActNone, ui.game.state
	case termbox.EventError:
		ui.log.Panic(e.Err)
		fallthrough
	default:
		return ActNone, GameInvalidState
	}
}

type boxStyle struct {
	horizontal rune
	vertical   rune
	corner     rune
}

var defaultBoxStyle = boxStyle{'-', '|', '+'}
