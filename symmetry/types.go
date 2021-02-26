package symm

import "fmt"

type axis int

// Cartesian axes/indices
const (
	X axis = iota
	Y
	Z
)

// Cylindrical coordinate indices, use Z from Cartesians
const (
	R axis = iota
	T
)

// return the plane perpendicular to the axis
func (a axis) not() plane {
	switch a {
	case X:
		return plane{Y, Z}
	case Y:
		return plane{X, Z}
	case Z:
		return plane{X, Y}
	default:
		panic("axis.not: invalid axis")
	}
}

func (a axis) String() string {
	return []string{"X", "Y", "Z"}[int(a)]
}

type plane struct {
	a, b axis
}

// return the axis perpindicular to the plane
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
