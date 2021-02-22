package main

import (
	"flag"
	"fmt"
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
	spectro = flag.Bool("s", false, "parse a spectro output file")
	intder  = flag.Bool("i", false, "parse an intder output file")
	tex     = flag.Bool("tex", false, "output summary in TeX table format")
	plain   = flag.Bool("plain", false, "disable Unicode characters in txt output")
	nohead  = flag.Bool("nohead", false, "disable printing of header info for TeX output")
)

func parseFlags() []string {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s", help)
		flag.PrintDefaults()
	}
	flag.Parse()
	initConst()
	return flag.Args()
}
