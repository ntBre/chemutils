package main

import (
	"fmt"
	"os"

	"strconv"

	"github.com/ntBre/chemutils/summarize"
)

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
	res.PrintAll(os.Stdout)
}
