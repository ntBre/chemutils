package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	help = `summarize is a tool for summarizing output from quantum chemistry
programs. Currently supported programs are
- spectro
- intder
Usage:
$ summarize <filename>
Flags:
`
)

var (
	freq  = flag.Bool("v", false, "only print vibrational frequencies and exit")
	lfreq = flag.Bool("vl", false,
		"only print the resonance-corrected vibrational frequencies")
	rot     = flag.Bool("r", false, "only print principal rotational constants and exit")
	spectro = flag.Bool("s", false, "parse a spectro output file")
	intder  = flag.Bool("i", false, "parse an intder output file")
	tex     = flag.Bool("tex", false, "output summary in TeX table format")
	plain   = flag.Bool("plain", false, "disable Unicode characters in txt output")
	nohead  = flag.Bool("nohead", false, "disable printing of header info for TeX output")
	cm      = flag.Bool("cm", false, "print principal rotational constants in cm-1")
)

func parseFlags() []string {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s", help)
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "summarize: not enough arguments")
		os.Exit(1)
	}
	for _, filename := range args {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "summarize: %q is not a file\n",
				filename)
			os.Exit(1)
		}
		if strings.Contains(filename, "spectro") {
			*spectro = true
		} else if strings.Contains(filename, "intder") {
			*intder = true
		}
	}
	initConst()
	return args
}
