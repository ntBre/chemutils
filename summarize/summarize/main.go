package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"bytes"

	"text/template"

	"github.com/ntBre/chemutils/summarize"
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
		fmt.Fprintf(&str, "%8s%10.6f\n",
			fmt.Sprintf(ABC[0], a), res.Rots[a][2])
		fmt.Fprintf(&str, "%8s%10.6f\n",
			fmt.Sprintf(ABC[1], a), res.Rots[a][0])
		fmt.Fprintf(&str, "%8s%10.6f",
			fmt.Sprintf(ABC[2], a), res.Rots[a][1])
		if a != len(res.Rots)-1 {
			fmt.Fprint(&str, "\n")
		}
	}
	return &Table{
		Caption:   "ABC (cm-1)",
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
	t.Execute(out, makeFreqs(res))
	t.Execute(out, makeABC(res))
	t.Execute(out, makeDeltas(res))
	t.Execute(out, makePhis(res))
	t.Execute(out, makeGeom(res))
	t.Execute(out, makeFermi(res))
}

func main() {
	args := parseFlags()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "summarize: not enough arguments")
		os.Exit(1)
	}
	filename := args[0]
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "summarize: %q is not a file\n",
			filename)
		os.Exit(1)
	}
	res := summarize.Spectro(filename)
	if *tex && !*nohead {
		fmt.Print("\\documentclass{article}\n\\begin{document}\n\n")
	}
	printAll(os.Stdout, res)
	if *tex && !*nohead {
		fmt.Print("\\end{document}\n")
	}
}
