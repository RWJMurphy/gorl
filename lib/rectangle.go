package gorl

import "fmt"

type RectangleI interface {
	TopLeft() Vector
	BottomRight() Vector
	TopRight() Vector
	BottomLeft() Vector
	Size() Vector
	Width() int
	Height() int
}

type Rectangle struct {
	topLeft Vector
	size    Vector
}

func (r Rectangle) TopLeft() Vector {
	return r.topLeft
}

func (r Rectangle) BottomRight() Vector {
	return r.topLeft.Add(r.size)
}

func (r Rectangle) TopRight() Vector {
	return Vector{r.BottomRight().x, r.TopLeft().y}
}

func (r Rectangle) BottomLeft() Vector {
	return Vector{r.TopLeft().x, r.BottomRight().y}
}

func (r Rectangle) Width() int {
	return r.size.x
}

func (r Rectangle) Height() int {
	return r.size.y
}

func (r Rectangle) Size() Vector {
	return r.size
}

func (r Rectangle) String() string {
	return fmt.Sprintf("<Rectangle topLeft:%s, bottomRight:%s, size:%s>", r.topLeft, r.BottomRight(), r.size)
}
