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

// Distance returns the number of moves the Vector covers
func (v Vector) Distance() uint {
	a := IntAbs(v.x)
	b := IntAbs(v.y)
	if a > b {
		return a
	}
	return b
}

// Not really a unit vector, but shut up
func (v Vector) Unit() Vector {
	if v.Distance() <= 1 {
		return v
	}

	if v.x < 0 {
		v.x = -1
	} else if v.x > 0 {
		v.x = 1
	}

	if v.y < 0 {
		v.y = -1
	} else if v.y > 0 {
		v.y = 1
	}
	return v
}
