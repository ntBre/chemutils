package main

import (
	"flag"
	"fmt"
	"os"

	"strconv"

	"github.com/ntBre/chemutils/summarize"
)

const (
	help = `summarize is a tool for summarizing output from quantum chemistry programs. Currently supported programs are
- spectro
Usage:
summarize <filename> <nfreqs>
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
	// zpt, harm, fund, corr,
	// 	rotABC, deltas, phis := summarize.Spectro(filename, nfreqs)
	fmt.Println(summarize.Spectro(filename, nfreqs))
}
