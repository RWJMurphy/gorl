package gorl

import "fmt"

// A Vec is a pair of x, y values.
type Vec struct {
	x, y int
}

func (v Vec) Plus(other Vec) Vec {
	return Vec{v.x + other.x, v.y + other.y}
}

func (c Vec) String() string {
	return fmt.Sprintf("<Vec x:%d, y:%d>", c.x, c.y)
}
