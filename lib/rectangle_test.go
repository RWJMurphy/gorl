package gorl

import "testing"

//type RectangleI interface {
//TopLeft() Vector
//BottomRight() Vector
//TopRight() Vector
//BottomLeft() Vector
//Size() Vector
//Width() int
//Height() int
//}

func TestRectangle(t *testing.T) {
	tests := []struct {
		rect                                       RectangleI
		topLeft, topRight, bottomLeft, bottomRight Vector
		size                                       Vector
		width, height                              int
	}{
		{
			Rectangle{Vector{0, 0}, Vector{0, 0}},
			Vector{0, 0},
			Vector{0, 0},
			Vector{0, 0},
			Vector{0, 0},
			Vector{0, 0},
			0,
			0,
		},
		{
			Rectangle{Vector{99, 99}, Vector{1, 1}},
			Vector{99, 99},
			Vector{100, 99},
			Vector{99, 100},
			Vector{100, 100},
			Vector{1, 1},
			1,
			1,
		},
		{
			Rectangle{Vector{10, 20}, Vector{30, 40}},
			Vector{10, 20},
			Vector{40, 20},
			Vector{10, 60},
			Vector{40, 60},
			Vector{30, 40},
			30,
			40,
		},
	}
	for _, test := range tests {
		if test.rect.TopLeft() != test.topLeft {
			t.Errorf("%#v.TopLeft() = %#v, want %#v",
				test.rect,
				test.rect.TopLeft(),
				test.topLeft)
		}
		if test.rect.TopRight() != test.topRight {
			t.Errorf("%#v.TopRight() = %#v, want %#v",
				test.rect,
				test.rect.TopRight(),
				test.topRight)
		}
		if test.rect.BottomLeft() != test.bottomLeft {
			t.Errorf("%#v.BottomLeft() = %#v, want %#v",
				test.rect,
				test.rect.BottomLeft(),
				test.bottomLeft)
		}
		if test.rect.BottomRight() != test.bottomRight {
			t.Errorf("%#v.BottomRight() = %#v, want %#v",
				test.rect,
				test.rect.BottomRight(),
				test.bottomRight)
		}
		if test.rect.Size() != test.size {
			t.Errorf("%#v.Size() = %#v, want %#v",
				test.rect,
				test.rect.Size(),
				test.size)
		}
		if test.rect.Width() != test.width {
			t.Errorf("%#v.Width() = %#v, want %#v",
				test.rect,
				test.rect.Width(),
				test.width)
		}
		if test.rect.Height() != test.height {
			t.Errorf("%#v.Height() = %#v, want %#v",
				test.rect,
				test.rect.Height(),
				test.height)
		}
	}
}
