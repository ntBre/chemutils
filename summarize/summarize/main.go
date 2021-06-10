package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"bytes"

	"text/template"

	"regexp"

	"github.com/ntBre/chemutils/summarize"
)

const (
	toMHz = 29979.2458
)

var (
	DeltaOrder []string
	PhiOrder   []string
	ABC        = []string{"A_%d", "B_%d", "C_%d"}
	t          *template.Template
)

func colPrint(format string, cols ...[]float64) string {
	var buf bytes.Buffer
	for i := range cols[0] {
		fmt.Fprintf(&buf, "%5d", i+1)
		for j := range cols {
			fmt.Fprintf(&buf, format, cols[j][i])
		}
		if i != len(cols[0])-1 {
			fmt.Fprint(&buf, "\n")
		}
	}
	return buf.String()
}

func makeFreqs(res *summarize.Result) *Table {
	var str strings.Builder
	if lh := len(res.Harm); !(lh == len(res.Fund) && lh == len(res.Corr)) {
		panic("dimension mismatch")
	}
	fmt.Fprint(&str, colPrint("%8.1f", res.Harm, res.Fund, res.Corr))
	return &Table{
		Caption:   fmt.Sprintf("Freqs, ZPT=%.1f (cm-1)", res.ZPT),
		Alignment: "crrr",
		Header:    fmt.Sprintf("%5s%8s%8s%8s", "Mode", "HARM", "FUND", "CORR"),
		Body:      str.String(),
	}
}

func makeABC(res *summarize.Result) *Table {
	var str strings.Builder
	for a := range res.Rots {
		fmt.Fprintf(&str, "%8s%10.1f\n",
			fmt.Sprintf(ABC[0], a), res.Rots[a][2]*toMHz)
		fmt.Fprintf(&str, "%8s%10.1f\n",
			fmt.Sprintf(ABC[1], a), res.Rots[a][0]*toMHz)
		fmt.Fprintf(&str, "%8s%10.1f",
			fmt.Sprintf(ABC[2], a), res.Rots[a][1]*toMHz)
		if a != len(res.Rots)-1 {
			fmt.Fprint(&str, "\n")
		}
	}
	return &Table{
		Caption:   "ABC (MHz)",
		Alignment: "cc",
		Header:    fmt.Sprintf("%8s%10s", "Constant", "Value"),
		Body:      str.String(),
	}
}

type Unit struct {
	Name  string
	Scale float64
}

// tryUnit takes a value and a scale to multiply by and returns the
// string version and its length
func tryUnit(val, scale float64) (string, int) {
	ret := fmt.Sprintf("%.3f", val*scale)
	l := len(ret)
	// don't count the sign toward the length
	if val < 0 {
		return ret, l - 1
	}
	return ret, l
}

func makeDeltas(res *summarize.Result) *Table {
	units := []Unit{
		{"GHz", 1e-3},
		{"MHz", 1.0},
		{"kHz", 1e3},
		{"Hz", 1e6},
		{"mHz", 1e9},
	}
	var str strings.Builder
	for d := range res.Deltas {
		for _, u := range units {
			s, l := tryUnit(res.Deltas[d], u.Scale)
			if l > 4 && s[0] != '0' && l <= 7 {
				fmt.Fprintf(&str, "%8s%10s%10s",
					DeltaOrder[d], u.Name, s)
				break
			}
		}
		if d != len(res.Deltas)-1 {
			fmt.Fprint(&str, "\n")
		}
	}
	return &Table{
		Caption:   "Deltas",
		Alignment: "llr",
		Header: fmt.Sprintf("%8s%10s%10s",
			"Constant", "Units", "Value"),
		Body: str.String(),
	}
}

func makePhis(res *summarize.Result) *Table {
	units := []Unit{
		{"kHz", 1e-3},
		{"Hz", 1.0},
		{"mHz", 1e3},
		{"uHz", 1e6},
		{"nHz", 1e9},
	}
	var str strings.Builder
	for p := range res.Phis {
		for _, u := range units {
			s, l := tryUnit(res.Phis[p], u.Scale)
			if l > 4 && s[0] != '0' && s[:2] != "-0" && l <= 7 {
				fmt.Fprintf(&str, "%8s%10s%10s",
					PhiOrder[p], u.Name, s)
				break
			}
		}
		if p != len(res.Phis)-1 {
			fmt.Fprint(&str, "\n")
		}
	}
	return &Table{
		Caption:   "Phis",
		Alignment: "llr",
		Header: fmt.Sprintf("%8s%10s%10s",
			"Constant", "Units", "Value"),
		Body: str.String(),
	}
}

func makeGeom(res *summarize.Result) *Table {
	var str strings.Builder
	for g := range res.Requil {
		fmt.Fprintf(&str, "%15s%15.7f%15.7f", res.Rhead[g],
			res.Requil[g], res.Ralpha[g])
		if g != len(res.Requil)-1 {
			fmt.Fprint(&str, "\n")
		}
	}
	return &Table{
		Caption:   "Geom (A or Deg)",
		Alignment: "rrr",
		Header: fmt.Sprintf("%15s%15s%15s",
			"COORD", "R(EQUIL)", "R(ALPHA)"),
		Body: str.String(),
	}
}

func makeFermi(res *summarize.Result) *Table {
	var str strings.Builder
	for r := range res.Fermi {
		fmt.Fprintf(&str, "%s", res.Fermi[r])
		if r != len(res.Fermi)-1 {
			fmt.Fprint(&str, "\n")
		}
	}
	return &Table{
		Caption:   "Fermi Resonances",
		Alignment: "c",
		Body:      str.String(),
	}
}

func printAll(out io.Writer, res *summarize.Result) {
	switch {
	case *freq && *rot:
		t.Execute(out, makeFreqs(res))
		t.Execute(out, makeABC(res))
	case *freq:
		t.Execute(out, makeFreqs(res))
	case *rot:
		t.Execute(out, makeABC(res))
	default:
		t.Execute(out, makeFreqs(res))
		t.Execute(out, makeABC(res))
		t.Execute(out, makeDeltas(res))
		t.Execute(out, makePhis(res))
		t.Execute(out, makeGeom(res))
		t.Execute(out, makeFermi(res))
	}
}

func printSiic(id *summarize.Intder, siic []int) string {
	var str strings.Builder
	last := siic[len(siic)-1]
	switch last {
	case summarize.STRE:
		fmt.Fprintf(&str, "r(%s_%d - %s_%d)",
			id.Geom[siic[0]].Sym, siic[0]+1,
			id.Geom[siic[1]].Sym, siic[1]+1,
		)
	case summarize.BEND:
		fmt.Fprintf(&str, "<(%s_%d - %s_%d - %s_%d)",
			id.Geom[siic[0]].Sym, siic[0]+1,
			id.Geom[siic[1]].Sym, siic[1]+1,
			id.Geom[siic[2]].Sym, siic[2]+1,
		)
	case summarize.TORS:
		fmt.Fprintf(&str, "t(%s_%d - %s_%d - %s_%d - %s_%d)",
			id.Geom[siic[0]].Sym, siic[0]+1,
			id.Geom[siic[1]].Sym, siic[1]+1,
			id.Geom[siic[2]].Sym, siic[2]+1,
			id.Geom[siic[3]].Sym, siic[3]+1,
		)
	case summarize.LIN1:
		geom := append(id.Geom, id.Dumm...)
		fmt.Fprintf(&str, "LIN(%s_%d - %s_%d - %s_%d - %s_%d)",
			geom[siic[0]].Sym, siic[0]+1,
			geom[siic[1]].Sym, siic[1]+1,
			geom[siic[2]].Sym, siic[2]+1,
			geom[siic[3]].Sym, siic[3]+1,
		)
	}
	return str.String()
}

func TrimNewline(str string) string {
	return strings.TrimRight(str, "\n")
}

func makeIntderGeom(id *summarize.Intder) *Table {
	tab := new(Table)
	var str strings.Builder
	tab.Caption = "Geometry"
	for _, atom := range append(id.Geom, id.Dumm...) {
		fmt.Fprintf(&str, "%s\n", atom.String())
	}
	tab.Body = TrimNewline(str.String())
	return tab
}

func makeSiIC(id *summarize.Intder) *Table {
	tab := new(Table)
	var str strings.Builder
	tab.Caption = "Simple Internals"
	for d, siic := range id.SiIC {
		fmt.Fprintf(&str, "%2d\t%s\n", d+1, printSiic(id, siic))
	}
	tab.Body = TrimNewline(str.String())
	return tab
}

// eqnify converts SICs to the format needed for a LaTeX eqnarray
func Eqnify(str string, end bool) string {
	split := strings.Split(str, "\t")
	atom := regexp.MustCompile(`([A-Z][a-z]?)_`)
	cord := regexp.MustCompile(`([^A-Z]r|<|t|LIN)`)
	cords := cord.FindAllString(str, -1)
	split[1] = strings.Replace(split[1], "<", `\angle`, -1)
	split[1] = strings.Replace(split[1], "t", `\tau`, -1)
	split[1] = atom.ReplaceAllString(split[1], `\text{$1}_`)
	term := `\\`
	// don't want \\ on last line of eqnarray
	if end {
		term = ""
	}
	switch len(cords) {
	case 1:
		return fmt.Sprintf(`S_{%s} &= &%s%s`,
			split[0], split[1], term)
	case 2:
		return fmt.Sprintf(`S_{%s} &= &\frac{1}{\sqrt{2}}[%s]%s`,
			split[0], split[1], term)
	default:
		fmt.Printf("%q -> %v\n", str, cords)
		panic("unrecognized number of SICs")
	}
}

func makeSyIC(id *summarize.Intder) *Table {
	tab := new(Table)
	var str, line strings.Builder
	tab.Caption = "Symmetry Internals"
	for d, syic := range id.SyIC {
		fmt.Fprintf(&line, "%2d\t", d+1)
		for i, j := range syic {
			if j < 0 {
				fmt.Fprint(&line, " - ")
				j = -j
			} else if i > 0 {
				fmt.Fprint(&line, " + ")
			}
			fmt.Fprint(&line, printSiic(id, id.SiIC[j]))
		}
		if *tex {
			fmt.Fprint(&str, Eqnify(line.String(), d == len(id.SyIC)-1))
		} else {
			fmt.Fprint(&str, line.String())
		}
		line.Reset()
		fmt.Fprint(&str, "\n")
	}
	tab.Body = TrimNewline(str.String())
	return tab
}

func makeVibs(id *summarize.Intder) *Table {
	tab := new(Table)
	var str strings.Builder
	tab.Caption = "Vibrational Assignments"
	tab.Header = fmt.Sprintf("%6s\tContribs", "Freq")
	vibs := strings.Split(strings.TrimSpace(id.Vibs), "\n")
	for i := range id.Freq {
		fmt.Fprintf(&str, "%6.1f\t%s\n", id.Freq[i], vibs[i])
	}
	tab.Body = TrimNewline(str.String())
	return tab
}

func printIntder(out io.Writer, id *summarize.Intder) {
	t.Execute(out, makeIntderGeom(id))
	t.Execute(out, makeSiIC(id))
	t.Execute(out, makeSyIC(id))
	t.Execute(out, makeVibs(id))
}

func main() {
	filename := parseFlags()
	if *tex && *spectro && !*nohead {
		fmt.Print("\\documentclass{article}\n\\begin{document}\n\n")
		defer func() {
			fmt.Print("\\end{document}\n")
		}()
	}
	if *spectro {
		res := summarize.Spectro(filename)
		printAll(os.Stdout, res)
	} else if *intder {
		id := summarize.ReadIntder(filename)
		printIntder(os.Stdout, id)
	}
}
