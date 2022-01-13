package symm

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Atom represents an atom, with atomic symbol and Cartesian coordinate
type Atom struct {
	Label string
	Coord []float64
}

func (a Atom) String() string {
	return fmt.Sprintf("%-2s%15.10f%15.10f%15.10f\n",
		a.Label, a.Coord[0], a.Coord[1], a.Coord[2])
}

func Labels(atoms []Atom) []string {
	ret := make([]string, len(atoms))
	for i, atom := range atoms {
		ret[i] = atom.Label
	}
	return ret
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
			if ApproxEqual(atom.Coord, btom.Coord) &&
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
	case C1:
		return E
	case C2:
		if IsRotAxis(atoms, 180.0, m.Axes[0]) {
			return A
		}
		return B
	case Cs:
		if IsRefPlane(atoms, m.Planes[0]) {
			return Ap
		}
		return App
	case C2v:
		sxz, syz := Contains(m.Axes[0])
		if IsRotAxis(atoms, 180.0, m.Axes[0]) {
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
	if len(mol.Axes) > 0 && IsRotAxis(mol.Atoms, 180.0, mol.Axes[0]) {
		ret = C2
		if len(mol.Planes) > 0 && IsRefPlane(mol.Atoms, mol.Planes[0]) &&
			(mol.Planes[0].a == mol.Axes[0] ||
				mol.Planes[0].b == mol.Axes[0]) {
			// => sigma_v
			ret = C2v
		}
	} else if len(mol.Planes) > 0 && IsRefPlane(mol.Atoms, mol.Planes[0]) {
		ret = Cs
	}
	return
}
func ReadAtoms(r io.Reader) (ret []Atom, sums [3]float64) {
	scanner := bufio.NewScanner(r)
	var (
		line string
		skip int
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
			ret = append(ret, *atom)
		}
	}
	return
}

// ReadXYZ reads an .xyz geometry from r and returns a Molecule. If
// the first line looks like the number of atoms skip it and the
// comment line. Otherwise start reading coordinates directly.
func ReadXYZ(r io.Reader) (ret Molecule) {
	var sums [3]float64
	ret.Atoms, sums = ReadAtoms(r)
	// Find all C2 axes and mirror planes
	var axes AxByMass
	for _, a := range []Axis{X, Y, Z} {
		if IsRotAxis(ret.Atoms, 180.0, a) {
			axes = append(axes, abm{a, sums})
		}
	}
	var planes ByMass
	for _, p := range []Plane{XY, XZ, YZ} {
		if IsRefPlane(ret.Atoms, p) {
			planes = append(planes, bm{p, sums})
		}
	}
	// Sort the axes and mirror planes by mass involved
	sort.Sort(sort.Reverse(axes))
	ret.Axes = axes.Axes()
	sort.Sort(sort.Reverse(planes))
	ret.Planes = planes.Planes()
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
