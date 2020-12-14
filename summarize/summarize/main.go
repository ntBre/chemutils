package main

import (
	"fmt"
	"os"

	"strconv"

	"bytes"

	"github.com/ntBre/chemutils/summarize"
)

var (
	DeltaOrder []string
	PhiOrder []string
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
	res := summarize.Spectro(filename, nfreqs)
	fmt.Printf("ZPT (cm-1): %.1f\n", res.ZPT)
	// TODO flag for specifying width and probably precision too
	// TODO flag for org, tex format
	width := "8"
	fmt.Println("Freqs (cm-1):")
	fmt.Printf("%"+width+"s"+"%"+width+"s"+"%"+width+"s"+"\n",
		"HARM", "FUND", "CORR")
	// TODO check dimension mismatch before calling this
	fmt.Print(colPrint("%"+width+".1f", res.Harm, res.Fund, res.Corr))
	// TODO flag for units/format here
	// TODO convert these to MHz
	fmt.Println("ABC (cm-1):")
	for a := range res.Rots {
		fmt.Printf("A_%d%10.6f\n", a, res.Rots[a][2])
		fmt.Printf("B_%d%10.6f\n", a, res.Rots[a][0])
		fmt.Printf("C_%d%10.6f\n", a, res.Rots[a][1])
	}
	fmt.Println("Deltas (GHz MHz kHz Hz mHz):")
	// TODO flag to disable unicode output?
	for d := range res.Deltas {
		fmt.Printf("%s%15.3f%15.3f%15.3f%15.3f%15.3f\n",
			DeltaOrder[d], res.Deltas[d]/1e3, res.Deltas[d],
			res.Deltas[d]*1e3, res.Deltas[d]*1e6, res.Deltas[d]*1e9)
	}
	fmt.Println("Phis (kHz Hz mHz uHz nHz):")
	for p := range res.Phis {
		fmt.Printf("%s%15.3f%15.3f%15.3f%15.3f%15.3f\n",
			PhiOrder[p], res.Phis[p]/1e3, res.Phis[p],
			res.Phis[p]*1e3, res.Phis[p]*1e6, res.Phis[p]*1e9)
	}
	fmt.Println("Geom (A or Deg):")
	fmt.Printf("%15s%15s%15s\n", "COORD", "R(EQUIL)", "R(ALPHA)")
	for g := range res.Requil {
		fmt.Printf("%15s%15.7f%15.7f\n", res.Rhead[g], res.Requil[g], res.Ralpha[g])
	}
	fmt.Println("Fermi Resonances:")
	for r := range res.Fermi {
		fmt.Println(res.Fermi[r])
	}
}
