// (wobuf "sxiv test.png")
package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Axis int

const (
	X Axis = iota
	Y
	Z
)

const (
	width  = 256
	height = 256
)

var (
	RED   = color.NRGBA{255, 0, 0, 255}
	BLACK = color.NRGBA{0, 0, 0, 255}
)

var ptable = map[string]color.NRGBA{
	"H": {0, 0, 0, 128},
	"C": BLACK,
	"O": RED,
}

type Atom struct {
	Symbol string
	X      float64
	Y      float64
	Z      float64
}

// Coords returns the X, Y, and Z coordinates of a as a slice of float
func (a Atom) Coords() []float64 {
	return []float64{a.X, a.Y, a.Z}
}

// Swap swaps axes i and j
func (a Atom) Swap(i, j Axis) Atom {
	coords := []float64{a.X, a.Y, a.Z}
	coords[i], coords[j] = coords[j], coords[i]
	a.X, a.Y, a.Z = coords[0], coords[1], coords[2]
	return a
}

// Invert negates the ith coordinate of a
func (a Atom) Invert(i Axis) Atom {
	c := a.Coords()
	c[i] *= -1.0
	a.X, a.Y, a.Z = c[0], c[1], c[2]
	return a
}

func (a Atom) String() string {
	return fmt.Sprintf("%s %8.4f %8.4f %8.4f",
		a.Symbol, a.X, a.Y, a.Z)
}

type Output struct {
	Geom []Atom
	max  float64
}

// Cart2D converts the point (x, y, z) to an image.Point
func Cart2D(x, y, z float64) image.Point {
	cw, ch := float64(width/2), float64(height/2)
	wx := math.Round(-math.Sqrt2 / 2 * x * cw)
	hx := math.Round(math.Sqrt2 / 2 * x * ch)
	return image.Point{int(cw + y*cw + wx), int(ch - z*ch + hx)}
}

// Normalize the geometries of atoms such that the largest coordinate
// has a magnitude of 0.75
func (o Output) NormalizeGeom() {
	scale := 2.0
	for i := range o.Geom {
		o.Geom[i].X /= scale * o.max
		o.Geom[i].Y /= scale * o.max
		o.Geom[i].Z /= scale * o.max
	}
}

func ReadOut(filename string) (out Output) {
	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	var (
		geom  bool
		coord [3]float64
	)
	atom := regexp.MustCompile(`^\s+[0-9]+`)
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.Contains(line, "ATOMIC COORDINATES"):
			geom = true
		case strings.Contains(line, "Bond lengths in Bohr"):
			geom = false
		case geom && atom.MatchString(line):
			fields := strings.Fields(line)
			for i, f := range fields[3:] {
				v, _ := strconv.ParseFloat(f, 64)
				coord[i] = v
				if vabs := math.Abs(v); vabs > out.max {
					out.max = math.Abs(v)
				}
			}
			out.Geom = append(out.Geom,
				Atom{fields[1], coord[0], coord[1], coord[2]})
		}
	}
	return
}

// plot axes

// identify the center of mass of the molecule

// translate the molecule so its center of mass is at the origin

// also have to translate the dipole vectors by this

// determine the moments of inertia to identify the rotation axes for
// each constant

// plot those axes

// identify plane of molecule, make that be plane of screen

// DrawCircle draws a circle of radius r at c
func DrawCircle(img *image.NRGBA, c image.Point, r int, color color.NRGBA) {
	var ex, ey int
	for y := c.Y - r; y <= c.Y+r; y++ {
		for x := c.X - r; x <= c.X+r; x++ {
			ex = x - c.X
			ey = y - c.Y
			if ex*ex+ey*ey < r*r {
				img.Set(x, y, color)
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
func DrawLine(img *image.NRGBA, from, to image.Point) int {
	// vertical line
	if from.X == to.X {
		for y := from.Y; y <= to.Y; y++ {
			img.Set(to.X, y, color.NRGBA{0, 0, 0, 255})
		}
		return to.Y - from.Y
	}
	// needed the precision from floating point here
	m := float64(to.Y-from.Y) / float64(to.X-from.X)
	b := float64(to.Y) - m*float64(to.X)
	for x := from.X; x <= to.X; x++ {
		img.Set(x, int(m*float64(x)+b), color.NRGBA{0, 0, 0, 255})
	}
	x := from.X - to.X
	y := from.Y - to.Y
	return int(math.Sqrt(float64(x*x + y*y)))
}

// PlotAxes draws axes onto img using DrawLine
func PlotAxes(img *image.NRGBA) {
	DrawLine(img, image.Point{width / 2, 0}, image.Point{width / 2, height})
	DrawLine(img,
		image.Point{0, height / 2},
		image.Point{width, height / 2},
	)
	// length of each part of the x axis = (width/2)*sqrt(2)/2
	// where width and height can be used interchangeably if they
	// are the same, need the sqrt(w^2 + h^2)/2 if they are
	// different
	sw := int(math.Round(float64(width/2) * math.Sqrt2 / 2))
	DrawLine(img,
		image.Point{width/2 - sw, height/2 + sw},
		image.Point{width/2 + sw, height/2 - sw},
	)
}

// I have swap and invert, but need to apply them programatically
// instead of by eyeballing

func main() {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	PlotAxes(img)
	out := ReadOut("tests/dip.out")
	out.NormalizeGeom()
	for _, atom := range out.Geom {
		// atom = atom.Swap(X, Z)
		atom = atom.Invert(Z)
		fmt.Println(atom)
		pt := Cart2D(atom.X, atom.Y, atom.Z)
		DrawCircle(img, pt, 5, ptable[atom.Symbol])
	}
	f, _ := os.Create("test.png")
	png.Encode(f, img)
	f.Close()
}
