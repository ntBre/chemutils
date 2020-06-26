package atom

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
	got := ReadXYZ(strings.NewReader(xyz), true)
	want := []Atom{
		{"H", []float64{0.0000000000, 0.7574590974, 0.5217905143}},
		{"O", []float64{0.0000000000, 0.0000000000, -0.0657441568}},
		{"H", []float64{0.0000000000, -0.7574590974, 0.5217905143}},
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
	got := ToCartesian([]float64{0.7574590974, math.Pi / 2, 0.5217905143})
	want := []float64{0.0000000000, 0.7574590974, 0.5217905143}
	if !approxEqual(got, want) {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

func TestRotate(t *testing.T) {
	xyz := `3
water
H 0.0000000000 0.7574590974 0.5217905143
O 0.0000000000 0.0000000000 -0.0657441568
H 0.0000000000 -0.7574590974 0.5217905143
`
	atoms := ReadXYZ(strings.NewReader(xyz), true)
	tests := []struct {
		deg  float64
		axis int
		want []Atom
	}{
		{180, Z,
			[]Atom{
				{"H", []float64{0.0000000000, -0.7574590974, 0.5217905143}},
				{"O", []float64{0.0000000000, 0.0000000000, -0.0657441568}},
				{"H", []float64{0.0000000000, 0.7574590974, 0.5217905143}},
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

	for i := range tests {
		got := Rotate(atoms, tests[i].deg, tests[i].axis)
		if !reflect.DeepEqual(got, tests[i].want) {
			t.Errorf("got %v, wanted %v\n", got, tests[i].want)
		}
	}
}
