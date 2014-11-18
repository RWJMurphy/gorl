package gorl

import (
	"testing"
)

func TestIntAbs(t *testing.T) {
	tests := []struct {
		in   int
		want uint
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{9999, 9999},

		{-1, 1},
		{-2, 2},
		{-9999, 9999},
	}
	for _, test := range tests {
		got := IntAbs(test.in)
		if got != test.want {
			t.Errorf("IntAbs(%d) = %d, want %d", test.in, got, test.want)
		}
	}
}
