package gorl

import (
	"testing"
)

func TestDistance(t *testing.T) {
	tests := []struct {
		in   Vector
		want uint
	}{
		{Vector{0, 0}, 0},
		{Vector{0, 1}, 1},
		{Vector{1, 0}, 1},
		{Vector{1, 1}, 1},

		{Vector{0, -1}, 1},
		{Vector{-1, 0}, 1},
		{Vector{-1, -1}, 1},

		{Vector{1, -1}, 1},
		{Vector{-1, 1}, 1},

		{Vector{10, 20}, 20},
		{Vector{-20, -10}, 20},
	}
	for _, test := range tests {
		got := test.in.Distance()
		if got != test.want {
			t.Errorf("vec = %#v; vec.Distance() = %d, want %d", test.in, got, test.want)
		}
	}
}
func TestUnit(t *testing.T) {
	tests := []struct {
		in, want Vector
	}{
		{Vector{0, 0}, Vector{0, 0}},

		{Vector{0, 1}, Vector{0, 1}},
		{Vector{1, 0}, Vector{1, 0}},
		{Vector{1, 1}, Vector{1, 1}},

		{Vector{0, -1}, Vector{0, -1}},
		{Vector{-1, 0}, Vector{-1, 0}},
		{Vector{-1, -1}, Vector{-1, -1}},

		{Vector{1, -1}, Vector{1, -1}},
		{Vector{-1, 1}, Vector{-1, 1}},
		{Vector{-1, -1}, Vector{-1, -1}},

		{Vector{99, 99}, Vector{1, 1}},
		{Vector{99, 1}, Vector{1, 1}},
		{Vector{1, 99}, Vector{1, 1}},

		{Vector{-99, -99}, Vector{-1, -1}},
		{Vector{-99, 1}, Vector{-1, 1}},
		{Vector{1, -99}, Vector{1, -1}},
	}
	for _, test := range tests {
		got := test.in.Unit()
		if got != test.want {
			t.Errorf("vec = %#v; vec.Unit() = %#v, want %#v", test.in, got, test.want)
		}
	}
}
