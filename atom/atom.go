// Package atom provides the Atom type and associated functions for working
// with atoms and their coordinates
package atom

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// Cartesian axes/indices
const (
	X axis = iota
	Y
	Z
)

type axis int

func (a axis) String() string {
	return []string{"X", "Y", "Z"}[int(a)]
}

type plane struct {
	a, b axis
}

func (p plane) not() axis {
	switch {
	case p.a == Y && p.b == Z || p.a == Z && p.b == Y:
		return X
	case p.a == X && p.b == Z || p.a == Z && p.b == X:
		return Y
	case p.a == X && p.b == Y || p.a == Y && p.b == X:
		return Z
	default:
		panic("plane.not: invalid axis")
	}
}

func (p plane) String() string {
	return fmt.Sprintf("{%s, %s}", p.a, p.b)
}

// Cylindrical coordinate indices, use Z from Cartesians
const (
	R = iota
	T
)

// Float comparison threshold
const (
	eps = 1e-15
)

// Atom represents an atom, with atomic symbol and Cartesian coordinate
type Atom struct {
	Label string
	Coord []float64
}

// ReadXYZ reads an .xyz geometry from r and returns a slice of Atoms.
// If header is true, skip the number of atoms and comment
// lines. Otherwise, assume coordinates start at the first line.
func ReadXYZ(r io.Reader, header bool) []Atom {
	scanner := bufio.NewScanner(r)
	var line string
	atoms := make([]Atom, 0)
	for i := 1; scanner.Scan(); i++ {
		if header && i < 3 {
			continue
		}
		line = scanner.Text()
		atom, _ := Atomize(line)
		atoms = append(atoms, atom)
	}
	return atoms
}

// Atomize takes a string line of a Cartesian geometry and returns an
// Atom with that label and coordinate
func Atomize(line string) (Atom, error) {
	atom := new(Atom)
	fields := strings.Fields(line)
	if len(fields) != 4 {
		return *atom, fmt.Errorf("atomize: input (%q) too short", line)
	}
	atom.Label = fields[0]
	var (
		err error
		f   float64
	)
	for _, v := range fields[1:] {
		f, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return *atom, fmt.Errorf("atomize: %v", err)
		}
		atom.Coord = append(atom.Coord, f)
	}
	return *atom, nil
}

// Rotate returns a copy of atoms, with its coordinates rotated by deg
// degrees about axis
func Rotate(atoms []Atom, deg float64, axis axis) []Atom {
	rad := deg * math.Pi / 180.0
	new := make([]Atom, 0, len(atoms))
	for i := range atoms {
		cyl := ToCylinder(atoms[i].Coord, axis)
		cyl[T] += rad
		new = append(new, Atom{Label: atoms[i].Label,
			Coord: ToCartesian(cyl, axis)})
	}
	return new
}

// Mirror returns a copy of atoms, with its coordinates mirrored
// across plane
func Mirror(atoms []Atom, plane plane) []Atom {
	ax := plane.not()
	new := make([]Atom, len(atoms))
	for i, atom := range atoms {
		new[i] = atom
		new[i].Coord = make([]float64, len(atom.Coord))
		copy(new[i].Coord, atom.Coord)
		new[i].Coord[ax] = -new[i].Coord[ax]
	}
	return new
}

// approxEqual checks approximate equality between float slices
func approxEqual(x, y []float64) bool {
	if len(x) != len(y) {
		panic("approxEqual: dimension mismatch")
	}
	for i := range x {
		if math.Abs(x[i]-y[i]) > eps {
			return false
		}
	}
	return true
}

// ToCylinder transforms Cartesian coordinates to cylindrical
// coordinates of the form (r, theta, z)
func ToCylinder(coords []float64, axis axis) []float64 {
	var (
		r, t, z float64
	)
	if len(coords) != 3 {
		panic(fmt.Errorf("tocylinder: wrong number of coords (%d/3)", len(coords)))
	}
	switch axis {
	case X:
		r = math.Hypot(coords[Y], coords[Z])
		t = math.Atan2(coords[Z], coords[Y])
		z = coords[X]
	case Y:
		r = math.Hypot(coords[Z], coords[X])
		t = math.Atan2(coords[X], coords[Z])
		z = coords[Y]
	case Z:
		r = math.Hypot(coords[X], coords[Y])
		t = math.Atan2(coords[Y], coords[X])
		z = coords[Z]
	default:
		panic(fmt.Errorf("tocylinder: improper axis"))
	}
	return []float64{r, t, z}
}

// ToCartesian transforms coordinates from cylindrical coordinates of
// the form (r, t, z) to Cartesians of the form (x, y, z)
func ToCartesian(coords []float64, axis axis) (res []float64) {
	var x, y, z float64
	if len(coords) != 3 {
		panic(fmt.Errorf("tocartesian: wrong number of coords (%d/3)", len(coords)))
	}
	x = coords[0] * math.Cos(coords[1])
	y = coords[0] * math.Sin(coords[1])
	z = coords[2]
	switch axis {
	case X:
		res = []float64{z, x, y}
	case Y:
		res = []float64{y, z, x}
	case Z:
		res = []float64{x, y, z}
	}
	for i := range res {
		if math.Abs(res[i]) < eps {
			res[i] = 0
		}
	}
	return res
}
