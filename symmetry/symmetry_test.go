package symm

import (
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestReadXYZ(t *testing.T) {
	xyz := `3
water
H 0.0000000000 0.7574590974 0.5217905143
O 0.0000000000 0.0000000000 -0.0657441568
H 0.0000000000 -0.7574590974 0.5217905143
`
	got := ReadXYZ(strings.NewReader(xyz))
	want := Molecule{
		Atoms: []Atom{
			{"H", []float64{
				0.0000000000,
				0.7574590974,
				0.5217905143,
			}},
			{"O", []float64{
				0.0000000000,
				0.0000000000,
				-0.0657441568,
			}},
			{"H", []float64{
				0.0000000000,
				-0.7574590974,
				0.5217905143,
			}},
		},
		Principal: Z,
		Main:      Plane{Y, Z},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

func TestToCylinder(t *testing.T) {
	got := ToCylinder([]float64{0.0000000000, 0.7574590974, 0.5217905143}, Z)
	want := []float64{0.7574590974, math.Pi / 2, 0.5217905143}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

func TestToCartesian(t *testing.T) {
	got := ToCartesian([]float64{0.7574590974, math.Pi / 2, 0.5217905143}, Z)
	want := []float64{0.0000000000, 0.7574590974, 0.5217905143}
	if !approxEqual(got, want) {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

func TestRotate(t *testing.T) {
	mole := LoadXYZ("tests/h2o.xyz")
	tests := []struct {
		deg  float64
		axis Axis
		want []Atom
	}{
		{180, Z,
			[]Atom{
				{"H", []float64{0.0000000000, -0.7574590974, 0.5217905143}},
				{"O", []float64{0.0000000000, 0.0000000000, -0.0657441568}},
				{"H", []float64{0.0000000000, 0.7574590974, 0.5217905143}},
			},
		},
		{180, Y,
			[]Atom{
				{"H", []float64{0.0000000000, 0.7574590974, -0.5217905143}},
				{"O", []float64{0.0000000000, 0.0000000000, 0.0657441568}},
				{"H", []float64{0.0000000000, -0.7574590974, -0.5217905143}},
			},
		},
		{180, X,
			[]Atom{
				{"H", []float64{0.0000000000, -0.7574590974, -0.5217905143}},
				{"O", []float64{0.0000000000, 0.0000000000, 0.0657441568}},
				{"H", []float64{0.0000000000, 0.7574590974, -0.5217905143}},
			},
		},
		{90, Z,
			[]Atom{
				{"H", []float64{-0.7574590974, 0.0000000000, 0.5217905143}},
				{"O", []float64{0.0000000000, 0.0000000000, -0.0657441568}},
				{"H", []float64{0.7574590974, 0.0000000000, 0.5217905143}},
			},
		},
	}
	for _, test := range tests {
		got := Rotate(mole.Atoms, test.deg, test.axis)
		for j := range got {
			if got[j].Label != test.want[j].Label ||
				!approxEqual(got[j].Coord, test.want[j].Coord) {
				t.Errorf("Rotate(%f, %s): got %v, wanted %v\n",
					test.deg, test.axis, got, test.want)
			}
		}
	}
}

func TestReflect(t *testing.T) {
	xyz := `3
water
H 0.0000000000 0.7574590974 0.5217905143
O 0.0000000000 0.0000000000 -0.0657441568
H 0.0000000000 -0.7574590974 0.5217905143
`
	atoms := ReadXYZ(strings.NewReader(xyz))
	tests := []struct {
		plane Plane
		want  []Atom
	}{
		{
			Plane{Y, Z},
			[]Atom{
				{"H", []float64{0.0000000000, 0.7574590974, 0.5217905143}},
				{"O", []float64{0.0000000000, 0.0000000000, -0.0657441568}},
				{"H", []float64{0.0000000000, -0.7574590974, 0.5217905143}},
			},
		},
		{
			Plane{Y, X},
			[]Atom{
				{"H", []float64{0.0000000000, 0.7574590974, -0.5217905143}},
				{"O", []float64{0.0000000000, 0.0000000000, 0.0657441568}},
				{"H", []float64{0.0000000000, -0.7574590974, -0.5217905143}},
			},
		},
		{
			Plane{X, Z},
			[]Atom{
				{"H", []float64{
					0.0000000000,
					-0.7574590974,
					0.5217905143,
				}},
				{"O", []float64{
					0.0000000000,
					0.0000000000,
					-0.0657441568,
				}},
				{"H", []float64{
					0.0000000000,
					0.7574590974,
					0.5217905143,
				}},
			},
		},
	}
	for i := range tests {
		got := Reflect(atoms.Atoms, tests[i].plane)
		for j := range got {
			if got[j].Label != tests[i].want[j].Label ||
				!approxEqual(got[j].Coord, tests[i].want[j].Coord) {
				t.Errorf("Reflect(%s): got\n%v, wanted\n%v\n",
					tests[i].plane, got, tests[i].want)
				break
			}
		}
	}
}

func TestRotaryReflect(t *testing.T) {
	tests := []struct {
		atoms string
		deg   float64
		axis  Axis
	}{
		{
			atoms: "tests/ethane.xyz",
			deg:   60.0,
			axis:  Z,
		},
	}
	tmp := eps
	eps = 1e-11
	defer func() {
		eps = tmp
	}()
	var found bool
	for _, test := range tests {
		// want it to give itself back
		wants := LoadXYZ(test.atoms)
		gots := RotaryReflect(wants.Atoms, test.deg, test.axis)
		for _, got := range gots {
			found = false
			for _, want := range wants.Atoms {
				if got.Label == want.Label &&
					approxEqual(got.Coord, want.Coord) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("got\n%v,\nwanted\n%v\n", gots, wants)
				break
			}
		}
	}
}

func TestInvert(t *testing.T) {
	tests := []struct {
		atoms string
		axis  Axis
	}{
		{
			atoms: "tests/ethane.xyz",
			axis:  Z,
		},
	}
	tmp := eps
	eps = 1e-11
	defer func() {
		eps = tmp
	}()
	var found bool
	for _, test := range tests {
		// want it to give itself back
		wants := LoadXYZ(test.atoms)
		gots := Invert(wants.Atoms, test.axis)
		for _, got := range gots {
			found = false
			for _, want := range wants.Atoms {
				if got.Label == want.Label &&
					approxEqual(got.Coord, want.Coord) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("got\n%v,\nwanted\n%v\n", gots, wants)
				break
			}
		}
	}
}
