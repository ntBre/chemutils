package main

import (
	"flag"
	"fmt"
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

var (
	tex   = flag.Bool("tex", false, "output summary in TeX table format")
	plain = flag.Bool("plain", false, "disable Unicode characters in txt output")
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
