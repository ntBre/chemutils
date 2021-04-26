package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var (
	nxo   = flag.Bool("no-xml-output", false, "dummy flag")
	procs = flag.Int("t", 1, "dummy flag")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		panic("not enough arguments in call to molpro")
	}
	infile, err := os.Open(args[0])
	defer infile.Close()
	if err != nil {
		panic(err)
	}
	// TODO grab geometry and maybe whether or not optg is there,
	// but first step is geometry to match to energy
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
}
