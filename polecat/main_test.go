package main

import (
	"image"
	"reflect"
	"testing"
)

func TestReadOut(t *testing.T) {
	tests := []struct {
		read string
		want Output
	}{
		{
			read: "tests/dip.out",
			want: Output{
				Geom: []Atom{
					{"O", 0.014408608, -0.000000000, -1.564404376},
					{"O", -0.127428598, -0.000000000, 1.347805656},
					{"H", 0.897003803, 1.477943087, 1.719075321},
					{"H", 0.897003803, -1.477943087, 1.719075321},
				},
				max: 1.719075321,
			},
		},
	}
	for _, test := range tests {
		got := ReadOut(test.read)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("got %v, wanted %v\n", got, test.want)
		}
	}
}

func TestNormalizeGeom(t *testing.T) {
	out := ReadOut("tests/dip.out")
	out.NormalizeGeom()
	want := Output{
		Geom: []Atom{
			{"O", 0.008381603658656676, -0.0, -0.9100266619440347},
			{"O", -0.07412624475690441, -0.0, 0.784029436951006},
			{"H", 0.5217943577237824, 0.8597314317445198, 1.0},
			{"H", 0.5217943577237824, -0.8597314317445198, 1.0},
		},
		max: 1.719075321,
	}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("got %v, wanted %v\n", out, want)
	}
}

func TestCart2D(t *testing.T) {
	tests := []struct {
		msg     string
		x, y, z float64
		want    image.Point
	}{
		{"origin", 0, 0, 0, image.Point{width / 2, height / 2}},
		{"x", 1, 0, 0, image.Point{37, 219}},
		{"y", 0, 1, 0, image.Point{width, height / 2}},
		{"z", 0, 0, 1, image.Point{width / 2, 0}},
		{"-x", -1, 0, 0, image.Point{219, 37}},
		{"-y", 0, -1, 0, image.Point{0, height / 2}},
		{"-z", 0, 0, -1, image.Point{width / 2, height}},
		{"yz", 0, 1, 1, image.Point{width, 0}},
		{"y,-z", 0, 1, -1, image.Point{width, height}},
		{"xy", 1, 1, 0, image.Point{width, 0}},
	}
	for _, test := range tests {
		got := Cart2D(test.x, test.y, test.z)
		if got != test.want {
			t.Errorf("%s: got %v, wanted %v\n", test.msg, got, test.want)
		}
	}
}
