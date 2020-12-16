package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"bytes"

	"os"

	"os/exec"

	"github.com/ntBre/chemutils/summarize"
)

func TestColPrint(t *testing.T) {
	datum0 := []float64{1, 2, 3}
	datum1 := []float64{4, 5, 6}
	got := colPrint("%4.1f", datum0, datum1)
	want := "    1 1.0 4.0\n    2 2.0 5.0\n    3 3.0 6.0"
	if got != want {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

func TestPrintResult(t *testing.T) {
	initConst()
	res := summarize.Spectro("testfiles/spectro.out")
	truefile := "testfiles/summary.txt"
	var got bytes.Buffer
	printAll(&got, res)
	want, _ := ioutil.ReadFile(truefile)
	if !reflect.DeepEqual(got.Bytes(), want) {
		t.Errorf("got:\n%v, wanted:\n%v\n", got.String(), string(want))
		out, _ := os.Create("testfiles/problem.out")
		defer out.Close()
		got.WriteTo(out)
	}
	// (diff "testfiles/problem.out" "testfiles/summary.txt")
}

func TestTex(t *testing.T) {
	// inverting the meaning of short
	if !testing.Short() {
		t.Skip("skipping TestTex")
	}
	dir := "testfiles/tex"
	os.MkdirAll(dir, 0755)
	hold := *tex
	*tex = true
	initConst()
	texname := filepath.Join(dir, "test.tex")
	texfile, _ := os.Create(texname)
	defer func() {
		*tex = hold
		texfile.Close()
	}()
	res := summarize.Spectro("testfiles/spectro.out")
	fmt.Fprint(texfile, "\\documentclass{article}\n\\begin{document}\n\n")
	printAll(texfile, res)
	fmt.Fprint(texfile, "\\end{document}\n")
	pdfcmd := exec.Command("pdflatex", "test.tex")
	pdfcmd.Dir = dir
	pdfcmd.Run()
	zathura := exec.Command("zathura", "test.pdf")
	zathura.Dir = dir
	zathura.Run()
}
