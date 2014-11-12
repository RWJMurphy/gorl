package gorl

import "fmt"

// A Vector is a pair of x, y values.
type Vector struct {
	x, y int
}

func (v Vector) String() string {
	return fmt.Sprintf("<Vector x:%d, y:%d>", v.x, v.y)
}

func (v Vector) Add(other Vector) Vector {
	return Vector{v.x + other.x, v.y + other.y}
}

func (v Vector) Sub(other Vector) Vector {
	return Vector{v.x - other.x, v.y - other.y}
}
