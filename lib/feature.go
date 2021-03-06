package gorl

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

// Feature is a visible, rendered part of a Dungeon
type Feature interface {
	Name() string
	Char() rune

	Color() termbox.Attribute
	SetColor(termbox.Attribute)

	Loc() Vector
	SetLoc(Vector)

	Flags() Flag

	LightRadius() int
	SetLightRadius(int)
}

type feature struct {
	loc         Vector
	name        string
	char        rune
	color       termbox.Attribute
	flags       Flag
	lightRadius int
}

// NewFeature returns a new Feature
func NewFeature(name string, char rune) Feature {
	f := &feature{}
	f.name = name
	f.char = char
	f.color = termbox.ColorDefault
	return f
}

func (f *feature) Loc() Vector {
	return f.loc
}

func (f *feature) SetLoc(loc Vector) {
	f.loc = loc
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

func (f *feature) SetColor(color termbox.Attribute) {
	f.color = color
}

func (f *feature) Flags() Flag {
	return f.flags
}

func (f *feature) LightRadius() int {
	return f.lightRadius
}

func (f *feature) SetLightRadius(radius int) {
	f.lightRadius = radius
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
