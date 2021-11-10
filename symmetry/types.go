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

type Molecule struct {
	Atoms  []Atom
	Axes   []Axis
	Planes []Plane
	Group  Group
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

const (
	// C2v
	A1 Irrep = iota
	B2
	B1
	A2
	// Cs, p == prime
	Ap
	App
	// C1
	E
	// C2
	A
	B
)

func (i Irrep) String() string {
	return []string{
		"A1",
		"B2",
		"B1",
		"A2",
	}[i]
}
