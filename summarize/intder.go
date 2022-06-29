package summarize

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	bohr2ang = 0.52917721067
)

type Atom struct {
	Sym     string
	X, Y, Z float64
}

func (a Atom) String() string {
	return fmt.Sprintf(
		"%2s%15.10f%15.10f%15.10f",
		a.Sym, a.X, a.Y, a.Z,
	)
}

// Intder is a struct for holding the information in a frequency
// intder output file. Geom is the geometry, SiIC is the simple
// internal coordinate system, SyIC is the symmetry internal
// coordinate system, and Vibs are the vibrational assignments in
// terms of the SyICs.
type Intder struct {
	Geom []Atom
	Dumm []Atom
	SiIC [][]int
	SyIC [][]int
	Freq []float64
	Vibs string
}

// ptable is a map from the default string masses in intder to the
// corresponding atomic symbols
var ptable = map[string]string{
	"1.007825":  "H",
	"4.002600":  "He",
	"11.009310": "B",
	"12.000000": "C",
	"14.003070": "N",
	"15.994910": "O",
	"18.998400": "F",
	"19.992435": "Ne",
	"21.991383": "Ne",
	"26.981530": "Al",
	"27.976927": "Si",
	"31.972070": "S",
	"34.968852": "Cl",
	"35.967545": "Ar",
	"37.962732": "Ar",
	"39.962384": "Ar",
}

const (
	STRE int = iota
	BEND
	TORS
	LIN1
	LINX
	LINY
	OUT
)

func ReadIntder(filename string) *Intder {
	id := new(Intder)
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	var (
		geom   bool
		siic   bool
		syic   bool
		vibs   bool
		dumm   bool
		skip   uint
		line   string
		fields []string
	)
	var str strings.Builder
	contrib := regexp.MustCompile(`(-?)([0-9]+) \([ -]*([0-9]{1,3}\.[0-9])\)`)
	sint := regexp.MustCompile(`(L|S)\( ?([0-9]{1,2})\)=?`)
	syicpat := regexp.MustCompile(`^ SYMMETRY INTERNAL COORDINATES$`)
	// the length of the last SIC; if 4 components, intder prints a blank
	// line that will otherwise terminate the SIC section
	var lsic_last int
	for scanner.Scan() {
		line = scanner.Text()
		fields = strings.Fields(line)
		switch {
		case skip > 0:
			skip--
		case strings.Contains(line,
			"NUCLEAR CARTESIAN COORDINATES (ANG.)"):
			skip += 3
			geom = true
			id.Geom = nil
		case strings.Contains(line,
			"DUMMY ATOM VECTORS"):
			skip++
			dumm = true
		case strings.Contains(line,
			"VIBRATIONAL ASSIGNMENTS"):
			skip += 4
			vibs = true
		case strings.Contains(line,
			"DEFINITION OF INTERNAL COORDINATES"):
			skip += 3
			siic = true
		case geom && len(fields) == 0:
			geom = false
		case dumm && len(fields) == 0:
			dumm = false
		case vibs && len(fields) == 0:
			vibs = false
			id.Vibs = str.String()
			str.Reset()
		case siic && len(fields) == 0:
			siic = false
		case syicpat.MatchString(line):
			syic = true
			skip += 1
		case syic && len(fields) == 0:
			if lsic_last != 4 {
				syic = false
			}
		case strings.Contains(line,
			"VALUES OF SIMPLE INTERNAL COORDINATES (ANG. OR DEG.)",
		):
			syic = false
		case geom:
			coords := make([]float64, 3)
			for i, v := range fields[2:] {
				coords[i], _ = strconv.ParseFloat(v, 64)
			}
			id.Geom = append(id.Geom, Atom{
				ptable[fields[1]],
				coords[0],
				coords[1],
				coords[2],
			})
		case dumm:
			coords := make([]float64, 3)
			for i, v := range fields {
				coords[i], _ = strconv.ParseFloat(v, 64)
			}
			id.Dumm = append(id.Dumm, Atom{
				"X",
				coords[0] * bohr2ang,
				coords[1] * bohr2ang,
				coords[2] * bohr2ang,
			})
		case vibs:
			v, _ := strconv.ParseFloat(fields[1], 64)
			id.Freq = append(id.Freq, v)
			contribs := contrib.ReplaceAllString(
				strings.Join(fields[2:], " "),
				"${1}${3} S_{${2}}")
			cf := strings.Fields(contribs)
			var ret strings.Builder
			for i, c := range cf {
				if i%2 == 0 {
					v, _ := strconv.ParseFloat(c, 64)
					if i > 0 {
						fmt.Fprintf(&ret, "%+5.3f", v/100)
					} else {
						fmt.Fprintf(&ret, "%5.3f", v/100)
					}
				} else {
					fmt.Fprintf(&ret, "%s", c)
				}
			}
			fmt.Fprintln(&str, ret.String())
		case siic:
			line = sint.ReplaceAllString(line, "")
			fields = strings.Fields(line)
			ids := make([]int, 4)
			for i, f := range fields[1:] {
				ids[i], _ = strconv.Atoi(f)
				ids[i]-- // index from zero
			}
			switch fields[0] {
			case "STRE":
				ids = append(ids, STRE)
			case "BEND":
				ids = append(ids, BEND)
			case "TORS":
				ids = append(ids, TORS)
			case "LIN1":
				ids = append(ids, LIN1)
			case "OUT":
				ids = append(ids, OUT)
			case "LINX":
				ids = append(ids, LINX)
			case "LINY":
				ids = append(ids, LINY)
			default:
				panic("this type of coordinate not implemented")
			}
			id.SiIC = append(id.SiIC, ids)
		case syic:
			lsic_last = 0
			ids := make([]int, 0)
			fields = strings.Fields(
				sint.ReplaceAllString(line, "${2}"),
			)
			var fac int = 1
			for i, f := range fields[1:] {
				if i%2 == 0 {
					v, _ := strconv.ParseFloat(f, 64)
					if math.Signbit(v) {
						fac = -1
					} else {
						fac = 1
					}
				} else {
					lsic_last++
					d, _ := strconv.Atoi(f)
					ids = append(ids, fac*(d-1))
				}
			}
			id.SyIC = append(id.SyIC, ids)
		}
	}
	return id
}
