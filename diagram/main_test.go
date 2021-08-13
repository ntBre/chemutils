package main

import (
	"image"
	"reflect"
	"testing"
)

func TestParseCaptions(t *testing.T) {
	got := ParseCaptions("tests/c2h4.cap")
	want := []Caption{
		{"H<sub>1</sub>", 84, image.Point{500, 800}},
		{"H<sub>2</sub>", 84, image.Point{500, 2400}},
		{"H<sub>3</sub>", 84, image.Point{2700, 800}},
		{"H<sub>4</sub>", 84, image.Point{2700, 2400}},
		{"C<sub>1</sub>", 84, image.Point{1000, 1600}},
		{"C<sub>2</sub>", 84, image.Point{2200, 1600}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

func BenchmarkDrawGrid(b *testing.B) {
	img := loadPic("tests/c2h4.png")
	for i := 0; i < b.N; i++ {
		DrawGrid(img, 16, 16)
	}
}
