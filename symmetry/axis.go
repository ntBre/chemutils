package symm

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

type abm struct {
	axis    Axis
	weights [3]float64
}

type AxByMass []abm

func (b AxByMass) Len() int {
	return len(b)
}

func (b AxByMass) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b AxByMass) Less(i, j int) bool {
	// let the weight of the plane be the product of the masses of
	// its components
	ia, ib := b[i].weights[b[i].axis], b[i].weights[b[i].axis]
	ja, jb := b[j].weights[b[j].axis], b[j].weights[b[j].axis]
	return ia*ib < ja*jb
}

func (b AxByMass) Axes() (ret []Axis) {
	for _, v := range b {
		ret = append(ret, v.axis)
	}
	return
}
