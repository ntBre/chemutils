package spectro

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"os/exec"
	"path"

	"io"

	"github.com/ntBre/chemutils/summarize"
)

var (
	atmNum = map[string]float64{
		"H": 1, "HE": 2, "LI": 3,
		"BE": 4, "B": 5, "C": 6,
		"N": 7, "O": 8, "F": 9,
		"NE": 10, "NA": 11, "MG": 12,
		"AL": 13, "SI": 14, "P": 15,
		"S": 16, "CL": 17, "AR": 18,
	}
	// SpectroCommand is the command used to run spectro
	SpectroCommand string
)

// Spectro holds the information for a spectro input file
type Spectro struct {
	Head     string // input directives
	Geometry string
	Body     string // weight and curvil
	Fermi1   string
	Fermi2   string
	Polyad   string
	Coriol   string
	Nfreqs   int
}

// LoadSpectro loads a spectro input file, assumes no resonances included
func LoadSpectro(filename string, names []string, coords string) (*Spectro, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var (
		buf  bytes.Buffer
		line string
		sp   Spectro
		geom bool
	)
	for scanner.Scan() {
		line = scanner.Text()
		if !geom {
			buf.WriteString(line + "\n")
		}
		if strings.Contains(line, "GEOM") {
			sp.Head = buf.String()
			buf.Reset()
			geom = true
		}
		// order agnostic
		if geom && (strings.Contains(line, "WEIGHT") ||
			strings.Contains(line, "CURVIL")) {
			geom = false
			buf.WriteString(line + "\n")
		}
	}
	sp.Body = buf.String()
	sp.FormatGeom(names, coords)
	return &sp, nil
}

// FormatGeom formats a slice of atom names and their corresponding
// coordinates for use in spectro
func (s *Spectro) FormatGeom(names []string, coords string) {
	// atomic numbers are 5.2f, 18.9f on coords
	var buf bytes.Buffer
	lines := strings.Split(coords, "\n")
	fmt.Fprintf(&buf, "%4d%4d\n", len(names), 1)
	for n := range names {
		fields := strings.Fields(lines[n])
		fmt.Fprintf(&buf, "%5.2f%18s%18s%18s\n",
			atmNum[strings.ToUpper(names[n])],
			fields[0], fields[1], fields[2])
	}
	s.Geometry = buf.String()
}

// WriteInput writes a Spectro to an input file for use
func (s *Spectro) WriteInput(filename string) error {
	var buf bytes.Buffer
	buf.WriteString(s.Head)
	buf.WriteString(s.Geometry)
	buf.WriteString(s.Body)
	buf.WriteString("# CORIOL #####\n")
	if s.Coriol != "" {
		buf.WriteString(s.Coriol)
	} else {
		fmt.Fprintf(&buf, "%5d\n", 0)
	}
	buf.WriteString("# FERMI1 ####\n")
	if s.Fermi1 != "" {
		buf.WriteString(s.Fermi1)
	} else {
		fmt.Fprintf(&buf, "%5d\n", 0)
	}
	buf.WriteString("# FERMI2 ####\n")
	if s.Fermi2 != "" {
		buf.WriteString(s.Fermi2)
	} else {
		fmt.Fprintf(&buf, "%5d\n", 0)
	}
	buf.WriteString("# RESIN ####\n")
	if s.Polyad != "" {
		buf.WriteString(s.Polyad)
	} else {
		fmt.Fprintf(&buf, "%5d\n", 0)
	}
	return ioutil.WriteFile(filename, buf.Bytes(), 0755)
}

// ParseCoriol parse a coriolis resonance from a spectro
// output line
func ParseCoriol(line string) string {
	letter := regexp.MustCompile(`A|B|C`)
	fields := strings.Fields(line)[1:3]
	// letters are only one character, so just take start index
	abcIndex := letter.FindStringIndex(fields[0])[0]
	abc := string(fields[0][abcIndex])
	switch {
	case abc == "A":
		abc = fmt.Sprintf("%5d%5d%5d", 1, 0, 0)
	case abc == "B":
		abc = fmt.Sprintf("%5d%5d%5d", 0, 1, 0)
	case abc == "C":
		abc = fmt.Sprintf("%5d%5d%5d", 0, 0, 1)
	}
	i := string(fields[0][:abcIndex])
	j := fields[1]
	return fmt.Sprintf("%5s%5s%s\n%5d\n", i, j, abc, 0)
}

// ParseFermi1 parses a type 1 fermi resonance from a spectro output line
func ParseFermi1(line string) string {
	fields := strings.Fields(line)[2:4]
	return fmt.Sprintf("%5s%5s\n", fields[0], fields[1])
}

// ParseFermi2 parses a type 2 fermi resonance from a spectro output line
func ParseFermi2(line string) string {
	fields := strings.Fields(line)[1:4]
	return fmt.Sprintf("%5s%5s%5s\n", fields[0], fields[1], fields[2])
}

// ReadOutput reads a spectro output and prepares resonance fields
// for rerunning spectro
func (s *Spectro) ReadOutput(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var (
		// buf  bytes.Buffer
		line        string
		coriol      bool
		fermi1      bool
		fermi2      bool
		skip        int
		coriolCount int
		fermi1Count int
		fermi2Count int
		freqs       bool
	)
	freq := regexp.MustCompile(`([0-9]+)\*?`)
	for scanner.Scan() {
		line = scanner.Text()
		switch {
		case strings.Contains(line, "FUNDAMENTAL"):
			freqs = true
		case freqs:
			fields := strings.Fields(line)
			if len(fields) == 4 &&
				freq.MatchString(fields[0]) {
				s.Nfreqs, _ = strconv.Atoi(
					freq.ReplaceAllString(
						fields[0], "$1",
					))
			}
		case strings.Contains(line, "DUNHAM"):
			freqs = false
		case skip > 0:
			skip--
			continue
		case coriol:
			if line == "" ||
				strings.Contains(line, "NO MODES FOUND") {
				coriol = false
			} else {
				coriolCount++
				s.Coriol += ParseCoriol(line)
			}
		case fermi1:
			if line == "" ||
				strings.Contains(line, "NOT FOUND") {
				fermi1 = false
			} else {
				fermi1Count++
				s.Fermi1 += ParseFermi1(line)
			}
		case fermi2:
			if line == "" ||
				strings.Contains(line, "NOT FOUND") {
				fermi2 = false
			} else {
				fermi2Count++
				s.Fermi2 += ParseFermi2(line)
			}
		case strings.Contains(line, "CORIOLIS RESONANCES"):
			skip = 3
			coriol = true
			// avoid fermi resonance in other contexts
		case strings.Contains(line, "  FERMI RESONANCE  "):
			fields := strings.Fields(line)
			if fields[3] == "1" {
				fermi1 = true
			} else {
				fermi2 = true
			}
			skip = 3
		}
	}
	// prepend the counts
	s.Coriol = fmt.Sprintf("%5d\n%5d\n", coriolCount, 0) + s.Coriol
	s.Fermi1 = fmt.Sprintf("%5d\n", fermi1Count) + s.Fermi1
	s.Fermi2 = fmt.Sprintf("%5d\n", fermi2Count) + s.Fermi2
	s.CheckPolyad()
}

// CheckPolyad checks for Fermi Polyads and set the Polyad field of s
// as necessary
func (s *Spectro) CheckPolyad() {
	f1 := strings.Split(s.Fermi1, "\n")
	f2 := strings.Split(s.Fermi2, "\n")
	rhSet := make(map[int]bool)
	lhSet := make(map[string]bool)
	var poly bool
	// skip count line and empty last split
	for _, line := range f1[1 : len(f1)-1] {
		lhs, rhs := EqnSeparate(line)
		if !rhSet[rhs] {
			rhSet[rhs] = true
		}
		if !lhSet[MakeKey(lhs)] {
			lhSet[MakeKey(lhs)] = true
		}
	}
	for _, line := range f2[1 : len(f2)-1] {
		lhs, rhs := EqnSeparate(line)
		if rhSet[rhs] {
			poly = true
		} else {
			rhSet[rhs] = true
		}
		if !lhSet[MakeKey(lhs)] {
			lhSet[MakeKey(lhs)] = true
		}
	}
	if !poly {
		return
	}
	var (
		resin string
		count int
	)
	for k := range rhSet {
		resin += ResinLine(s.Nfreqs, 1, k)
		count++
	}
	for k := range lhSet {
		num := 1
		ints := make([]int, 0)
		for _, f := range strings.Fields(k) {
			i, _ := strconv.Atoi(f)
			ints = append(ints, i)
		}
		if len(ints) == 1 {
			num = 2
		}
		resin += ResinLine(s.Nfreqs, num, ints...)
		count++
	}
	s.Polyad = fmt.Sprintf("%5d\n%5d\n%s", 1, count, resin)
}

// ResinLine formats a frequency number as a spectro RESIN line
func ResinLine(nfreqs, fill int, freqs ...int) string {
	var (
		buf   bytes.Buffer
		wrote bool
	)
	for j := 1; j <= nfreqs; j++ {
		for _, i := range freqs {
			if i == j {
				fmt.Fprintf(&buf, "%5d", fill)
				wrote = true
			}
		}
		if !wrote {
			fmt.Fprintf(&buf, "%5d", 0)
		}
		wrote = false
	}
	return buf.String() + "\n"
}

// MakeKey makes a mappable key from a slice of ints
func MakeKey(ints []int) string {
	var buf bytes.Buffer
	for i, v := range ints {
		fmt.Fprintf(&buf, "%d", v)
		if i < len(ints)-1 {
			fmt.Fprint(&buf, " ")
		}
	}
	return buf.String()
}

// EqnSeparate separates the fields of a spectro Fermi resonance into
// a  left- and right-hand side
func EqnSeparate(line string) (lhs []int, rhs int) {
	fields := strings.Fields(line)
	last := len(fields) - 1
	ints := make([]int, last+1)
	for i, val := range fields {
		val, _ := strconv.Atoi(val)
		ints[i] = val
	}
	lhs = append(lhs, ints[:last]...)
	rhs = ints[last]
	return
}

// RunSpectro runs SpectroCommand on filename
func RunSpectro(filename string) (err error) {
	file := path.Base(filename)
	cmd := exec.Command(SpectroCommand)
	cmd.Dir = path.Dir(filename)
	cmd.Stdin, err = os.Open(filepath.Join(cmd.Dir, file+".in"))
	if err != nil {
		return fmt.Errorf("stdin: %w", err)
	}
	var stdout io.Reader
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout: %w", err)
	}
	err = cmd.Start()
	outfile, err := os.Create(filepath.Join(cmd.Dir, file+".out"))
	if err != nil {
		return err
	}
	nbytes, err := io.Copy(outfile, stdout)
	if err != nil {
		return fmt.Errorf("%w after %d bytes", err, nbytes)
	}
	cmd.Wait()
	return nil
}

// UpdateHeader turns on resonance accounting if the resonance fields
// of s are non-empty
func (s *Spectro) UpdateHeader() {
	lines := strings.Split(s.Head, "\n")
	top := lines[0]
	bot := lines[3]
	fields := strings.Fields(lines[1] + lines[2])
	if s.Coriol != "" {
		fields[15] = "1"
	}
	if s.Fermi1 != "" {
		if s.Polyad != "" {
			fields[16] = "4"
		} else {
			fields[16] = "1"
		}
	}
	if s.Fermi2 != "" {
		if s.Polyad != "" {
			fields[17] = "4"
		} else {
			fields[17] = "1"
		}
	}
	var str strings.Builder
	str.WriteString(top + "\n")
	for i, f := range fields {
		fmt.Fprintf(&str, "%5s", f)
		if (i+1)%15 == 0 {
			fmt.Fprint(&str, "\n")
		}
	}
	str.WriteString(bot + "\n")
	s.Head = str.String()
}

// DoSpectro runs spectro dir, assuming there are nharms harmonic
// frequencies
func (s *Spectro) DoSpectro(dir string) (float64, []float64, []float64, []float64) {
	err := s.WriteInput(filepath.Join(dir, "spectro.in"))
	if err != nil {
		panic(err)
	}
	err = RunSpectro(filepath.Join(dir, "spectro"))
	if err != nil {
		panic(err)
	}
	s.ReadOutput(filepath.Join(dir, "spectro.out"))
	s.UpdateHeader()
	err = s.WriteInput(filepath.Join(dir, "spectro2.in"))
	if err != nil {
		panic(err)
	}
	err = RunSpectro(filepath.Join(dir, "spectro2"))
	if err != nil {
		panic(err)
	}
	res := summarize.Spectro(filepath.Join(dir, "spectro2.out"))
	return res.ZPT, res.Harm, res.Fund, res.Corr
}
