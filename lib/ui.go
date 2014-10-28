package gorl

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/nsf/termbox-go"
	"strings"
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
	State         UiState
	game          *Game
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
	ui.messageWidget = &MessageLogWidget{
		0, height - height/4,
		width, height / 4,
		strings.Repeat("Hello, termbox. ", 4),
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
			camera.ui.PutRune(out_x, out_y, tile)
		}
	}
	for _, m := range camera.dungeon.mobs {
		camera.ui.PutRune(m.loc.x-ne.x, m.loc.y-ne.y, m.c)
	}
	camera.ui.PaintBorder(camera.x, camera.y, camera.x+camera.width-1, camera.y+camera.height-1, DefaultBoxStyle)
}

type MessageLogWidget struct {
	x, y          int
	width, height int
	s             string
	ui            *UI
}

func (messageLog *MessageLogWidget) Paint() {
	messageLog.ui.PrintAt(messageLog.x+1, messageLog.y+1, messageLog.s)
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
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for _, p := range ui.Paintables {
		p.Paint()
	}
	termbox.Flush()
}

func (ui *UI) PointCameraAt(d *Dungeon, c Coord) {
	ui.cameraWidget.dungeon = d
	ui.cameraWidget.center = c
}

func (ui *UI) HandleMovementKey(key termbox.Key) {
	var movement Movement
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
		panic(fmt.Sprintf("Not a supported movement key: %s", key))
	}
	ui.game.Move(movement)
}

func (ui *UI) HandleEvent(e termbox.Event) bool {
	dirty := false
	switch e.Type {
	case termbox.EventKey:
		if e.Ch != 0 {
			switch e.Ch {
			case 'q':
				ui.State = StateClosed
			}
		} else {
			switch e.Key {
			case termbox.KeyCtrlC, termbox.KeyEsc:
				ui.State = StateClosed
			case termbox.KeyArrowUp, termbox.KeyArrowRight, termbox.KeyArrowDown, termbox.KeyArrowLeft:
				ui.HandleMovementKey(e.Key)
				dirty = true
			}
		}
	case termbox.EventResize:
		dirty = true
	case termbox.EventError:
		panic(e.Err)
	}
	return dirty
}

func (ui *UI) Tick() {
	event := termbox.PollEvent()
	dirty := ui.HandleEvent(event)
	if dirty {
		ui.Paint()
	}
}
