// Package summarize is for summarizing output from quantum chemistry
// programs
// (let ((compile-command "go test .")) (my-recompile))
package summarize

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Result struct {
	ZPT    float64
	Harm   []float64
	Fund   []float64
	Corr   []float64
	Rots   [][]float64
	Deltas []float64
	Phis   []float64
	Rhead  []string
	Ralpha []float64
	Requil []float64
	Fermi  []string
}

// FreqReport gathers harmonic, anharmonic, and resonance-corrected
// frequencies from a spectro  output file for reporting
func Spectro(filename string) *Result {
	res := new(Result)
	fermiMap := make(map[string][]string)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var (
		line     string
		fields   []string
		skip     int
		corr     bool
		holdFreq float64
		holdZPT  float64
		pos      int
		state    []string
		good     bool
		one      bool
		harmFund bool
		rot      bool
		geom     bool
		fermi1   bool
		fermi2   bool
		nobza    bool = true
		buf      bytes.Buffer
		// this is cute
		gparams = []string{"", "", "r", "<"}
	)
	freq := regexp.MustCompile(`[0-9]+`)
	delta := regexp.MustCompile(`(?i)delta (J|(JK)|K)`)
	phi := regexp.MustCompile(`(?i)phi (J|(JK)|(KJ)|K)`)
	icn := regexp.MustCompile(`\([ 0-9]+\)\s+(BOND|ANGLE)`)
	atom := regexp.MustCompile(`([0-9]+)\(([A-Za-z ]+)\)`)
	for scanner.Scan() {
		line = scanner.Text()
		fields = strings.Fields(line)
		switch {
		case skip > 0:
			skip--
		case strings.Contains(line, "BAND CENTER ANALYSIS"):
			skip += 3
			harmFund = true
		case harmFund && len(line) > 1:
			if strings.Contains(line, "DUNHAM") ||
				strings.Contains(line, "VIBRATIONAL ENERGY AND") {
				harmFund = false
			} else if freq.MatchString(fields[0]) && len(fields) > 2 {
				h, _ := strconv.ParseFloat(fields[1], 64)
				f, _ := strconv.ParseFloat(fields[2], 64)
				res.Harm = append(res.Harm, h)
				res.Fund = append(res.Fund, f)
			}
		case strings.Contains(line, "STATE NO.") &&
			!strings.Contains(line, "SPECTRUM"):
			skip += 2
			corr = true
		case corr && strings.Contains(line, "*******************"):
			corr = false
			good = true
		case corr && strings.Contains(line, "NON-DEG (Vs)"):
			// need to grab freq, then check if the next line
			// contains more state info before appending
			holdZPT, _ = strconv.ParseFloat(fields[1], 64)
			holdFreq, _ = strconv.ParseFloat(fields[2], 64)
			pos, _ = strconv.Atoi(fields[0])
			state = nil
			good = true
			one = false
			for _, f := range fields[6:] {
				// if there's a 2 or a second one it's bad
				if f == "2" || (one && f == "1") {
					good = false
					break
				} else if f == "1" {
					one = true
				}
				state = append(state, f)
			}
		case corr && good && strings.Contains(line, "DEGEN   (Vt)"):
			for _, f := range fields[3:] {
				// if there's a 2 or a second one it's bad
				if f == "2" || (one && f == "1") {
					good = false
					break
				} else if f == "1" {
					one = true
				}
				state = append(state, f)
			}
		case corr && good && strings.Contains(line, "DEGEN   (Vl)"):
		case corr && good && len(fields) > 0:
			if filename == "testfiles/degen.out" {
				fmt.Println(good, holdFreq, line)
			}
			for _, f := range fields {
				if f == "2" || (one && f == "1") {
					good = false
					break
				} else if f == "1" {
					one = true
				}
				state = append(state, f)
			}
		case corr && len(fields) == 0 && good:
			if !one {
				res.ZPT = holdZPT
			} else if pos < 2 {
				continue
			} else {
				for pos-2 >= len(res.Corr) {
					res.Corr = append(res.Corr, 0)
				}
				res.Corr[pos-2] = holdFreq
			}
		case strings.Contains(line, "NON-DEG(Vt)"):
			for _, f := range fields[2:] {
				if f == "2" || (one && f == "1") {
					good = false
					break
				} else if f == "1" {
					one = true
				}
				state = append(state, f)
			}
		case good && strings.Contains(line, "BZA"):
			rot = true
			nobza = false
		case good && nobza && strings.Contains(line, "BZS"):
			rot = true
		case rot && good:
			state = nil
			rot = false
			one = false
			tmp := make([]float64, 0, 3)
			for f := range fields {
				v, err := strconv.ParseFloat(fields[f], 64)
				if err != nil {
					v = math.NaN()
				}
				tmp = append(tmp, v)
			}
			res.Rots = append(res.Rots, tmp)
		case delta.MatchString(line):
			// order is DELTA J, K, JK, delta J, K
			// in MHz
			fields := strings.Fields(line)
			f, _ := strconv.ParseFloat(fields[len(fields)-1],
				64)
			res.Deltas = append(res.Deltas, f)
		case phi.MatchString(line):
			// order is PHI J, K, JK, KJ, phi j, jk, k
			// in Hz
			// may need this in delta too
			line := strings.ReplaceAll(line, "D", "E")
			fields := strings.Fields(line)
			f, _ := strconv.ParseFloat(fields[len(fields)-1],
				64)
			res.Phis = append(res.Phis, f)
		case strings.Contains(line, "INT COORD TYPE") &&
			!geom && res.Requil == nil:
			geom = true
			skip++
		case geom && !strings.Contains(line, "LINEAR"):
			if line == "" {
				geom = false
				continue
			}
			fields := strings.Fields(line)
			e, _ := strconv.ParseFloat(fields[2], 64)
			a, _ := strconv.ParseFloat(fields[4], 64)
			res.Requil = append(res.Requil, e)
			res.Ralpha = append(res.Ralpha, a)
		case icn.MatchString(line):
			// Torsions do not appear in r(equil|alpha) part so
			// neglect here as well
			match := atom.FindAllStringSubmatch(line, -1)
			fmt.Fprintf(&buf, "%s(", gparams[len(match)])
			for l, p := range match {
				fmt.Fprintf(&buf, "%s%s",
					strings.TrimSpace(p[2]),
					strings.TrimSpace(p[1]))
				if l < len(match)-1 {
					fmt.Fprint(&buf, "-")
				}
			}
			fmt.Fprint(&buf, ")")
			res.Rhead = append(res.Rhead, buf.String())
			buf.Reset()
		case strings.Contains(line, "FERMI RESONANCE   "):
			skip += 3
			if strings.Contains(line, "TYPE 1") {
				fermi1 = true
			} else {
				fermi2 = true
			}
		case (fermi1 || fermi2) && line == "":
			// just set them both instead of checking
			fermi1 = false
			fermi2 = false
		case fermi1 && !strings.Contains(line, "NOT FOUND"):
			fields := strings.Fields(line)
			key := fields[3]
			fermiMap[key] = append(fermiMap[key],
				fmt.Sprintf("2v_%s", fields[1]))
		case fermi2 && !strings.Contains(line, "NOT FOUND"):
			fields := strings.Fields(line)
			key := fields[3]
			fermiMap[key] = append(fermiMap[key],
				fmt.Sprintf("v_%s+v_%s", fields[1],
					fields[2]))
		}
		// TODO option for BZA and/or BZS
		// TODO option for D in addition to DELTA
	}
	sorter := make([]string, 0, len(fermiMap))
	for k := range fermiMap {
		sorter = append(sorter, k)
	}
	sort.Strings(sorter)
	for _, k := range sorter {
		for _, r := range fermiMap[k] {
			fmt.Fprintf(&buf, "%s=", r)
		}
		fmt.Fprintf(&buf, "v_%s", k)
		res.Fermi = append(res.Fermi, buf.String())
		buf.Reset()
	}
	return res
}

type Atom struct {
	Sym     string
	X, Y, Z float64
}

func (a Atom) String() string {
	return fmt.Sprintf(
		"%2s%15.10f%15.10f%15.10f\n",
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
	SiIC [][]int
	SyIC [][]int
	Freq []float64
	Vibs string
}

func (id Intder) PrintSiic(siic []int) string {
	var str strings.Builder
	last := siic[len(siic)-1]
	switch last {
	case STRE:
		fmt.Fprintf(&str, "r(%s_%d - %s_%d)",
			id.Geom[siic[0]].Sym, siic[0]+1,
			id.Geom[siic[1]].Sym, siic[1]+1,
		)
	case BEND:
		fmt.Fprintf(&str, "<(%s_%d - %s_%d - %s_%d)",
			id.Geom[siic[0]].Sym, siic[0]+1,
			id.Geom[siic[1]].Sym, siic[1]+1,
			id.Geom[siic[2]].Sym, siic[2]+1,
		)
	case TORS:
		fmt.Fprintf(&str, "t(%s_%d - %s_%d - %s_%d - %s_%d)",
			id.Geom[siic[0]].Sym, siic[0]+1,
			id.Geom[siic[1]].Sym, siic[1]+1,
			id.Geom[siic[2]].Sym, siic[2]+1,
			id.Geom[siic[3]].Sym, siic[3]+1,
		)
	}
	return str.String()
}

func (id Intder) String() string {
	var str strings.Builder
	str.WriteString("Geometry:\n")
	for _, atom := range id.Geom {
		str.WriteString(atom.String())
	}
	str.WriteString("Simple Internals:\n")
	for d, siic := range id.SiIC {
		fmt.Fprintf(&str, "%2d\t%s\n", d+1, id.PrintSiic(siic))
	}
	str.WriteString("Symmetry Internals:\n")
	for d, syic := range id.SyIC {
		fmt.Fprintf(&str, "%2d\t", d+1)
		for i, j := range syic {
			if j < 0 {
				fmt.Fprint(&str, " - ")
				j = -j
			} else if i > 0 {
				fmt.Fprint(&str, " + ")
			}
			str.WriteString(id.PrintSiic(id.SiIC[j]))
		}
		str.WriteString("\n")
	}
	str.WriteString("Vibrational Assignments:\n")
	vibs := strings.Split(strings.TrimSpace(id.Vibs), "\n")
	for i := range id.Freq {
		fmt.Fprintf(&str, "%6.1f\t%s\n", id.Freq[i], vibs[i])
	}
	return str.String()
}

// ptable is a map from the default string masses in intder to the
// corresponding atomic symbols
var ptable = map[string]string{
	"12.000000": "C",
	"1.007825":  "H",
}

const (
	STRE int = iota
	BEND
	TORS
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
		skip   uint
		line   string
		fields []string
	)
	var str strings.Builder
	contrib := regexp.MustCompile(`(-?)([0-9]+) \( *([0-9]{1,3}\.[0-9])\)`)
	sint := regexp.MustCompile(`(L|S)\( ?([0-9]{1,2})\)=?`)
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
			"VIBRATIONAL ASSIGNMENTS"):
			skip += 4
			vibs = true
		case strings.Contains(line,
			"DEFINITION OF INTERNAL COORDINATES"):
			skip += 3
			siic = true
		case geom && len(fields) == 0:
			geom = false
		case vibs && len(fields) == 0:
			vibs = false
			id.Vibs = str.String()
			str.Reset()
		case siic && len(fields) == 0:
			siic = false
			skip += 2
			syic = true
		case syic && len(fields) == 0:
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
			default:
				panic("this type of coordinate not implemented")
			}
			id.SiIC = append(id.SiIC, ids)
		case syic:
			ids := make([]int, 0)
			fields = strings.Fields(sint.ReplaceAllString(line, "${2}"))
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
					d, _ := strconv.Atoi(f)
					ids = append(ids, fac*(d-1))
				}
			}
			id.SyIC = append(id.SyIC, ids)
		}
	}
	return id
}
