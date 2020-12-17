package spectro

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ntBre/chemutils/summarize"
)

var (
	names  = []string{"Al", "O", "O", "Al"}
	coords = `0.000000000        2.391678166        0.000000000
     -2.274263181        0.000000000        0.000000000
      2.274263181        0.000000000        0.000000000
      0.000000000       -2.391678166        0.000000000
`
)

func TestLoad(t *testing.T) {
	got, _ := Load("testfiles/spectro.in")
	got.FormatGeom([]string{"N", "C", "O", "H"},
		`0.0000000000       -0.0115666469        2.4598228639
      0.0000000000      -0.0139207809       0.2726915161
      0.0000000000       0.1184234620      -2.1785371074
      0.0000000000      -1.5591967852      -2.8818447886
`)
	want := &Spectro{
		Head: `# SPECTRO ##########################################
    1    1    3    2    0    0    1    4    0    1    0    0    0    0    0
    0    0    0    0    0    1    0    0    0    0    0    0    0    0    0
# GEOM #######################################
`,
		Geometry: `   4   1
 7.00      0.0000000000     -0.0115666469      2.4598228639
 6.00      0.0000000000     -0.0139207809      0.2726915161
 8.00      0.0000000000      0.1184234620     -2.1785371074
 1.00      0.0000000000     -1.5591967852     -2.8818447886
`,
		Body: `# WEIGHT ###### 
    4    
    1   14.003074
    2   12.0    
    3   15.9949146
    4    1.007825 
# CURVIL ##########################################
    1    2      
    2    3      
    3    4      
    3    2    4
    4    3    2    1
    4    3    2    1
`,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got\n%#v\nwanted\n%#v\n", got, want)
	}
}

func TestWriteSpectroInput(t *testing.T) {
	tests := []struct {
		load   string
		names  []string
		coords string
		write  string
		right  string
	}{
		{
			load:   "testfiles/spectro.in",
			names:  names,
			coords: coords,
			right:  "testfiles/right.in",
		},
	}
	for _, test := range tests {
		spec, _ := Load(test.load)
		spec.FormatGeom(test.names, test.coords)
		temp := t.TempDir()
		write := filepath.Join(temp, "spectro.in")
		spec.WriteInput(write)
		got, _ := ioutil.ReadFile(write)
		want, _ := ioutil.ReadFile(test.right)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got\n%s, wanted\n%s\n", got, want)
		}
	}
}

func TestReadSpectroOutput(t *testing.T) {
	tests := []struct {
		msg  string
		load string
		read string
		// assume other fields tested by Load
		fermi1 string
		fermi2 string
		polyad string
		coriol string
		darlin string
		nfreqs int
	}{
		{
			msg:  "all resonances present",
			load: "testfiles/spectro.in",
			read: "testfiles/spectro.out",
			fermi1: `    2
    4    2
    5    4
`,
			fermi2: `    1
    4    3    2
`,
			polyad: `    1
    5
    0    1    0    0    0    0
    0    0    0    1    0    0
    0    0    0    2    0    0
    0    0    0    0    2    0
    0    0    1    1    0    0
`,
			coriol: `    1
    0
    6    5    1    0    0
`,
			darlin: `    1
    6    5
`,
			nfreqs: 6,
		},
		{
			msg:  "no fermi 2 resonances present",
			load: "testfiles/spectro.in",
			read: "testfiles/spectro.prob",
			fermi1: `    1
    6    5
`,
			fermi2: "",
			polyad: "",
			coriol: `    3
    0
    3    2    0    0    1
    0
    4    1    0    0    1
    0
    5    4    0    0    1
`,
			darlin: `    6
    2    1
    3    1
    3    2
    4    2
    4    3
    5    4
`,
			nfreqs: 6,
		},
		{
			msg:    "no coriolis resonances present",
			load:   "testfiles/spectro.in",
			read:   "testfiles/spectro.nocoriol",
			fermi1: "",
			fermi2: "",
			polyad: "",
			coriol: "",
			darlin: `    1
    2    1
`,
			nfreqs: 3,
		},
		{
			msg:  "no fermi 2 but polyad",
			load: "testfiles/prob.in",
			read: "testfiles/prob.out",
			fermi1: `    4
    3    1
    7    3
    8    3
    9    3
`,
			fermi2: "",
			polyad: `    1
    6
    1    0    0    0    0    0    0    0    0
    0    0    1    0    0    0    0    0    0
    0    0    2    0    0    0    0    0    0
    0    0    0    0    0    0    2    0    0
    0    0    0    0    0    0    0    2    0
    0    0    0    0    0    0    0    0    2
`,
			coriol: `    7
    0
    6    5    1    0    0
    0
    7    6    1    0    0
    0
    8    5    0    0    1
    0
    8    6    0    1    0
    0
    8    7    0    0    1
    0
    9    7    0    1    0
    0
    9    8    1    0    0
`,
			darlin: `    7
    2    1
    6    5
    7    6
    8    6
    8    7
    9    7
    9    8
`,
			nfreqs: 9,
		},
	}
	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			spec, _ := Load(test.load)
			spec.ReadOutput(test.read)
			if spec.Fermi1 != test.fermi1 {
				t.Errorf("got %v, wanted %v\n", spec.Fermi1, test.fermi1)
			}
			if spec.Fermi2 != test.fermi2 {
				t.Errorf("got %v, wanted %v\n", spec.Fermi2, test.fermi2)
			}
			if !polyEqual(spec.Polyad, test.polyad) {
				t.Errorf("got %v, wanted %v\n", spec.Polyad, test.polyad)
			}
			if spec.Coriol != test.coriol {
				t.Errorf("got %v, wanted %v\n", spec.Coriol, test.coriol)
			}
			if spec.Darlin != test.darlin {
				t.Errorf("got\n%v, wanted\n%v\n", spec.Darlin, test.darlin)
			}
			if spec.Nfreqs != test.nfreqs {
				t.Errorf("got %v, wanted %v\n", spec.Nfreqs, test.nfreqs)
			}
		})
	}
}

// check if polyads are equal even though they come from unordered
// maps
func polyEqual(p1, p2 string) bool {
	if len(p1) != len(p2) {
		return false
	}
	if len(p1) == len(p2) && len(p1) == 0 {
		return true
	}
	sp1 := strings.Split(p1, "\n")
	sp2 := strings.Split(p2, "\n")
	if sp1[0] != sp2[0] || sp1[1] != sp2[1] {
		return false
	}
	sp1 = sp1[2:]
	sp2 = sp2[2:]
	var found bool
	for i := range sp1 {
		found = false
		for j := range sp2 {
			if sp1[i] == sp2[j] {
				found = true
				sp2 = append(sp2[:j], sp2[j+1:]...)
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestCheckPolyad(t *testing.T) {
	spec, err := Load("testfiles/spectro.in")
	if err != nil {
		t.Errorf("Load failed")
	}
	spec.ReadOutput("testfiles/spectro.out")
	got := spec.Polyad
	want := `    1
    5
    0    0    0    1    0    0
    0    1    0    0    0    0
    0    0    0    2    0    0
    0    0    0    0    2    0
    0    0    1    1    0    0
`
	if !polyEqual(got, want) {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

func TestMakeKey(t *testing.T) {
	got := MakeKey([]int{1, 2, 3})
	want := "1 2 3"
	if got != want {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

func TestResinLine(t *testing.T) {
	t.Run("One frequency on lhs", func(t *testing.T) {
		got := ResinLine(6, 2, 2)
		want := "    0    2    0    0    0    0\n"
		if got != want {
			t.Errorf("got %v, wanted %v\n", got, want)
		}
	})
	t.Run("two frequencies on lhs", func(t *testing.T) {
		got := ResinLine(6, 1, 2, 1)
		want := "    1    1    0    0    0    0\n"
		if got != want {
			t.Errorf("got %v, wanted %v\n", got, want)
		}
	})
}

func TestDoSpectro(t *testing.T) {
	SpectroCommand = "/home/brent/Downloads/spec3jm/backup/spectro.x"
	tests := []struct {
		load  string
		forts []string
		zpt   float64
		harm  []float64
		fund  []float64
		corr  []float64
	}{
		{
			load:  "testfiles/spectro.in",
			forts: []string{"fort.15", "fort.30", "fort.40"},
			zpt:   4682.7491,
			harm: []float64{
				3811.360, 2337.700, 1267.577,
				1086.351, 496.788, 437.756,
			},
			fund: []float64{
				3623.015, 2294.998, 1231.309,
				1071.641, 513.228, 454.579,
			},
			corr: []float64{
				3623.0149, 2298.5272, 1231.3094,
				1087.3762, 513.2276, 454.5787,
			},
		},
		{
			load:  "testfiles/prob.in",
			forts: []string{"prob.15", "prob.30", "prob.40"},
			zpt:   6974.5686,
			harm: []float64{
				3281.244, 3247.542, 1623.324,
				1307.596, 1090.695, 992.978,
				908.490, 901.527, 785.379,
			},
			fund: []float64{
				3140.115, 3113.252, 1589.153,
				1273.128, 1059.592, 967.523,
				887.009, 845.919, 769.511,
			},
			corr: []float64{
				3128.0166, 3113.2520, 1590.4451,
				1273.1281, 1059.5924, 967.5230,
				887.0092, 845.9192, 769.5109,
			},
		},
	}
	dests := []string{"fort.15", "fort.30", "fort.40"}
	for _, test := range tests {
		spec, _ := Load(test.load)
		tmp := t.TempDir()
		for i, file := range test.forts {
			src, _ := os.Open(filepath.Join("testfiles", file))
			dst, _ := os.Create(filepath.Join(tmp, dests[i]))
			io.Copy(dst, src)
		}
		spec.WriteInput(filepath.Join(tmp, "spectro.in"))
		spec.DoSpectro(tmp)
		res := summarize.Spectro(filepath.Join(tmp, "spectro2.out"))
		if res.ZPT != test.zpt {
			t.Errorf("got %v, wanted %v\n", res.ZPT, test.zpt)
		}
		if !reflect.DeepEqual(res.Harm, test.harm) {
			t.Errorf("got %v, wanted %v\n", res.Harm, test.harm)
		}
		if !reflect.DeepEqual(res.Fund, test.fund) {
			t.Errorf("got %v, wanted %v\n", res.Fund, test.fund)
		}
		if !reflect.DeepEqual(res.Corr, test.corr) {
			t.Errorf("got %v, wanted %v\n", res.Corr, test.corr)
		}
	}
}

func TestUpdateHeader(t *testing.T) {
	spec, _ := Load("testfiles/spectro.in")
	spec.ReadOutput("testfiles/spectro.out")
	spec.UpdateHeader()
	got := spec.Head
	want := `# SPECTRO ##########################################
    1    1    3    2    0    0    1    4    0    1    0    0    0    0    0
    1    4    4    0    0    1    0    0    0    0    0    1    0    0    0
# GEOM #######################################
`
	if got != want {
		t.Errorf("got\n%v, wanted\n%v\n", got, want)
	}
}
