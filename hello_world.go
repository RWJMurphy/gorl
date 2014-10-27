package main

import (
	"github.com/nsf/termbox-go"
	"strings"
	"unicode/utf8"
)

type UI struct{}

func (ui UI) Close() {
	termbox.Close()
}

func NewUI() (*UI, error) {
	ui := &UI{}
	err := termbox.Init()
	return ui, err
}

func (ui UI) PrintAt(x, y int, s string) {
	for i, r := range s {
		termbox.SetCell(x+i, y, r, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func (ui UI) PrintCentered(s string) {
	width, height := termbox.Size()
	mid_w, mid_h := width/2, height/2
	s_len := utf8.RuneCountInString(s)
	ui.PrintAt(mid_w-s_len/2, mid_h, s)
}

func (ui UI) Draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	ui.PrintCentered(strings.Repeat("Hello, termbox. ", 4))
	termbox.Flush()
}

func main() {
	redraw := false
	ui, err := NewUI()
	if err != nil {
		panic(err)
	}
	defer ui.Close()
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

