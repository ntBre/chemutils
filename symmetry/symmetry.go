// Package symm is for determining molecular point group symmetry
package symm

import (
	"fmt"
	"math"
)

// Float comparison threshold
var (
	EPS = 1e-15
)

// Rotate returns a copy of atoms, with its coordinates rotated by deg
// degrees about axis
func Rotate(atoms []Atom, deg float64, axis Axis) []Atom {
	rad := deg * math.Pi / 180.0
	new := make([]Atom, len(atoms))
	for i := range atoms {
		cyl := ToCylinder(atoms[i].Coord, axis)
		cyl[T] += rad
		new[i] = Atom{
			Label: atoms[i].Label,
			Coord: ToCartesian(cyl, axis),
		}
	}
	return new
}

// Reflect returns a copy of atoms, with its coordinates mirrored
// across plane
func Reflect(atoms []Atom, plane Plane) []Atom {
	ax := plane.Not()
	new := make([]Atom, len(atoms))
	for i, atom := range atoms {
		new[i] = atom
		new[i].Coord = make([]float64, len(atom.Coord))
		copy(new[i].Coord, atom.Coord)
		new[i].Coord[ax] = -new[i].Coord[ax]
	}
	return new
}

// RotaryReflect returns a copy of atoms, with its coordinates rotated
// about axis and then mirrored through the plane perpendicular to
// axis
func RotaryReflect(atoms []Atom, deg float64, axis Axis) []Atom {
	rot := Rotate(atoms, deg, axis)
	pl := axis.Not()
	return Reflect(rot, pl)
}

// Invert uses the fact that S_2 = i to return a copy of atoms, with
// its coordinates inverted or equivalently rotated 180 degrees about
// ax and then mirrored through the plane perpendicular to ax
func Invert(atoms []Atom, ax Axis) []Atom {
	return RotaryReflect(atoms, 180.0, ax)
}

// ToCylinder transforms Cartesian coordinates to cylindrical
// coordinates of the form (r, theta, z)
func ToCylinder(coords []float64, axis Axis) []float64 {
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

// approxEqual checks approximate equality between float slices
func approxEqual(x, y []float64) bool {
	if len(x) != len(y) {
		panic("approxEqual: dimension mismatch")
	}
	for i := range x {
		if math.Abs(x[i]-y[i]) > EPS {
			return false
		}
	}
	return true
}

// ToCartesian transforms coordinates from cylindrical coordinates of
// the form (r, t, z) to Cartesians of the form (x, y, z)
func ToCartesian(coords []float64, axis Axis) (res []float64) {
	var x, y, z float64
	if len(coords) != 3 {
		panic(fmt.Errorf("tocartesian: wrong number of coords (%d/3)",
			len(coords)))
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
		if math.Abs(res[i]) < EPS {
			res[i] = 0
		}
	}
	return res
}
