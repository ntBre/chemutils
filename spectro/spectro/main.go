package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/ntBre/chemutils/spectro"
)

// Flags
var (
	cmd = flag.String("cmd",
		"/ddn/home6/r2533/programs/spec3jm.ifort-O0.static.x",
		"command to use for running SPECTRO")
)

func ParseFlags() []string {
	flag.Parse()
	spectro.Command = *cmd
	return flag.Args()
}

func main() {
	args := ParseFlags()
	if len(args) < 1 {
		log.Fatal("spectro: not enough input arguments\n")
	}
	dir := filepath.Dir(args[0])
	spec, err := spectro.Load(args[0])
	if err != nil {
		log.Fatalf("spectro: %v\n", err)
	}
	err = spec.DoSpectro(dir)
	if err != nil {
		log.Fatalf("spectro: %v\n", err)
	}
}
