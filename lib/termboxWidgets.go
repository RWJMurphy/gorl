package gorl

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type widget struct {
	Rectangle
	ui TermboxUI
}

// Paint paints the Widget to the UI
func (w *widget) Paint() {
	w.ui.PaintBorder(w, defaultBoxStyle)
}

type cameraWidget struct {
	widget
	dungeon *Dungeon
	center  Vector
}

// Paint paints the cameraWidget to the TermboxUI
func (camera *cameraWidget) Paint() {
	var (
		tile   *Tile
		offset Vector
		loc    Vector
		out    Vector
		x, y   int
		char   rune
		color  termbox.Attribute
	)

	ne := camera.center.Add(Vector{-camera.widget.Width() / 2, -camera.widget.Height() / 2})

	for x = 0; x < camera.widget.Width(); x++ {
		for y = 0; y < camera.widget.Height(); y++ {
			offset = Vector{x, y}
			loc = ne.Add(offset)
			out = camera.TopLeft().Add(offset)
			tile = camera.dungeon.Tile(loc)
			if tile.Seen() || tile.Visible() {
				if tile.Visible() {
					fg := camera.dungeon.FeatureGroup(loc)
					if fg.mob != nil {
						char = fg.mob.Char()
						color = fg.mob.Color()
					} else if fg.feature != nil {
						char = fg.feature.Char()
						color = fg.feature.Color()
					} else if len(fg.items) > 0 {
						char = fg.items[len(fg.items)-1].Char()
						color = fg.items[len(fg.items)-1].Color()
					} else {
						char = tile.c
						color = tile.color | termbox.AttrBold
					}
					camera.ui.PutRuneColor(out, char, color, termbox.ColorDefault)
				} else {
					camera.ui.PutRuneColor(out, tile.c, tile.color, termbox.ColorDefault)
				}
			}
		}
	}
	camera.widget.Paint()
}

type logWidget struct {
	widget
	messages []string
}

// Paint paints the logWidget to the TermboxUI
func (lw *logWidget) Paint() {
	var loc Vector
	for i, m := range lw.ui.Messages() {
		loc = lw.TopLeft().Add(Vector{1, 1 + i})
		lw.ui.PrintAt(loc, m)
	}
	lw.widget.Paint()
}

// A menuWidget in theory displays a menu.
type menuWidget struct {
	widget
}

// Paint paints the MenuWidget to the UI
func (mw *menuWidget) Paint() {
	mw.widget.Paint()
}

type inventoryWidget struct {
	widget
	owner Mob
}

func (iw *inventoryWidget) SetOwner(m Mob) {
	iw.owner = m
}

func (iw *inventoryWidget) Paint() {
	var loc Vector
	iw.ui.PrintAt(
		iw.TopLeft().Add(Vector{1, 1}),
		"Inventory",
	)
	for i, item := range iw.owner.Inventory() {
		loc = iw.TopLeft().Add(Vector{1, 3 + i})
		iw.ui.PrintAt(loc, fmt.Sprintf("%c) %s", 'a' + i, item.Name()))
	}
	iw.widget.Paint()
}
