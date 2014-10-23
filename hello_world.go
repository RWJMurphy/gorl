package main

import (
    "strings"
    "unicode/utf8"
    "github.com/nsf/termbox-go"
)

func print_at(x, y int, s string) {
    for i, r := range s {
        termbox.SetCell(x + i, y, r, termbox.ColorDefault, termbox.ColorDefault)
    }
}

func print_centered(s string) {
    width, height := termbox.Size()
    mid_w, mid_h := width / 2, height / 2
    s_len := utf8.RuneCountInString(s)
    print_at(mid_w - s_len / 2, mid_h, s)
}

func draw() {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
    print_centered(strings.Repeat("Hello, termbox. ", 4))
    termbox.Flush()
}

func main() {
    err := termbox.Init()
    if err != nil {
        panic(err)
    }
    defer termbox.Close()

    draw()
loop:
    for {
        event := termbox.PollEvent()
        switch event.Type {
        case termbox.EventKey:
            switch event.Key {
            case termbox.KeyCtrlC, termbox.KeyEsc:
                break loop
            }
        case termbox.EventResize:
            draw()
        case termbox.EventError:
            panic(event.Err)
        }
    }
}
