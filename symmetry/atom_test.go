package symm

import "testing"

func TestSymmetry(t *testing.T) {
	tests := []struct {
		msg   string
		atoms []Atom
		want  Irrep
	}{
		{
			msg: "symmetric stretch",
			atoms: []Atom{
				{"H", []float64{0.0000000000, 1.4366694195, 0.9874061512}},
				{"O", []float64{0.0000000000, -0.0000000000, -0.1269683874}},
				{"H", []float64{0.0000000000, -1.4366694195, 0.9874061512}},
			},
			want: A1,
		},
		{
			msg: "anti-symmetric stretch",
			atoms: []Atom{
				{"H", []float64{0.0000000000, 1.4337425639, 0.9878589402}},
				{"O", []float64{0.0000000000, -0.0046953159, -0.1242319281}},
				{"H", []float64{0.0000000000, -1.4290472480, 0.9842169028}},
			},
			want: B2,
		},
		{
			msg: "bend",
			atoms: []Atom{
				{"H", []float64{0.0000000000, 1.4341614301, 0.9848472035}},
				{"H", []float64{0.0000000000, -0.0000000000, -0.1218504921}},
				{"H", []float64{0.0000000000, -1.4341614301, 0.9848472035}},
			},
			want: A1,
		},
	}
	mol := LoadXYZ("tests/intder_h2o.xyz")
	// TODO uncomment for full range
	for _, test := range tests[1:2] {
		got := mol.Symmetry(test.atoms)
		if got != test.want {
			t.Errorf("%s: got %v, wanted %v\n",
				test.msg, got, test.want)
		}
	}
}

func TestIsSame(t *testing.T) {
	tests := []struct {
		atoms []Atom
		btoms []Atom
		want  bool
	}{
		{
			atoms: []Atom{
				{"H", []float64{0.000000, 1.433743, 0.987859}},
				{"O", []float64{0.000000, -0.004695, -0.124232}},
				{"H", []float64{0.000000, -1.429047, 0.984217}},
			},
			btoms: []Atom{
				{"H", []float64{0.000000, -1.433743, 0.987859}},
				{"O", []float64{0.000000, 0.004695, -0.124232}},
				{"H", []float64{0.000000, 1.429047, 0.984217}},
			},
			want: false,
		},
	}
	for _, test := range tests {
		got := IsSame(test.atoms, test.btoms)
		if got != test.want {
			t.Errorf("got %v, wanted %v\n", got, test.want)
		}
	}
}
