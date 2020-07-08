package main

import (
	"flag"
	"fmt"
	"os"

	"strconv"

	"bytes"

	"github.com/ntBre/chemutils/summarize"
)

const (
	help = `summarize is a tool for summarizing output from quantum chemistry
programs. Currently supported programs are
- spectro
Usage:
$ summarize <filename> <nfreqs>
Flags:
`
)

func parseFlags() []string {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s", help)
		flag.PrintDefaults()
	}
	flag.Parse()
	return flag.Args()
}

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
		rotABC, deltas, phis := summarize.Spectro(filename, nfreqs)
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
	fmt.Println("ABC (cm-1):")
	for a := range rotABC {
		fmt.Printf("A_%d%10.6f\n", a, rotABC[a][2])
		fmt.Printf("B_%d%10.6f\n", a, rotABC[a][0])
		fmt.Printf("C_%d%10.6f\n", a, rotABC[a][1])
	}
	fmt.Println("Deltas (MHz):")
	// TODO flag to disable unicode output?
	for d := range deltas {
		fmt.Printf("%s%15.10f\n", summarize.DeltaOrder[d],
			deltas[d])
	}
	fmt.Println("Phis (Hz):")
	for p := range phis {
		fmt.Printf("%s%20.10e\n", summarize.PhiOrder[p],
			phis[p])
	}
}
