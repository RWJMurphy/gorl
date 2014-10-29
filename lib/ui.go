package gorl

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/nsf/termbox-go"
	"unicode/utf8"
)

type Paintable interface {
	Paint()
}

type UiState int

const (
	StateGame UiState = iota
	StateClosed
)

func (state UiState) String() string {
	switch state {
	case StateGame:
		return "StateGame"
	case StateClosed:
		return "StateClosed"
	default:
		return string(int(state))
	}
}

type UI struct {
	Paintables    []Paintable
	cameraWidget  *CameraWidget
	menuWidget    *MenuWidget
	messageWidget *MessageLogWidget
	messages      []string
	State         UiState
	game          *Game
	dirty         bool
}

type BoxStyle struct {
	horizontal rune
	vertical   rune
	corner     rune
}

var DefaultBoxStyle = BoxStyle{'-', '|', '+'}

func (ui *UI) Close() {
	termbox.Close()
}

func NewUI() (*UI, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}

	width, height := termbox.Size()
	ui := new(UI)
	ui.messages = make([]string, 0, 10)
	ui.messageWidget = &MessageLogWidget{
		0, height - height/4,
		width, height / 4,
		nil,
		ui,
	}
	ui.cameraWidget = &CameraWidget{
		0, 0,
		width - width/4, height - height/4,
		nil,
		Coord{0, 0},
		ui,
	}
	ui.menuWidget = &MenuWidget{
		width - width/4, 0,
		width / 4, height - height/4,
		ui,
	}
	ui.Paintables = []Paintable{
		ui.messageWidget,
		ui.cameraWidget,
		ui.menuWidget,
	}
	ui.State = StateGame
	ui.dirty = true
	return ui, nil
}

type CameraWidget struct {
	x, y          int
	width, height int
	dungeon       *Dungeon
	center        Coord
	ui            *UI
}

func (camera *CameraWidget) Paint() {
	ne := Coord{camera.center.x - camera.width/2, camera.center.y - camera.height/2}
	for x := 0; x < camera.width; x++ {
		for y := 0; y < camera.height; y++ {
			tile_x, tile_y := ne.x+x, ne.y+y
			out_x, out_y := camera.x+x, camera.y+y
			tile := camera.dungeon.Tile(tile_x, tile_y)
			if tile.flags&FlagVisible != 0 {
				if tile.flags&FlagLit != 0 {
					camera.ui.PutRune(out_x, out_y, tile.c)
				}
			}
		}
	}
	for _, m := range camera.dungeon.mobs {
		out_x, out_y := m.Loc().x-ne.x, m.Loc().y-ne.y
		if out_x > 0 && out_x < camera.width && out_y > 0 && out_y < camera.height {
			camera.ui.PutRune(m.Loc().x-ne.x, m.Loc().y-ne.y, m.Char())
		}
	}
	camera.ui.PaintBorder(camera.x, camera.y, camera.x+camera.width-1, camera.y+camera.height-1, DefaultBoxStyle)
}

type MessageLogWidget struct {
	x, y          int
	width, height int
	messages      []string
	ui            *UI
}

func (messageLog *MessageLogWidget) Paint() {
	for i, m := range messageLog.ui.messages {
		messageLog.ui.PrintAt(messageLog.x+1, messageLog.y+1+i, m)
	}
	messageLog.ui.PaintBorder(messageLog.x, messageLog.y, messageLog.x+messageLog.width-1, messageLog.y+messageLog.height-1, DefaultBoxStyle)
}

type MenuWidget struct {
	x, y          int
	width, height int
	ui            *UI
}

func (menu *MenuWidget) Paint() {
	menu.ui.PaintBorder(menu.x, menu.y, menu.x+menu.width-1, menu.y+menu.height-1, DefaultBoxStyle)
}

func (ui *UI) PutRune(x, y int, r rune) {
	termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
}

func (ui *UI) PrintAt(x, y int, s string) {
	for i, r := range s {
		ui.PutRune(x+i, y, r)
	}
}

func (ui *UI) PrintCentered(s string) {
	width, height := termbox.Size()
	mid_w, mid_h := width/2, height/2
	s_len := utf8.RuneCountInString(s)
	ui.PrintAt(mid_w-s_len/2, mid_h, s)
}

func (ui *UI) PaintBox(x1, y1, x2, y2 int, r rune) {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			ui.PutRune(x, y, r)
		}
	}
}

func (ui *UI) PaintVerticalLine(x, y1, y2 int, r rune) {
	for y := y1; y <= y2; y++ {
		ui.PutRune(x, y, r)
	}
}

func (ui *UI) PaintHorizontalLine(x1, y, x2 int, r rune) {
	for x := x1; x <= x2; x++ {
		ui.PutRune(x, y, r)
	}
}

func (ui *UI) PaintBorder(x1, y1, x2, y2 int, style BoxStyle) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	if err := mergo.Merge(&style, DefaultBoxStyle); err != nil {
		panic(err)
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

func (ui *UI) Paint() {
	if ! ui.dirty {
		return
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for _, p := range ui.Paintables {
		p.Paint()
	}
	termbox.Flush()
	ui.dirty = false
}

func (ui *UI) PointCameraAt(d *Dungeon, c Coord) {
	ui.cameraWidget.dungeon = d
	ui.cameraWidget.center = c
}

func (ui *UI) HandleKey(char rune, key termbox.Key) {
	switch char {
	case 'q':
		ui.State = StateClosed
	case 'h', 'j', 'k', 'l':
		ui.HandleMovementKey(char, key)
	case 0:
		switch key {
		case termbox.KeyCtrlC, termbox.KeyEsc:
			ui.State = StateClosed
		case termbox.KeyArrowUp, termbox.KeyArrowRight, termbox.KeyArrowDown, termbox.KeyArrowLeft:
			ui.HandleMovementKey(char, key)
		default:
			ui.game.AddMessage(fmt.Sprintf("Unhandled key: %s", key))
		}
	default:
		ui.game.AddMessage(fmt.Sprintf("Unhandled key: %c", char))
	}
}

func (ui *UI) HandleMovementKey(char rune, key termbox.Key) {
	var movement Movement
	switch char {
	case 'h':
			movement = Movement{-1, 0}
	case 'j':
			movement = Movement{0, 1}
	case 'k':
			movement = Movement{0, -1}
	case 'l':
			movement = Movement{1, 0}
	case 0:
		switch key {
		case termbox.KeyArrowUp:
			movement = Movement{0, -1}
		case termbox.KeyArrowRight:
			movement = Movement{1, 0}
		case termbox.KeyArrowDown:
			movement = Movement{0, 1}
		case termbox.KeyArrowLeft:
			movement = Movement{-1, 0}
		default:
			panic(fmt.Sprintf("Not a movement key: %s", key))
		}
	default:
		panic(fmt.Sprintf("Not a movement key: %s", char))
	}
	ui.game.Move(movement)
}

func (ui *UI) HandleEvent(e termbox.Event) {
	switch e.Type {
	case termbox.EventKey:
		ui.HandleKey(e.Ch, e.Key)
	case termbox.EventResize:
		ui.dirty = true
	case termbox.EventError:
		panic(e.Err)
	}
}

func (ui *UI) Tick() {
	event := termbox.PollEvent()
	ui.HandleEvent(event)
	ui.Paint()
}
