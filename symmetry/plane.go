package symm

import "fmt"

type Plane struct {
	a, b Axis
}

// Planes
var (
	XY = Plane{X, Y}
	XZ = Plane{X, Z}
	YZ = Plane{Y, Z}
)

// return the Axis perpendicular to the Plane
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

type bm struct {
	plane   Plane
	weights [3]float64
}

type ByMass []bm

func (b ByMass) Len() int {
	return len(b)
}

func (b ByMass) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByMass) Less(i, j int) bool {
	// let the weight of the plane be the product of the masses of
	// its components
	ia, ib := b[i].weights[b[i].plane.a], b[i].weights[b[i].plane.b]
	ja, jb := b[j].weights[b[j].plane.a], b[j].weights[b[j].plane.b]
	return ia*ib < ja*jb
}

func (b ByMass) Planes() (ret []Plane) {
	for _, v := range b {
		ret = append(ret, v.plane)
	}
	return
}
