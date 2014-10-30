package gorl

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

type Feature interface {
	Loc() Coord
	Name() string
	Char() rune
	Color() termbox.Attribute
	Flags() Flag
	LightRadius() int
}

type feature struct {
	loc          Coord
	name         string
	char         rune
	color        termbox.Attribute
	flags        Flag
	lightRadius  int
}

func NewFeature(name string, char rune) *feature {
	f := &feature{}
	f.name = name
	f.char = char
	f.color = termbox.ColorDefault
	return f
}

func (f *feature) Loc() Coord {
	return f.loc
}

func (f *feature) Name() string {
	return f.name
}

func (f *feature) Char() rune {
	return f.char
}

func (f *feature) Color() termbox.Attribute {
	return f.color
}

func (f *feature) Flags() Flag {
	return f.flags
}

func (f *feature) LightRadius() int {
	return f.lightRadius
}

func (f *feature) String() string {
	return fmt.Sprintf(
		"<feature %s char:%c, loc:%s, flags:%s, lightRadius:%d>",
		f.name,
		f.char,
		f.loc,
		f.flags,
		f.lightRadius,
	)
}
