package main

import (
	"io/ioutil"
	"reflect"
	"testing"

	"bytes"

	"github.com/ntBre/chemutils/summarize"
)

func TestColPrint(t *testing.T) {
	datum0 := []float64{1, 2, 3}
	datum1 := []float64{4, 5, 6}
	got := colPrint("%4.1f", datum0, datum1)
	want := " 1.0 4.0\n 2.0 5.0\n 3.0 6.0\n"
	if got != want {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

func TestPrintResult(t *testing.T) {
	initConst()
	res := summarize.Spectro("testfiles/spectro.out", 6)
	truefile := "testfiles/summary.txt"
	var got bytes.Buffer
	printResult(&got, res)
	want, _ := ioutil.ReadFile(truefile)
	if !reflect.DeepEqual(got.Bytes(), want) {
		t.Errorf("got %v, wanted %v\n", got.String(), string(want))
	}
}
