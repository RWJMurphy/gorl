package gorl

import "fmt"

type RectangleI interface {
	TopLeft() Vec
	BottomRight() Vec
	TopRight() Vec
	BottomLeft() Vec
	Size() Vec
	Width() int
	Height() int
}

type Rectangle struct {
	topLeft Vec
	size    Vec
}

func (r Rectangle) TopLeft() Vec {
	return r.topLeft
}

func (r Rectangle) BottomRight() Vec {
	return r.topLeft.Plus(r.size).Plus(Vec{-1, -1})
}

func (r Rectangle) TopRight() Vec {
	return Vec{r.BottomRight().x, r.TopLeft().y}
}

func (r Rectangle) BottomLeft() Vec {
	return Vec{r.TopLeft().x, r.BottomRight().y}
}

func (r Rectangle) Width() int {
	return r.size.x
}

func (r Rectangle) Height() int {
	return r.size.y
}

func (r Rectangle) Size() Vec {
	return r.size
}

func (r Rectangle) String() string {
	return fmt.Sprintf("<Rectangle topLeft:%s, bottomRight:%s, size:%s>", r.topLeft, r.BottomRight(), r.size)
}
