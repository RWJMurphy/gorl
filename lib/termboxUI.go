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
	PutRuneColor(Vector, rune, termbox.Attribute, termbox.Attribute)
	PrintAt(Vector, string)
	PutRune(Vector, rune)
}

type termboxUI struct {
	paintables      []Paintable
	cameraWidget    *cameraWidget
	menuWidget      *menuWidget
	logWidget       *logWidget
	inventoryWidget *inventoryWidget
	messages        []string
	state           State
	game            *Game
	dirty           bool
	log             *log.Logger
	// ugh this is hacky
	stateAction MobAction
}

func NewTermboxUI(game *Game) (TermboxUI, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}

	ui := &termboxUI{}
	ui.game = game
	ui.log = game.log
	ui.messages = make([]string, 0, 10)
	ui.logWidget = &logWidget{
		widget{Rectangle{}, ui},
		nil,
	}
	ui.cameraWidget = &cameraWidget{
		widget{Rectangle{}, ui},
		nil,
		Vector{0, 0},
	}
	ui.menuWidget = &menuWidget{
		widget{Rectangle{}, ui},
	}
	ui.inventoryWidget = &inventoryWidget{
		widget{Rectangle{}, ui},
		game.player,
	}
	ui.Resize()
	ui.setState(StateGame, MobAction{ActNone, nil})
	return ui, nil
}

func (ui *termboxUI) Resize() {
	width, height := termbox.Size()

	ui.cameraWidget.topLeft = Vector{0, 0}
	ui.cameraWidget.size = Vector{width - width/4, height - height/4}

	ui.menuWidget.topLeft = Vector{width - width/4, 0}
	ui.menuWidget.size = Vector{width / 4, height - height/4}

	ui.logWidget.topLeft = Vector{0, height - height/4}
	ui.logWidget.size = Vector{width, height / 4}

	ui.inventoryWidget.topLeft = Vector{0, 0}
	ui.inventoryWidget.size = Vector{width, height - height/4}

	ui.log.Println(ui.cameraWidget)
	ui.log.Println(ui.menuWidget)
	ui.log.Println(ui.logWidget)
	ui.log.Println(ui.inventoryWidget)

	ui.MarkDirty()
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

func (ui *termboxUI) DoEvent() (MobAction, GameState) {
	action := MobAction{ActNone, nil}
	nextState := ui.game.state

	switch ui.State() {
	case StateGame, StateInventory:
		event := termbox.PollEvent()
		action, nextState = ui.HandleEvent(event)
	case StateClosed:
		ui.log.Panic("Can't handle event while closed")
	default:
		ui.log.Panicf("Can't handle event in state %s", ui.State())
	}
	return action, nextState
}

// PointCameraAt sets the dungeon and center for the CameraWidget
func (ui *termboxUI) PointCameraAt(d *Dungeon, c Vector) {
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

// HandleKey handles a termbox.KeyEvent
//
// XXX I really dislike how coupled this is to GameState :(
func (ui *termboxUI) HandleKey(char rune, key termbox.Key) (MobAction, GameState) {
	switch ui.State() {
	case StateGame:
		switch char {
		// Quit
		case 'q':
			return MobAction{ActNone, nil}, GameClosed
		// Move
		case 'h', 'j', 'k', 'l', 'y', 'u', 'b', 'n':
			action := ui.HandleMovementKey(char, key)
			if action.action == ActMove {
				return action, GameWorldTurn
			}
			return action, ui.game.state
		case 'i':
			ui.setState(StateInventory, MobAction{ActNone, nil})
			return MobAction{ActNone, nil}, GamePlayerTurn
		// Drop
		case 'd':
			ui.setState(StateInventory, MobAction{ActDrop, nil})
			return MobAction{ActNone, nil}, GamePlayerTurn
		// Drop all
		case 'D':
			return MobAction{ActDropAll, nil}, GameWorldTurn
		case ',', 'g':
			return MobAction{ActPickUpAll, nil}, GameWorldTurn
		case 0:
			switch key {
			// Quit
			case termbox.KeyCtrlC, termbox.KeyEsc:
				return MobAction{ActNone, nil}, GameClosed
			// Wait
			case termbox.KeySpace:
				return MobAction{ActWait, nil}, GameWorldTurn
			// Move
			case termbox.KeyArrowUp, termbox.KeyArrowRight, termbox.KeyArrowDown, termbox.KeyArrowLeft:
				action := ui.HandleMovementKey(char, key)
				if action.action == ActMove {
					return action, GameWorldTurn
				}
				return action, ui.game.state
			}
		}
	case StateInventory:
		if char != 0 {
			inventoryIndex := int(char - 'a')
			inventory := ui.game.player.Inventory()
			if inventoryIndex >= 0 && inventoryIndex < len(inventory) {
				stateAction := ui.stateAction
				stateAction.target = inventory[inventoryIndex]
				ui.setState(StateGame, MobAction{ActNone, nil})
				return stateAction, GameWorldTurn
			}
			return MobAction{ActNone, nil}, ui.game.state
		}
		switch key {
		case termbox.KeyEsc:
			ui.setState(StateGame, MobAction{ActNone, nil})
			return MobAction{ActNone, nil}, ui.game.state
		}
	case StateClosed:
		ui.log.Panic("am closed, can't handle keys :(")
	}
	if char == 0 {
		ui.game.AddMessage(fmt.Sprintf("Unhandled key: %s", string(key)))
	} else {
		ui.game.AddMessage(fmt.Sprintf("Unhandled key: %c", char))
	}
	return MobAction{ActNone, nil}, ui.game.state
}

// HandleMovementKey maps a key to its respective Vector, and passes it
// to Game.Move. Returns true if the move was successful.
func (ui *termboxUI) HandleMovementKey(char rune, key termbox.Key) MobAction {
	var movement Vector
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
			return MobAction{ActNone, nil}
		}
	default:
		ui.log.Panicf("Not a movement key: %c", char)
		return MobAction{ActNone, nil}
	}
	return MobAction{ActMove, movement}
}

// HandleEvent handles a termbox.Event
func (ui *termboxUI) HandleEvent(e termbox.Event) (MobAction, GameState) {
	switch e.Type {
	case termbox.EventResize:
		ui.Resize()
		return MobAction{ActNone, nil}, ui.game.state
	case termbox.EventKey:
		return ui.HandleKey(e.Ch, e.Key)
	case termbox.EventError:
		ui.log.Panic(e.Err)
	}
	ui.log.Panicf("Unhandled event: %s", e)
	return MobAction{ActNone, nil}, GameInvalidState
}

func (ui *termboxUI) setState(state State, stateAction MobAction) {
	ui.log.Printf("termboxUI state change: %s -> %s", ui.state, state)
	ui.log.Printf("state expects action: %s", stateAction)
	ui.stateAction = stateAction
	if ui.state == state {
		return
	}
	ui.state = state
	ui.MarkDirty()
	switch state {
	case StateGame:
		ui.paintables = []Paintable{
			ui.cameraWidget,
			ui.logWidget,
			ui.menuWidget,
		}
	case StateInventory:
		ui.paintables = []Paintable{
			ui.inventoryWidget,
			ui.logWidget,
		}
	case StateClosed:
		ui.paintables = []Paintable{}
	default:
		ui.log.Panicf("Don't know how to setState(%s)", state)
	}
}

func (ui *termboxUI) Messages() []string {
	return ui.messages
}

// PutRuneColor paints the rune r at position x, y with foreground color fg and
// background bg.
func (ui *termboxUI) PutRuneColor(loc Vector, r rune, fg, bg termbox.Attribute) {
	termbox.SetCell(loc.x, loc.y, r, fg, bg)
}

// PutRune paints the rune r at positions x, y in the default colors.
func (ui *termboxUI) PutRune(loc Vector, r rune) {
	ui.PutRuneColor(loc, r, termbox.ColorDefault, termbox.ColorDefault)
}

// PrintAt paints a string to the UI, left to right starting at x, y
func (ui *termboxUI) PrintAt(loc Vector, s string) {
	for i, r := range s {
		ui.PutRune(loc.Add(Vector{i, 0}), r)
	}
}

// PaintBox fills a rectangular section with rune r. The rectangle is defined by
// its corners (x1, y1) and (x2, y2).
func (ui *termboxUI) PaintBox(rect RectangleI, r rune) {
	x1, y1 := rect.TopLeft().x, rect.TopLeft().y
	x2, y2 := rect.BottomRight().x, rect.BottomRight().y

	for x := x1; x < x2; x++ {
		for y := y1; y < y2; y++ {
			ui.PutRune(Vector{x, y}, r)
		}
	}
}

// PaintBorder paints a border along the rectangle (x1, y2), (x2, y2)
// with the runes defines by style.
func (ui *termboxUI) PaintBorder(rect RectangleI, style boxStyle) {
	x1, y1 := rect.TopLeft().x, rect.TopLeft().y
	x2, y2 := rect.BottomRight().x - 1, rect.BottomRight().y - 1

	if err := mergo.Merge(&style, defaultBoxStyle); err != nil {
		ui.log.Panic(err)
	}

	ui.PaintBox(Rectangle{Vector{x1 + 1, y1}, Vector{rect.Width() - 2, 1}}, style.horizontal)
	ui.PaintBox(Rectangle{Vector{x1 + 1, y2}, Vector{rect.Width() - 2, 1}}, style.horizontal)

	ui.PaintBox(Rectangle{Vector{x1, y1 + 1}, Vector{1, rect.Height() - 2}}, style.vertical)
	ui.PaintBox(Rectangle{Vector{x2, y1 + 1}, Vector{1, rect.Height() - 2}}, style.vertical)

	ui.PutRune(rect.TopLeft(), style.corner)
	ui.PutRune(rect.TopRight().Sub(Vector{1, 0}), style.corner)
	ui.PutRune(rect.BottomRight().Sub(Vector{1, 1}), style.corner)
	ui.PutRune(rect.BottomLeft().Sub(Vector{0, 1}), style.corner)
}

type boxStyle struct {
	horizontal rune
	vertical   rune
	corner     rune
}

var defaultBoxStyle = boxStyle{'-', '|', '+'}
