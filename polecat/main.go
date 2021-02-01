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
	GREEN = color.NRGBA{0, 255, 0, 255}
	BLUE  = color.NRGBA{0, 0, 255, 255}
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
func (a Atom) Coords() Vec {
	return Vec{a.X, a.Y, a.Z}
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

type Vec [3]float64

func (v Vec) Add(w Vec) (ret Vec) {
	for i := range v {
		ret[i] = v[i] + w[i]
	}
	return
}

func (v Vec) Sub(w Vec) (ret Vec) {
	for i := range v {
		ret[i] = v[i] - w[i]
	}
	return
}

func (v Vec) Mul(a float64) (ret Vec) {
	for i := range v {
		ret[i] = v[i] * a
	}
	return
}

func (v Vec) Dot(w Vec) (ret float64) {
	for i := range v {
		ret += v[i] * w[i]
	}
	return
}

func (v Vec) Cross(w Vec) (ret Vec) {
	return Vec{
		v[1]*w[2] - v[2]*w[1],
		v[2]*w[0] - v[0]*w[2],
		v[0]*w[1] - v[1]*w[0],
	}
}

func (v Vec) Size() (ret float64) {
	for i := range v {
		ret += v[i] * v[i]
	}
	return math.Sqrt(ret)
}

func (v Vec) Unit() Vec {
	size := v.Size()
	if size == 0 {
		return v
	}
	return v.Mul(1 / v.Size())
}

// Order returns the axes of v in descending order
func (v Vec) Order() []Axis {
	switch {
	case v[X] > v[Y] && v[X] > v[Z]:
		if v[Y] >= v[Z] {
			return []Axis{X, Y, Z}
		}
		return []Axis{X, Z, Y}
	case v[X] > v[Y]:
		return []Axis{Z, X, Y}
	case v[X] > v[Z]:
		return []Axis{Y, X, Z}
	case v[Y] > v[Z]:
		return []Axis{Y, Z, X}
	case v[Z] > v[Y]:
		return []Axis{Z, Y, X}
	default:
		panic("order not found")
	}
}

type Output struct {
	Geom []Atom
	max  float64
}

// Cart2D converts the point (x, y, z) to an image.Point
func Cart2D(vec Vec) image.Point {
	cw, ch := float64(width/2), float64(height/2)
	wx := math.Round(-math.Sqrt2 / 2 * vec[X] * cw)
	hx := math.Round(math.Sqrt2 / 2 * vec[X] * ch)
	return image.Point{int(cw + vec[Y]*cw + wx), int(ch - vec[Z]*ch + hx)}
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

// DONE plot axes

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
		if from.Y > to.Y {
			to, from = from, to
		}
		for y := from.Y; y <= to.Y; y++ {
			img.Set(to.X, y, color.NRGBA{0, 0, 0, 255})
		}
		return to.Y - from.Y
	}
	// needed the precision from floating point here
	m := float64(to.Y-from.Y) / float64(to.X-from.X)
	b := float64(to.Y) - m*float64(to.X)
	if from.X > to.X {
		to, from = from, to
	}
	for x := from.X; x <= to.X; x++ {
		img.Set(x, int(m*float64(x)+b), color.NRGBA{0, 0, 0, 255})
	}
	x := from.X - to.X
	y := from.Y - to.Y
	return int(math.Sqrt(float64(x*x + y*y)))
}

// Rodrigues applies Rodrigues' rotation formula to v using the unit
// vector k as the axis of rotation and theta as the rotation angle in
// radians
func Rodrigues(v, k Vec, theta float64) (ret Vec) {
	return v.Mul(math.Cos(theta)).
		Add(k.Cross(v).Mul(math.Sin(theta))).
		Add(k.Mul(k.Dot(v)).Mul(1 - math.Cos(theta)))
}

// DrawVec calls DrawLine and adds an arrow tip at to
func DrawVec(img *image.NRGBA, from, to Vec) int {
	// need to take this vector and do some trig on it to get the
	// tips
	v := to.Sub(from)
	// w is a vector perpendicular to v
	w := Vec{v[0], v[1], v[2]}
	ord := v.Order()
	w[ord[0]], w[ord[1]] = w[ord[1]], w[ord[0]]
	w = w.Mul(-1)
	if v.Dot(w) > 1e-4 {
		panic("nonzero dot product")
	}
	// k is a unit vector perpendicular to the v-w plane
	k := v.Cross(w).Unit()
	rod := Rodrigues(v, k, 7.5*math.Pi/6).Unit().Mul(0.1)
	mrod := Rodrigues(v, k, -7.5*math.Pi/6).Unit().Mul(0.1)
	DrawLine(img, Cart2D(to), Cart2D(to.Add(rod)))
	DrawLine(img, Cart2D(to), Cart2D(to.Add(mrod)))
	return DrawLine(img, Cart2D(from), Cart2D(to))
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
	// for _, atom := range out.Geom {
	// 	// atom = atom.Swap(X, Z)
	// 	atom = atom.Invert(Z)
	// 	// fmt.Println(atom)
	// 	pt := Cart2D(atom.Coords())
	// 	DrawCircle(img, pt, 5, ptable[atom.Symbol])
	// }
	DrawVec(img, Vec{0, 0.5, 0}, Vec{0.5, 0.5, 0})
	DrawVec(img, Vec{0, 0.5, 0}, Vec{0, 0.5, 0.5})
	f, _ := os.Create("test.png")
	png.Encode(f, img)
	f.Close()
}
