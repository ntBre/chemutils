// (wobuf "sxiv test.png")
// (zap "~/School/Research/Pubs/Oxywater/" "main")
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
	"sort"
	"strconv"
	"strings"

	"flag"

	"gonum.org/v1/gonum/mat"
)

// flags
var (
	noax = flag.Bool("noax", false, "turn off axes")
)

const (
	h    = 6.62607554e-34 // m^2 kg/s
	c    = 299792458.0    // m/s
	toCm = 1 / 1.6605402e-27 / 5.29177249e-11 / 5.29177249e-11 / 100 / c
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
	RED    = color.NRGBA{255, 0, 0, 255}
	GREEN  = color.NRGBA{0, 255, 0, 255}
	BLUE   = color.NRGBA{0, 0, 255, 255}
	BLACK  = color.NRGBA{0, 0, 0, 255}
	ORANGE = color.NRGBA{252, 186, 3, 255}

	Origin = Vec{0, 0, 0}
)

type Element struct {
	Mass  float64
	Size  int
	Color color.NRGBA
}

var ptable = map[string]Element{
	"H": {
		Mass:  1.00782503223,
		Size:  6,
		Color: color.NRGBA{0, 0, 0, 64},
	},
	"C": {
		Mass:  12.0,
		Size:  8,
		Color: BLACK,
	},
	"O": {
		Mass:  15.99491462957,
		Size:  8,
		Color: RED,
	},
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

func (a Atom) Translate(v Vec) Atom {
	a.X -= v[X]
	a.Y -= v[Y]
	a.Z -= v[Z]
	return a
}

func (a Atom) Dist(b Atom) float64 {
	vec := a.Coords().Sub(b.Coords())
	return vec.Size()
}

func (a Atom) String() string {
	return fmt.Sprintf("%s %12.8f %12.8f %12.8f",
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

// Output is a type for holding information from a Molpro dipole
// calculation output file. Geom is the cartesian geometry expressed
// as a slice of Atoms; Dip{x,y,z} are the dipole moments (in AU) in
// the respective directions; and max is the maximum value in the
// geometry, used for normalization.
type Output struct {
	Geom []Atom
	Dipx float64
	Dipy float64
	Dipz float64
	max  float64
}

func (o Output) Dips() []float64 {
	return []float64{o.Dipx, o.Dipy, o.Dipz}
}

// Cart2D converts the point (x, y, z) to an image.Point
func Cart2D(vec Vec) image.Point {
	cw, ch := float64(width/2), float64(height/2)
	wx := math.Round(-math.Sqrt2 / 2 * vec[X] * cw)
	hx := math.Round(math.Sqrt2 / 2 * vec[X] * ch)
	return image.Point{int(cw + vec[Y]*cw + wx), int(ch - vec[Z]*ch + hx)}
}

// Normalize the geometries and dipoles of atoms such that the largest
// coordinate has a magnitude of 0.5
func (o *Output) Normalize() {
	scale := 2.0
	for i := range o.Geom {
		o.Geom[i].X /= scale * o.max
		o.Geom[i].Y /= scale * o.max
		o.Geom[i].Z /= scale * o.max
	}
	dips := []float64{o.Dipx, o.Dipy, o.Dipz}
	sort.Sort(sort.Reverse(sort.Float64Slice(dips)))
	max := dips[0]
	o.Dipx /= scale * max
	o.Dipy /= scale * max
	o.Dipz /= scale * max
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
		case strings.Contains(line, "SETTING DIP"):
			fields := strings.Fields(line)
			// replace D scientific notation with E
			fields[3] = strings.Replace(fields[3], "D", "E", -1)
			switch fields[1] {
			case "DIPX":
				out.Dipx, _ = strconv.ParseFloat(fields[3], 64)
			case "DIPY":
				out.Dipy, _ = strconv.ParseFloat(fields[3], 64)
			case "DIPZ":
				out.Dipz, _ = strconv.ParseFloat(fields[3], 64)
			}
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

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// DrawLine draws a line from from to to
func DrawLine(img *image.NRGBA, color color.NRGBA, from, to image.Point) int {
	// // vertical line
	if from.X == to.X {
		if from.Y > to.Y {
			to, from = from, to
		}
		for y := from.Y; y <= to.Y; y++ {
			img.Set(to.X, y, color)
		}
		return to.Y - from.Y
	}
	// needed the precision from floating point here
	m := float64(to.Y-from.Y) / float64(to.X-from.X)
	b := float64(to.Y) - m*float64(to.X)
	// loop over the larger component
	if Abs(from.X-to.X) > Abs(from.Y-to.Y) {
		if from.X > to.X {
			to, from = from, to
		}
		for x := from.X; x <= to.X; x++ {
			img.Set(x, int(m*float64(x)+b), color)
		}
	} else {
		if from.Y > to.Y {
			to, from = from, to
		}
		// y = mx + b => (y - b)/m = x
		for y := from.Y; y <= to.Y; y++ {
			img.Set(int((float64(y)-b)/m), y, color)
		}
	}
	x := from.X - to.X
	y := from.Y - to.Y
	return int(math.Sqrt(float64(x*x + y*y)))
}

// Rodrigues applies Rodrigues' rotation formula to v using the unit
// vector k as the axis of rotation and theta as the rotation angle in
// radians. Idea from
// stackoverflow.com/questions/14607640/rotating-a-vector-in-3d-space
func Rodrigues(v, k Vec, theta float64) (ret Vec) {
	return v.Mul(math.Cos(theta)).
		Add(k.Cross(v).Mul(math.Sin(theta))).
		Add(k.Mul(k.Dot(v)).Mul(1 - math.Cos(theta)))
}

// DrawVec calls DrawLine and adds an arrow tip at to
func DrawVec(img *image.NRGBA, color color.NRGBA, from, to Vec) int {
	l := DrawLine(img, color, Cart2D(from), Cart2D(to))
	// for very short vectors, the head is larger so don't draw
	if l <= 1 {
		return l
	}
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
	DrawLine(img, color, Cart2D(to), Cart2D(to.Add(rod)))
	DrawLine(img, color, Cart2D(to), Cart2D(to.Add(mrod)))
	return l
}

// PlotAxes draws axes onto img using DrawLine
func PlotAxes(img *image.NRGBA) {
	DrawLine(img, BLUE, image.Point{width / 2, 0}, image.Point{width / 2, height})
	DrawLine(img, GREEN,
		image.Point{0, height / 2},
		image.Point{width, height / 2},
	)
	// length of each part of the x axis = (width/2)*sqrt(2)/2
	// where width and height can be used interchangeably if they
	// are the same, need the sqrt(w^2 + h^2)/2 if they are
	// different
	sw := int(math.Round(float64(width/2) * math.Sqrt2 / 2))
	DrawLine(img, RED,
		image.Point{width/2 - sw, height/2 + sw},
		image.Point{width/2 + sw, height/2 - sw},
	)
}

// COM computes the center of mass vector for atoms
func COM(atoms []Atom) Vec {
	var (
		m, mtot, xcm, ycm, zcm float64
	)
	for _, atom := range atoms {
		m = ptable[atom.Symbol].Mass
		mtot += m
		xcm += m * atom.X
		ycm += m * atom.Y
		zcm += m * atom.Z
	}
	return Vec{
		xcm / mtot,
		ycm / mtot,
		zcm / mtot,
	}
}

// MOI computes the principal moments of inertia in amu * bohr^2 and
// returns the eigenvectors corresponding to each moment
func MOI(atoms []Atom) (ia, ib, ic float64, eva, evb, evc Vec) {
	moi := mat.NewDense(3, 3, nil)
	var (
		m    float64
		c    Vec
		f, g int
	)
	for _, atom := range atoms {
		m = ptable[atom.Symbol].Mass
		c = atom.Coords()
		// x,y,z -> 0,1,2
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if i == j {
					f = (i + 1) % 3
					g = (i + 2) % 3
					moi.Set(i, j, moi.At(i, j)+
						m*(c[f]*c[f]+
							c[g]*c[g]))
				} else {
					moi.Set(i, j, moi.At(i, j)-
						m*c[i]*c[j])
				}
			}
		}
	}
	fmt.Println("Moment of inertia tensor:")
	printMat(moi)
	var eig mat.Eigen
	ok := eig.Factorize(moi, mat.EigenRight)
	if !ok {
		panic("eigen decomposition failed")
	}
	dst := eig.Values(nil)
	// extract eigenvectors to determine what the values match to
	// columns should be the vectors
	var vecs mat.CDense
	eig.VectorsTo(&vecs)
	rows, cols := vecs.Dims()
	realVecs := mat.NewDense(rows, cols, nil)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			realVecs.Set(i, j,
				real(vecs.At(i, j)))
		}
	}
	fmt.Println("Eigenvectors")
	printMat(realVecs)
	evs := make([]Vec, 3)
	for i := 0; i < rows; i++ {
		col := realVecs.ColView(i)
		evs[i] = Vec{col.AtVec(0), col.AtVec(1), col.AtVec(2)}
	}
	for _, v := range dst {
		if imag(v) > 0 {
			panic("imaginary moment of inertia")
		}
	}
	return real(dst[0]), real(dst[1]), real(dst[2]), evs[0], evs[1], evs[2]
}

func printMat(matrix *mat.Dense) {
	r, _ := matrix.Dims()
	for row := 0; row < r; row++ {
		fmt.Printf("%14.8f%14.8f%14.8f\n",
			matrix.At(row, 0),
			matrix.At(row, 1),
			matrix.At(row, 2),
		)
	}
}

// Rot takes a principal moment of inertia and returns the
// associated principal rotational constant in cm-1
func Rot(I float64) float64 {
	return toCm * h / (8 * math.Pi * math.Pi * I)
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		panic("no input file given")
	}
	infile := args[0]
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	if !*noax {
		PlotAxes(img)
	}
	out := ReadOut(infile)
	com := COM(out.Geom)
	fmt.Println("Translated geometry:")
	for i := range out.Geom {
		out.Geom[i] = out.Geom[i].Translate(com)
		fmt.Println(out.Geom[i])
	}
	ia, ib, ic, eva, evb, evc := MOI(out.Geom)
	fmt.Println("Moments of inertia (amu bohr^2):")
	fmt.Printf("%14.8f%14.8f%14.8f\n", ia, ib, ic)
	fmt.Println("Rotational constants (cm-1):")
	fmt.Printf("%14.8f%14.8f%14.8f\n", Rot(ia), Rot(ib), Rot(ic))
	out.Normalize()
	for i := 0; i < len(out.Geom); i++ {
		a := Cart2D(out.Geom[i].Coords())
		for j := i + 1; j < len(out.Geom); j++ {
			dist := out.Geom[i].Dist(out.Geom[j])
			if dist < 1.0 {
				DrawLine(img, BLACK, a,
					Cart2D(out.Geom[j].Coords()))
			}
		}
		element := ptable[out.Geom[i].Symbol]
		DrawCircle(img, a, element.Size, element.Color)
	}
	// dipole vectors
	DrawVec(img, BLACK, Origin.Add(com), Vec{out.Dipx, 0, 0}.Add(com))
	DrawVec(img, BLACK, Origin.Add(com), Vec{0, out.Dipy, 0}.Add(com))
	DrawVec(img, BLACK, Origin.Add(com), Vec{0, 0, out.Dipz}.Add(com))
	null, _ := os.Open(os.DevNull)
	fmt.Fprint(null, eva, evb, evc)
	// moment of inertia axes
	// DrawVec(img, RED, Origin.Add(com), eva.Add(com))
	// DrawVec(img, GREEN, Origin.Add(com), evb.Add(com))
	// DrawVec(img, BLUE, Origin.Add(com), evc.Add(com))
	f, _ := os.Create("test.png")
	png.Encode(f, img)
	f.Close()
}
