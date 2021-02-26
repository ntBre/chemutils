package symm

import "fmt"

type Axis int

// Cartesian axes
const (
	X Axis = iota
	Y
	Z
)

// Cylindrical axes
const (
	R Axis = iota
	T
)

// Planes
var (
	XY = Plane{X, Y}
	XZ = Plane{X, Z}
	YZ = Plane{Y, Z}
)

// return the Plane perpendicular to the Axis
func (a Axis) not() Plane {
	switch a {
	case X:
		return Plane{Y, Z}
	case Y:
		return Plane{X, Z}
	case Z:
		return Plane{X, Y}
	default:
		panic("Axis.not: invalid Axis")
	}
}

func (a Axis) String() string {
	return []string{"X", "Y", "Z"}[int(a)]
}

type Plane struct {
	a, b Axis
}

// return the Axis perpindicular to the Plane
func (p Plane) not() Axis {
	switch {
	case p.a == Y && p.b == Z || p.a == Z && p.b == Y:
		return X
	case p.a == X && p.b == Z || p.a == Z && p.b == X:
		return Y
	case p.a == X && p.b == Y || p.a == Y && p.b == X:
		return Z
	default:
		panic("Plane.not: invalid Axis")
	}
}

func (p Plane) String() string {
	return fmt.Sprintf("{%s, %s}", p.a, p.b)
}

type Molecule struct {
	Atoms     []Atom
	Principal Axis
	Main      Plane
}
