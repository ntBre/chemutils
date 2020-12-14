package main

import "testing"

func TestColPrint(t *testing.T) {
	datum0 := []float64{1, 2, 3}
	datum1 := []float64{4, 5, 6}
	got := colPrint("%4.1f", datum0, datum1)
	want := " 1.0 4.0\n 2.0 5.0\n 3.0 6.0\n"
	if got != want {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

func TestOutput(t *testing.T) {
	if got != want {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}
