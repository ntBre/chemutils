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
func (a Axis) Not() Plane {
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
	return []string{"X", "Y", "Z"}[a]
}

type Plane struct {
	a, b Axis
}

// return the Axis perpindicular to the Plane
func (p Plane) Not() Axis {
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
	Group     Group
}

type Group int

func (g Group) String() string {
	return []string{
		"C1",
		"Cs",
		"C2",
		"C2v",
	}[g]
}

const (
	C1 Group = iota
	Cs
	C2
	C2v
)

type Irrep int

// C2v symmetry elements
const (
	A1 Irrep = iota
	B2
	B1
	A2
)

func (i Irrep) String() string {
	return []string{
		"A1",
		"B2",
		"B1",
		"A2",
	}[i]
}
