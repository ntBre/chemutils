package main

import (
	"image"
	"math"
	"reflect"
	"sort"
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
				Dipx: 0.69892364,
				Dipy: -1.36097356e-06,
				Dipz: 1.65200259,
				max:  1.719075321,
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

func compAtom(a, b []Atom) bool {
	eps := 1e-4
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Symbol != b[i].Symbol {
			return false
		}
		ac, bc := a[i].Coords(), b[i].Coords()
		for c := range ac {
			if math.Abs(ac[c]-bc[c]) > eps {
				return false
			}
		}
	}
	return true
}

func TestNormalizeGeom(t *testing.T) {
	out := ReadOut("tests/dip.out")
	out.Normalize()
	want := Output{
		Geom: []Atom{
			{"O", 0.0041908, 0, -0.455013},
			{"O", -0.0370631, 0, 0.392015},
			{"H", 0.260897, 0.429866, 0.5},
			{"H", 0.260897, -0.429866, 0.5},
		},
		max: 1.719075321,
	}
	if !compAtom(out.Geom, want.Geom) {
		t.Errorf("got %v, wanted %v\n", out.Geom, want.Geom)
	}
	if math.Abs(out.max-want.max) > 1e-4 {
		t.Errorf("got %v, wanted %v\n", out.max, want.max)
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
		{"xy", 1, 1, 0, image.Point{165, 219}},
	}
	for _, test := range tests {
		got := Cart2D(Vec{test.x, test.y, test.z})
		if got != test.want {
			t.Errorf("%s: got %v, wanted %v\n", test.msg,
				got, test.want)
		}
	}
}

func TestOrder(t *testing.T) {
	tests := []struct {
		v    Vec
		want []Axis
	}{
		{
			v:    Vec{1, 2, 3},
			want: []Axis{Z, Y, X},
		},
	}
	for _, test := range tests {
		got := test.v.Order()
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("got %v, wanted %v\n", got, test.want)
		}
	}
}

func compareFloats(a, b []float64, eps float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > eps {
			return false
		}
	}
	return true
}

func TestMOI(t *testing.T) {
	eps := 1e-6
	atoms := []Atom{
		{"C", 0.000000000000, 0.000000000000, 0.000000000000},
		{"C", 0.000000000000, 0.000000000000, 2.845112131228},
		{"O", 1.899115961744, 0.000000000000, 4.139062527233},
		{"H", -1.894048308506, 0.000000000000, 3.747688672216},
		{"H", 1.942500819960, 0.000000000000, -0.701145981971},
		{"H", -1.007295466862, -1.669971842687, -0.705916966833},
		{"H", -1.007295466862, 1.669971842687, -0.705916966833},
	}
	com := COM(atoms)
	for a := range atoms {
		atoms[a] = atoms[a].Translate(com)
	}
	ia, ib, ic := MOI(atoms)
	got := []float64{ia, ib, ic}
	sort.Float64s(got)
	want := []float64{31.964078, 178.649562, 199.371127}
	if !compareFloats(got, want, eps) {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

func TestRot(t *testing.T) {
	eps := 1e-4
	atoms := []Atom{
		{"C", 0.000000000000, 0.000000000000, 0.000000000000},
		{"C", 0.000000000000, 0.000000000000, 2.845112131228},
		{"O", 1.899115961744, 0.000000000000, 4.139062527233},
		{"H", -1.894048308506, 0.000000000000, 3.747688672216},
		{"H", 1.942500819960, 0.000000000000, -0.701145981971},
		{"H", -1.007295466862, -1.669971842687, -0.705916966833},
		{"H", -1.007295466862, 1.669971842687, -0.705916966833},
	}
	com := COM(atoms)
	for a := range atoms {
		atoms[a] = atoms[a].Translate(com)
	}
	ia, ib, ic := MOI(atoms)
	tests := []struct {
		I    float64
		want float64
	}{
		{ib, 1.8834},
		{ia, 0.3370},
		{ic, 0.3019},
	}
	for _, test := range tests {
		got := Rot(test.I)
		if math.Abs(got-test.want) > eps {
			t.Errorf("got %v, wanted %v\n", got, test.want)
		}
	}
}
