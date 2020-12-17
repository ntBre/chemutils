package main

import (
	"flag"
	"log"

	"github.com/ntBre/chemutils/spectro"
)

// Flags
var (
	cmd = flag.String("cmd",
		"~r2533/programs/spec3jm.ifort-O0.static.x",
		"command to use for running SPECTRO")
)

func ParseFlags() []string {
	flag.Parse()
	spectro.SpectroCommand = *cmd
	return flag.Args()
}

func main() {
	args := ParseFlags()
	if len(args) < 1 {
		log.Fatal("spectro: not enough input arguments\n")
	}
	spec, err := spectro.Load(args[0])
	if err != nil {
		log.Fatalf("spectro: %v\n", err)
	}
	err = spec.DoSpectro(".")
	if err != nil {
		log.Fatalf("spectro: %v\n", err)
	}
}
