package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	nxo   = flag.Bool("no-xml-output", false, "dummy flag")
	procs = flag.Int("t", 1, "dummy flag")
)

//go:embed test.json
var js string

func main() {
	geoms := make(map[string]float64)
	f := strings.NewReader(js)
	byts, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(byts, &geoms)
	if err != nil {
		panic(err)
	}
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
	base := args[0][:len(args[0])-len(filepath.Ext(args[0]))]
	outfile, err := os.Create(base + ".out")
	defer outfile.Close()
	if err != nil {
		panic(err)
	}
	// TODO include gradients
	scanner := bufio.NewScanner(infile)
	var (
		geom bool
		str  strings.Builder
	)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		switch {
		case strings.Contains(line, "geometry={"):
			geom = true
		case strings.Contains(line, "}") && geom:
			geom = false
		case geom && len(fields) == 4:
			str.WriteString(strings.Join(fields, " ") + "\n")
		}
	}
	val, ok := geoms[str.String()]
	if !ok {
		panic("geometry not found")
	}
	fmt.Fprintf(outfile, "dummy output\nenergy= %20.12f\n", val)
}
