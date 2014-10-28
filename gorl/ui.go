package gorl

import (
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
	Paintables []Paintable
	cameraWidget *CameraWidget
	menuWidget *MenuWidget
	messageWidget *MessageLogWidget
	State UiState
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

	w, h := termbox.Size()
	ui := new(UI)
	ui.messageWidget = &MessageLogWidget{
		0, h - h/4,
		w, h / 4,
		strings.Repeat("Hello, termbox. ", 4),
		ui,
	}
	ui.cameraWidget = &CameraWidget{
		0, 0,
		w - w/4, h - h/4,
		ui,
	}
	ui.menuWidget = &MenuWidget{
		w - w/4, 0,
		w / 4, h - h/4,
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
	x, y 	int
	w, h 	int
	ui   	*UI
}

func (w *CameraWidget) Paint() {
	w.ui.PaintBorder(w.x, w.y, w.x+w.w-1, w.y+w.h-1, DefaultBoxStyle)
}

type MessageLogWidget struct {
	x, y int
	w, h int
	s    string
	ui   *UI
}

func (w *MessageLogWidget) Paint() {
	w.ui.PrintAt(w.x+1, w.y+1, w.s)
	w.ui.PaintBorder(w.x, w.y, w.x+w.w-1, w.y+w.h-1, DefaultBoxStyle)
}

type MenuWidget struct {
	x, y int
	w, h int
	ui   *UI
}

func (w *MenuWidget) Paint() {
	w.ui.PaintBorder(w.x, w.y, w.x+w.w-1, w.y+w.h-1, DefaultBoxStyle)
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

func (ui *UI) HandleEvent(e termbox.Event) bool {
	dirty := false
	switch e.Type {
	case termbox.EventKey:
		switch e.Key {
		case termbox.KeyCtrlC, termbox.KeyEsc:
			ui.State = StateClosed
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
