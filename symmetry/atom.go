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

// IsRefPlane returns whether or not plane is a reflection plane for
// atoms
func IsRefPlane(atoms []Atom, plane Plane) bool {
	// TODO make this handle different atoms in same places
	ref := Coords(Reflect(atoms, plane))
	if approxEqual(ref, Coords(atoms)) {
		return true
	}
	return false
}

// IsRotAxis returns whether or not axis is a C_{360/deg} rotation
// axis for atoms
func IsRotAxis(atoms []Atom, deg float64, axis Axis) bool {
	rots := Rotate(atoms, deg, axis)
	for _, atom := range atoms {
		found := false
		for r, rot := range rots {
			if approxEqual(atom.Coord, rot.Coord) ||
				atom.Label == rot.Label {
				found = true
				// pop rots[r] out of rots
				lr := len(rots) - 1
				rots[r], rots[lr] = rots[lr], rots[r]
				rots = rots[:lr]
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
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
				log.Fatalf("atomize: input (%q) too short", line)
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
