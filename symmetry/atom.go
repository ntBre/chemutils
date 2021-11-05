package symm

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// Atom represents an atom, with atomic symbol and Cartesian coordinate
type Atom struct {
	Label string
	Coord []float64
}

func (a Atom) String() string {
	return fmt.Sprintf("%-2s%10.6f%10.6f%10.6f\n",
		a.Label, a.Coord[0], a.Coord[1], a.Coord[2])
}

func Coords(atoms []Atom) (ret []float64) {
	for _, atom := range atoms {
		ret = append(ret, atom.Coord...)
	}
	return
}

func IsSame(atoms, btoms []Atom) bool {
	// resume testing here
	for _, atom := range atoms {
		found := false
		for b, btom := range btoms {
			if approxEqual(atom.Coord, btom.Coord) &&
				atom.Label == btom.Label {
				found = true
				// pop btoms[b] out of btoms
				lr := len(btoms) - 1
				btoms[b], btoms[lr] = btoms[lr], btoms[b]
				btoms = btoms[:lr]
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// IsRotAxis returns whether or not axis is a C_{360/deg} rotation
// axis for atoms
func IsRotAxis(atoms []Atom, deg float64, axis Axis) bool {
	return IsSame(atoms, Rotate(atoms, deg, axis))
}

// IsRefPlane returns whether or not plane is a reflection plane for
// atoms
func IsRefPlane(atoms []Atom, plane Plane) bool {
	return IsSame(atoms, Reflect(atoms, plane))
}

// Contains returns the two planes containing ax
func Contains(ax Axis) (sxz, syz Plane) {
	not := ax.Not()
	if ax < not.a {
		sxz = Plane{ax, not.a}
	} else {
		sxz = Plane{not.a, ax}
	}
	if ax < not.b {
		syz = Plane{ax, not.b}
	} else {
		syz = Plane{not.b, ax}
	}
	return
}

// Symmetry returns the Irrep corresponding to atoms within the point
// group, principal axis, and major plane defined by m
func (m Molecule) Symmetry(atoms []Atom) Irrep {
	// can I make this more generic? some kind of data structure
	// holding functions to loop over
	switch m.Group {
	case Cs:
		if IsRefPlane(atoms, m.Main) {
			return Ap
		}
		return App
	case C2v:
		sxz, syz := Contains(m.Principal)
		if IsRotAxis(atoms, 180.0, m.Principal) {
			// A something
			if IsRefPlane(atoms, sxz) {
				if IsRefPlane(atoms, syz) {
					return A1
				} else {
					panic("impossible irrep")
				}
			} else {
				return A2
			}
		} else {
			// B something
			if IsRefPlane(atoms, sxz) && !IsRefPlane(atoms, syz) {
				return B1
			} else if IsRefPlane(atoms, syz) && !IsRefPlane(atoms, sxz) {
				return B2
			}
		}
		return A1
	default:
		panic("Unrecognized point group")
	}
}

func Negate(atoms []Atom) []Atom {
	ret := make([]Atom, 0, len(atoms))
	for _, atom := range atoms {
		coords := make([]float64, 0, 3)
		for _, c := range atom.Coord {
			coords = append(coords, -1*c)
		}
		ret = append(ret, Atom{
			Label: atom.Label,
			Coord: coords,
		})
	}
	return ret
}

// PointGroup determines the point group of mol
func PointGroup(mol Molecule) (ret Group) {
	// check for rotation axis first
	if IsRotAxis(mol.Atoms, 180.0, mol.Principal) {
		ret = C2
		if IsRefPlane(mol.Atoms, mol.Main) {
			if mol.Main.a == mol.Principal ||
				mol.Main.b == mol.Principal {
				// => sigma_v
				ret = C2v
			}
		}
	} else if IsRefPlane(mol.Atoms, mol.Main) {
		ret = Cs
	}
	return
}

// ReadXYZ reads an .xyz geometry from r and returns a slice of Atoms.
// If the first line looks like the number of atoms skip it and the
// comment line. Otherwise start reading coordinates directly.
func ReadXYZ(r io.Reader) (ret Molecule) {
	scanner := bufio.NewScanner(r)
	var (
		line string
		skip int
		sums [3]float64
	)
	for i := 1; scanner.Scan(); i++ {
		line = scanner.Text()
		switch {
		case skip > 0:
			skip--
		case i == 1 && len(strings.Fields(line)) == 1:
			skip = 1
		default:
			atom := new(Atom)
			fields := strings.Fields(line)
			if len(fields) != 4 {
				continue
			}
			atom.Label = fields[0]
			var (
				err error
				f   float64
			)
			for i, v := range fields[1:] {
				f, err = strconv.ParseFloat(v, 64)
				if err != nil {
					log.Fatalf("atomize: %v", err)
				}
				atom.Coord = append(atom.Coord, f)
				f *= ptable[atom.Label]
				sums[i] += math.Abs(f)
			}
			ret.Atoms = append(ret.Atoms, *atom)
		}
	}
	fst := sums[0]
	var a, b Axis
	for i, v := range sums {
		switch {
		case v > fst:
			fst = v
			a, b = Axis(i), a
		}
	}
	ret.Principal = a
	ret.Main = Plane{b, a}
	ret.Group = PointGroup(ret)
	return
}

// LoadXYZ is a convenience function for calling ReadXYZ on a file
func LoadXYZ(filename string) Molecule {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	return ReadXYZ(f)
}
