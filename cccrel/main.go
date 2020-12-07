package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const (
	VTZ = iota
	VQZ
	V5Z
	MT
	MTC
	DK
	DKR
	NBASIS
)

var (
	BASIS = []string{"VTZ", "VQZ", "V5Z", "MT", "MTC", "DK", "DKR"}
)

// Flags
var (
	eline = flag.String("e", "energy= ", "pattern to match for energy")
	base  = flag.String("b", "pts", "base directory")
	dirs  = flag.String("d", "avtz avqz av5z mt mtc dk dkr",
		"component directories to read, preserve the default order if changed")
	inp  = flag.String("i", "inp", "input directory")
	suff = flag.String("s", "out", "output file suffix")
)

type Job struct {
	filename string
	energy   float64
}

func ReadDir(dir, suffix string) []Job {
	jobs := make([]Job, 0)
	files, err := filepath.Glob(dir + "/*." + suffix)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		e, err := ReadOut(f)
		if err != nil {
			panic(err)
		}
		jobs = append(jobs,
			Job{
				filename: filepath.Base(f),
				energy:   e,
			})
	}
	return jobs
}

func ReadOut(filename string) (float64, error) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, *eline) {
			fields := strings.Fields(line)
			f, _ := strconv.ParseFloat(fields[len(fields)-1], 64)
			return f, nil
		}
	}
	return 0, errors.New("energy not found")
}

func CcCR(tz, qz, fz, mt, mtc, dk, dkr float64) float64 {
	three54 := 1.0 / (3.5 * 3.5 * 3.5 * 3.5)
	three56 := 1.0 / (3.5 * 3.5 * 3.5 * 3.5 * 3.5 * 3.5)
	five56 := 1.0 / (5.5 * 5.5 * 5.5 * 5.5 * 5.5 * 5.5)
	magic := 0.7477488413
	a := 1.0/(4.5*4.5*4.5*4.5) - three54
	b := three54 - 1.0/(5.5*5.5*5.5*5.5)
	return tz - (qz-tz)*three54/a +
		(magic*three54-three56)*((fz-tz+(qz-tz)*b/a)/(magic*b-three56+five56)) +
		mtc - mt + dkr - dk
}

func main() {
	flag.Parse()
	toread := strings.Fields(*dirs)
	e := make([][]Job, NBASIS)
	var wg sync.WaitGroup
	fmt.Println("READING OUTPUT FILES")
	fmt.Println("--------------------")
	for b, r := range toread {
		read := filepath.Join(*base, r, *inp)
		fmt.Printf("Reading %s/*.%s as %s input\n", read, *suff, BASIS[b])
		wg.Add(1)
		go func(b int, read string) {
			defer wg.Done()
			e[b] = ReadDir(read, *suff)
		}(b, read)
	}
	wg.Wait()
	fmt.Println("\nRAW ENERGIES")
	fmt.Println("------------")
	fmt.Printf("%15s", "Basis")
	for b := range toread {
		fmt.Printf("%20s", BASIS[b])
	}
	fmt.Print("\n")
	fmt.Printf("%15s", "Directory")
	for _, r := range toread {
		fmt.Printf("%20s", r)
	}
	fmt.Print("\n")
	for i := range e[VTZ] {
		compare := e[VTZ][i].filename
		fmt.Printf("%15s", compare)
		for b := range BASIS {
			if e[b][i].filename != compare {
				panic("filename order mismatch")
			}
			fmt.Printf("%20.12f", e[b][i].energy)
		}
		fmt.Print("\n")
	}

	fmt.Println("\nRELATIVE ENERGIES")
	fmt.Println("-----------------")
	energies := make([]float64, len(e[VTZ]))
	min := 0.0
	for i := range e[VTZ] {
		energy := CcCR(e[VTZ][i].energy, e[VQZ][i].energy,
			e[V5Z][i].energy, e[MT][i].energy, e[MTC][i].energy,
			e[DK][i].energy, e[DKR][i].energy)
		if energy < min {
			min = energy
		}
		energies[i] = energy
	}
	for _, e := range energies {
		fmt.Printf("%20.12f\n", e-min)
	}
}
