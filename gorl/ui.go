package gorl

import (
	"github.com/imdario/mergo"
	"github.com/nsf/termbox-go"
	"strings"
	"unicode/utf8"
)

type UI struct{}

type BoxStyle struct {
	horizontal rune
	vertical   rune
	corner     rune
}

var DefaultBoxStyle = BoxStyle{'-', '|', '+'}

func (ui UI) Close() {
	termbox.Close()
}

func NewUI() (*UI, error) {
	ui := &UI{}
	err := termbox.Init()
	return ui, err
}

func (ui UI) PutRune(x, y int, r rune) {
	termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
}
func (ui UI) PrintAt(x, y int, s string) {
	for i, r := range s {
		ui.PutRune(x+i, y, r)
	}
}

func (ui UI) PrintCentered(s string) {
	width, height := termbox.Size()
	mid_w, mid_h := width/2, height/2
	s_len := utf8.RuneCountInString(s)
	ui.PrintAt(mid_w-s_len/2, mid_h, s)
}

func (ui UI) DrawVerticalLine(x, y1, y2 int, r rune) {
	for y := y1; y <= y2; y++ {
		termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func (ui UI) DrawHorizontalLine(x1, y, x2 int, r rune) {
	for x := x1; x <= x2; x++ {
		termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func (ui UI) DrawRectangle(x1, y1, x2, y2 int, style BoxStyle) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	if err := mergo.Merge(&style, DefaultBoxStyle); err != nil {
		panic(err)
	}

	ui.DrawHorizontalLine(x1+1, y1, x2-1, style.horizontal)
	ui.DrawHorizontalLine(x1+1, y2, x2-1, style.horizontal)

	ui.DrawVerticalLine(x1, y1+1, y2-1, style.vertical)
	ui.DrawVerticalLine(x2, y1+1, y2-1, style.vertical)

	ui.PutRune(x1, y1, style.corner)
	ui.PutRune(x1, y2, style.corner)
	ui.PutRune(x2, y1, style.corner)
	ui.PutRune(x2, y2, style.corner)
}

func (ui UI) Draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	w, h := termbox.Size()
	ui.PrintCentered(strings.Repeat("Hello, termbox. ", 4))
	ui.DrawRectangle(0, 0, w-1, h-1, DefaultBoxStyle)
	termbox.Flush()
}

func (ui UI) MainLoop() {
	redraw := false
	ui.Draw()

mainLoop:
	for {
		redraw = false

		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey:
			switch event.Key {
			case termbox.KeyCtrlC, termbox.KeyEsc:
				break mainLoop
			}
		case termbox.EventResize:
			redraw = true
		case termbox.EventError:
			panic(event.Err)
		}

		if redraw {
			ui.Draw()
		}
	}
}
