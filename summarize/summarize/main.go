package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"strconv"

	"bytes"

	"text/template"

	"github.com/ntBre/chemutils/summarize"
)

var (
	DeltaOrder []string
	PhiOrder   []string
	t          *template.Template
)

func colPrint(format string, cols ...[]float64) string {
	var buf bytes.Buffer
	for i := range cols[0] {
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
		Alignment: "ccc",
		Header:    fmt.Sprintf("%8s%8s%8s", "HARM", "FUND", "CORR"),
		Body:      str.String(),
	}
}

func makeABC(res *summarize.Result) *Table {
	var str strings.Builder
	for a := range res.Rots {
		fmt.Fprintf(&str, "%8s%10.6f\n",
			fmt.Sprintf("A_%d", a), res.Rots[a][2])
		fmt.Fprintf(&str, "%8s%10.6f\n",
			fmt.Sprintf("B_%d", a), res.Rots[a][0])
		fmt.Fprintf(&str, "%8s%10.6f",
			fmt.Sprintf("C_%d", a), res.Rots[a][1])
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

func makeDeltas(res *summarize.Result) *Table {
	var str strings.Builder
	for d := range res.Deltas {
		fmt.Fprintf(&str, "%8s%15.3f%15.3f%15.3f%15.3f%18.3f",
			DeltaOrder[d], res.Deltas[d]/1e3, res.Deltas[d],
			res.Deltas[d]*1e3, res.Deltas[d]*1e6, res.Deltas[d]*1e9)
		if d != len(res.Deltas)-1 {
			fmt.Fprint(&str, "\n")
		}
	}
	return &Table{
		Caption:   "Deltas",
		Alignment: "cccccc",
		Header: fmt.Sprintf("%8s%15s%15s%15s%15s%18s",
			"", "GHz", "MHz", "kHz", "Hz", "mHz"),
		Body: str.String(),
	}
}

func makePhis(res *summarize.Result) *Table {
	var str strings.Builder
	for p := range res.Phis {
		fmt.Fprintf(&str, "%8s%15.3f%15.3f%15.3f%15.3f%18.3f",
			PhiOrder[p], res.Phis[p]/1e3, res.Phis[p],
			res.Phis[p]*1e3, res.Phis[p]*1e6, res.Phis[p]*1e9)
		if p != len(res.Phis)-1 {
			fmt.Fprint(&str, "\n")
		}
	}
	return &Table{
		Caption:   "Phis",
		Alignment: "cccccc",
		Header: fmt.Sprintf("%8s%15s%15s%15s%15s%18s",
			"", "kHz", "Hz", "mHz", "uHz", "nHz"),
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
		Alignment: "ccc",
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
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "summarize: not enough arguments")
		os.Exit(1)
	}
	filename := args[0]
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "summarize: %q is not a file\n",
			filename)
		os.Exit(1)
	}
	nfreqs, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "summarize: %v\n", err)
		os.Exit(1)
	}
	res := summarize.Spectro(filename, nfreqs)
	printAll(os.Stdout, res)
}
