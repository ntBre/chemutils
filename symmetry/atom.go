package symm

import (
	"bufio"
	"fmt"
	"io"
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

// Atomize takes a string line of a Cartesian geometry and returns an
// Atom with that label and coordinate
func Atomize(line string) (Atom, error) {
	atom := new(Atom)
	fields := strings.Fields(line)
	if len(fields) != 4 {
		return *atom,
			fmt.Errorf("atomize: input (%q) too short", line)
	}
	atom.Label = fields[0]
	var (
		err error
		f   float64
	)
	for _, v := range fields[1:] {
		f, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return *atom, fmt.Errorf("atomize: %v", err)
		}
		atom.Coord = append(atom.Coord, f)
	}
	return *atom, nil
}

// ReadXYZ reads an .xyz geometry from r and returns a slice of Atoms.
// If the first line looks like the number of atoms skip it and the
// comment line. Otherwise start reading coordinates directly.
func ReadXYZ(r io.Reader) []Atom {
	scanner := bufio.NewScanner(r)
	var (
		line string
		skip int
	)
	atoms := make([]Atom, 0)
	for i := 1; scanner.Scan(); i++ {
		line = scanner.Text()
		switch {
		case skip > 0:
			skip--
		case i == 1 && len(strings.Fields(line)) == 1:
			skip = 1
		default:
			atom, _ := Atomize(line)
			atoms = append(atoms, atom)
		}
	}
	return atoms
}

// LoadXYZ is a convenience function for calling ReadXYZ on a file
func LoadXYZ(filename string) []Atom {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	return ReadXYZ(f)
}
