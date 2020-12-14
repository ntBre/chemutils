package main

import (
	"fmt"
	"os"

	"strconv"

	"bytes"

	"github.com/ntBre/chemutils/summarize"
)

// unicode characters
const (
	upperDelta = "\u0394"
	lowerDelta = "\u03B4"
	upperPhi   = "\u03A6"
	lowerPhi   = "\u03C6"
)

// Exported variables
var (
	DeltaOrder = []string{
		upperDelta + "_J ",
		upperDelta + "_K ",
		upperDelta + "_JK",
		lowerDelta + "_J ",
		lowerDelta + "_K ",
	}
	PhiOrder = []string{
		upperPhi + "_J ",
		upperPhi + "_K ",
		upperPhi + "_JK",
		upperPhi + "_KJ",
		lowerPhi + "_j ",
		lowerPhi + "_jk",
		lowerPhi + "_k ",
	}
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
	zpt, harm, fund, corr,
		rotABC, deltas, phis,
		requil, ralpha, rhead, fermi := summarize.Spectro(filename, nfreqs)
	fmt.Printf("ZPT (cm-1): %.1f\n", zpt)
	// TODO flag for specifying width and probably precision too
	// TODO flag for org, tex format
	width := "8"
	fmt.Println("Freqs (cm-1):")
	fmt.Printf("%"+width+"s"+"%"+width+"s"+"%"+width+"s"+"\n",
		"HARM", "FUND", "CORR")
	// TODO check dimension mismatch before calling this
	fmt.Print(colPrint("%"+width+".1f", harm, fund, corr))
	// TODO flag for units/format here
	// TODO convert these to MHz
	fmt.Println("ABC (cm-1):")
	for a := range rotABC {
		fmt.Printf("A_%d%10.6f\n", a, rotABC[a][2])
		fmt.Printf("B_%d%10.6f\n", a, rotABC[a][0])
		fmt.Printf("C_%d%10.6f\n", a, rotABC[a][1])
	}
	fmt.Println("Deltas (GHz MHz kHz Hz mHz):")
	// TODO flag to disable unicode output?
	for d := range deltas {
		fmt.Printf("%s%15.3f%15.3f%15.3f%15.3f%15.3f\n",
			DeltaOrder[d], deltas[d]/1e3, deltas[d],
			deltas[d]*1e3, deltas[d]*1e6, deltas[d]*1e9)
	}
	fmt.Println("Phis (kHz Hz mHz uHz nHz):")
	for p := range phis {
		fmt.Printf("%s%15.3f%15.3f%15.3f%15.3f%15.3f\n",
			PhiOrder[p], phis[p]/1e3, phis[p],
			phis[p]*1e3, phis[p]*1e6, phis[p]*1e9)
	}
	fmt.Println("Geom (A or Deg):")
	fmt.Printf("%15s%15s%15s\n", "COORD", "R(EQUIL)", "R(ALPHA)")
	for g := range requil {
		fmt.Printf("%15s%15.7f%15.7f\n", rhead[g], requil[g], ralpha[g])
	}
	fmt.Println("Fermi Resonances:")
	for r := range fermi {
		fmt.Println(fermi[r])
	}
}
