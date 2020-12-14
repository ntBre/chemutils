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
		fmt.Fprint(&buf, "\n")
	}
	return buf.String()
}

func printDeltas(out io.Writer, res *summarize.Result) {
	fmt.Fprintln(out, "Deltas (GHz MHz kHz Hz mHz):")
	for d := range res.Deltas {
		fmt.Fprintf(out, "%s%15.3f%15.3f%15.3f%15.3f%15.3f\n",
			DeltaOrder[d], res.Deltas[d]/1e3, res.Deltas[d],
			res.Deltas[d]*1e3, res.Deltas[d]*1e6, res.Deltas[d]*1e9)
	}
}

func printABC(out io.Writer, res *summarize.Result) {
	fmt.Fprintln(out, "ABC (cm-1):")
	for a := range res.Rots {
		fmt.Fprintf(out, "A_%d%10.6f\n", a, res.Rots[a][2])
		fmt.Fprintf(out, "B_%d%10.6f\n", a, res.Rots[a][0])
		fmt.Fprintf(out, "C_%d%10.6f\n", a, res.Rots[a][1])
	}
}

func printFreqs(out io.Writer, res *summarize.Result) {
	fmt.Fprintf(out, "ZPT (cm-1): %.1f\n", res.ZPT)
	width := "8"
	fmt.Fprintln(out, "Freqs (cm-1):")
	fmt.Fprintf(out, "%"+width+"s"+"%"+width+"s"+"%"+width+"s"+"\n",
		"HARM", "FUND", "CORR")
	if lh := len(res.Harm); !(lh == len(res.Fund) && lh == len(res.Corr)) {
		panic("dimension mismatch")
	}
	fmt.Fprint(out, colPrint("%"+width+".1f", res.Harm, res.Fund, res.Corr))
}

func printPhis(out io.Writer, res *summarize.Result) {
	fmt.Fprintln(out, "Phis (kHz Hz mHz uHz nHz):")
	for p := range res.Phis {
		fmt.Fprintf(out, "%s%15.3f%15.3f%15.3f%15.3f%15.3f\n",
			PhiOrder[p], res.Phis[p]/1e3, res.Phis[p],
			res.Phis[p]*1e3, res.Phis[p]*1e6, res.Phis[p]*1e9)
	}
}

func printGeom(out io.Writer, res *summarize.Result) {
	fmt.Fprintln(out, "Geom (A or Deg):")
	fmt.Fprintf(out, "%15s%15s%15s\n", "COORD", "R(EQUIL)", "R(ALPHA)")
	for g := range res.Requil {
		fmt.Fprintf(out, "%15s%15.7f%15.7f\n", res.Rhead[g], res.Requil[g], res.Ralpha[g])
	}
}

func makeFermi(res *summarize.Result) *Table {
	var str strings.Builder
	for r := range res.Fermi {
		fmt.Fprintln(&str, res.Fermi[r])
	}
	return &Table{
		Caption:   "Fermi Resonances",
		Alignment: "c",
		Body:      str.String(),
	}
}

func printAll(out io.Writer, res *summarize.Result) {
	printFreqs(out, res)
	printABC(out, res)
	printDeltas(out, res)
	printPhis(out, res)
	printGeom(out, res)
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
