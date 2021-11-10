package symm

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
