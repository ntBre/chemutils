package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestReadDir(t *testing.T) {
	got := ReadDir("testfiles/hcn/pts/av5z/inp", "out")
	want := []Job{
		{"hcn.0001.out", -93.309090707524},
		{"hcn.0002.out", -93.309530672208},
		{"hcn.0003.out", -93.309553104193},
		{"hcn.0004.out", -93.309539599098},
		{"hcn.0005.out", -93.309811147136},
		{"hcn.0006.out", -93.309871778340},
		{"hcn.0007.out", -93.309895431532},
		{"hcn.0008.out", -93.309889588542},
		{"hcn.0009.out", -93.309883140989},
		{"hcn.0010.out", -93.309835914818},
		{"hcn.0011.out", -93.309934902390},
		{"hcn.0012.out", -93.310034807104},
		{"hcn.0013.out", -93.310096665720},
		{"hcn.0014.out", -93.310090879562},
		{"hcn.0015.out", -93.310121539622},
		{"hcn.0016.out", -93.310115769192},
		{"hcn.0017.out", -93.310110463313},
		{"hcn.0018.out", -93.310104708619},
		{"hcn.0019.out", -93.310064445131},
		{"hcn.0020.out", -93.309984467953},
		{"hcn.0021.out", -93.309904577714},
		{"hcn.0022.out", -93.310044858011},
		{"hcn.0023.out", -93.310145996206},
		{"hcn.0024.out", -93.310140267015},
		{"hcn.0025.out", -93.310209081567},
		{"hcn.0026.out", -93.310203367943},
		{"hcn.0027.out", -93.310235175708},
		{"hcn.0028.out", -93.310229477655},
		{"hcn.0029.out", -93.310212381188},
		{"hcn.0030.out", -93.310225313350},
		{"hcn.0031.out", -93.310219630874},
		{"hcn.0032.out", -93.310180503051},
		{"hcn.0033.out", -93.310174836154},
		{"hcn.0034.out", -93.310101727902},
		{"hcn.0035.out", -93.309989946220},
		{"hcn.0036.out", -93.310045990750},
		{"hcn.0037.out", -93.310148361600},
		{"hcn.0038.out", -93.310212673068},
		{"hcn.0039.out", -93.310207031789},
		{"hcn.0040.out", -93.310239986997},
		{"hcn.0041.out", -93.310234361135},
		{"hcn.0042.out", -93.310231338328},
		{"hcn.0043.out", -93.310225727886},
		{"hcn.0044.out", -93.310187735828},
		{"hcn.0045.out", -93.310110162795},
		{"hcn.0046.out", -93.310045453105},
		{"hcn.0047.out", -93.310110990072},
		{"hcn.0048.out", -93.310139523362},
		{"hcn.0049.out", -93.310133969505},
		{"hcn.0050.out", -93.310132088134},
		{"hcn.0051.out", -93.310089693363},
		{"hcn.0052.out", -93.309907487712},
		{"hcn.0053.out", -93.309937239960},
		{"hcn.0054.out", -93.309931017942},
		{"hcn.0055.out", -93.309636499780},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

func TestReadOut(t *testing.T) {
	got, _ := ReadOut("testfiles/hcn/pts/av5z/inp/hcn.0001.out")
	want := -93.309090707524
	if got != want {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

func TestCcCR(t *testing.T) {
	var (
		tz  = -114.342836279341
		qz  = -114.372358595640
		fz  = -114.381621067350
		mt  = -114.230064143533
		mtc = -114.341229956925
		dk  = -114.328478778591
		dkr = -114.400489098870
	)
	got := CcCR(tz, qz, fz, mt, mtc, dk, dkr)
	want := -114.572402781901
	eps := 1e-10
	if math.Abs(got-want) >= eps {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

func loadEDat(filename string) []float64 {
	fmt.Println(filename)
	f, _ := os.Open(filename)
	scanner := bufio.NewScanner(f)
	energies := make([]float64, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			fl, _ := strconv.ParseFloat(strings.TrimSpace(line), 64)
			energies = append(energies, fl)
		}
	}
	return energies
}

func jobEnergies(jobs []Job) []float64 {
	energies := make([]float64, 0, len(jobs))
	for _, j := range jobs {
		energies = append(energies, j.energy)
	}
	return energies
}

// test LoadDirs against the contents of the energy.dat files in each
// directory
func TestLoadDirs(t *testing.T) {
	toread := []string{
		"testfiles/hno/pts/avtz/",
		"testfiles/hno/pts/avqz/",
		"testfiles/hno/pts/av5z/",
		"testfiles/hno/pts/mt/",
		"testfiles/hno/pts/mtc/",
		"testfiles/hno/pts/dk/",
		"testfiles/hno/pts/dkr/",
	}
	*base = ""
	got := LoadDirs(toread)
	for i, line := range *got {
		efile := filepath.Join(toread[i], "inp", "energy.dat")
		want := loadEDat(efile)
		nrg := jobEnergies(line)
		if !reflect.DeepEqual(nrg, want) {
			t.Errorf("got %v, wanted %v\n", nrg, want)
		}
	}
}
