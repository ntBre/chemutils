package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

const (
	width  = 256
	height = 256
)

// plot axes

// identify the center of mass of the molecule

// translate the molecule so its center of mass is at the origin

// determine the moments of inertia to identify the rotation axes for
// each constant

// plot those axes

// identify plane of molecule, make that be plane of screen

// (wobuf "sxiv test.png")

// DrawCircle draws a circle of radius r at c
func DrawCircle(img *image.NRGBA, c image.Point, r int) {
	var ex, ey int
	for y := c.Y - r; y <= c.Y+r; y++ {
		for x := c.X - r; x <= c.X+r; x++ {
			ex = x - c.X
			ey = y - c.Y
			if ex*ex+ey*ey < r*r {
				img.Set(x, y, color.NRGBA{0, 255, 0, 255})
			}
		}
	}
}

// DrawSquare draws a square from (fx, fy) to (tx, ty)
func DrawRect(img *image.NRGBA, from, to image.Point) {
	for y := from.Y; y <= to.Y; y++ {
		for x := from.X; x <= to.X; x++ {
			img.Set(x, y, color.NRGBA{255, 0, 0, 255})
		}
	}
}

// DrawLine draws a line from from to to
func DrawLine(img *image.NRGBA, from, to image.Point) {
	// vertical line
	if from.X == to.X {
		for y := from.Y; y <= to.Y; y++ {
			img.Set(to.X, y, color.NRGBA{0, 0, 0, 255})
		}
		return
	}
	m := (to.Y - from.Y) / (to.X - from.X)
	b := to.Y - m*to.X
	for x := from.X; x <= to.X; x++ {
		img.Set(x, m*x+b, color.NRGBA{0, 0, 0, 255})
	}
}

// PlotAxes draws axes onto img using DrawLine
func PlotAxes(img *image.NRGBA) {
	DrawLine(img, image.Point{width / 2, 0}, image.Point{width / 2, height})
	DrawLine(img, image.Point{0, height / 2}, image.Point{width, height / 2})
	DrawLine(img, image.Point{0, height}, image.Point{width, 0})
}

func main() {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	PlotAxes(img)
	DrawRect(img, image.Point{0, 0}, image.Point{5, 5})
	DrawCircle(img, image.Point{width / 2, height / 2}, 5)
	f, _ := os.Create("test.png")
	png.Encode(f, img)
	f.Close()
}
