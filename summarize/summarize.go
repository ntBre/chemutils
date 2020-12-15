// Package summarize is for summarizing output from quantum chemistry
// programs
// (let ((compile-command "go test .")) (my-recompile))
package summarize

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Result struct {
	ZPT    float64
	Harm   []float64
	Fund   []float64
	Corr   []float64
	Rots   [][]float64
	Deltas []float64
	Phis   []float64
	Rhead  []string
	Ralpha []float64
	Requil []float64
	Fermi  []string
}

// FreqReport gathers harmonic, anharmonic, and resonance-corrected
// frequencies from a spectro  output file for reporting
func Spectro(filename string, nfreqs int) *Result {
	res := new(Result)
	res.Corr = make([]float64, nfreqs, nfreqs)
	fermiMap := make(map[string][]string)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var (
		line     string
		skip     int
		freqs    int
		nrot     int
		harmFund bool
		rot      bool
		geom     bool
		fermi1   bool
		fermi2   bool
		buf      bytes.Buffer
		// this is cute
		gparams = []string{"", "", "r", "<"}
	)
	freq := regexp.MustCompile(`[0-9]+`)
	delta := regexp.MustCompile(`(?i)delta (J|(JK)|K)`)
	phi := regexp.MustCompile(`(?i)phi (J|(JK)|(KJ)|K)`)
	icn := regexp.MustCompile(`\([ 0-9]+\)\s+(BOND|ANGLE)`)
	atom := regexp.MustCompile(`([0-9]+)\(([A-Za-z ]+)\)`)
	for scanner.Scan() {
		line = scanner.Text()
		switch {
		case skip > 0:
			skip--
		case strings.Contains(line, "BAND CENTER ANALYSIS"):
			skip += 3
			freqs = nfreqs
			harmFund = true
		case harmFund && freqs > 0 && len(line) > 1:
			fields := strings.Fields(line)
			if freq.MatchString(fields[0]) {
				h, _ := strconv.ParseFloat(fields[1], 64)
				f, _ := strconv.ParseFloat(fields[2], 64)
				res.Harm = append(res.Harm, h)
				res.Fund = append(res.Fund, f)
				freqs--
			}
			if freqs == 0 {
				harmFund = false
			}
		case strings.Contains(line, "STATE NO."):
			skip += 2
			freqs = nfreqs + 1 // add ZPT
		case !harmFund && freqs > 0 && len(line) > 1:
			fields := strings.Fields(line)
			if strings.Contains(line, "NON-DEG") &&
				freq.MatchString(fields[0]) {
				state, _ := strconv.Atoi(fields[0])
				if state == 1 {
					res.ZPT, _ = strconv.ParseFloat(fields[1], 64)
					freqs--
				} else if state <= nfreqs+1 {
					f, _ := strconv.ParseFloat(fields[2], 64)
					res.Corr[state-2] = f
					freqs--
				}
			}
		case strings.Contains(line, "NON-DEG(Vt)"):
			if nrot < nfreqs+1 { /* include 0th */
				if nfreqs > 10 { /* two lines of NON-DEG(Vt) if > 10 */
					skip += 7
				} else {
					skip += 3
				}
				rot = true
				nrot++
			}
		case rot:
			// order is A0 -> An
			// in cm-1
			rot = false
			fields := strings.Fields(line)
			tmp := make([]float64, 0, 3)
			for f := range fields {
				v, err := strconv.ParseFloat(fields[f], 64)
				if err != nil {
					v = math.NaN()
				}
				tmp = append(tmp, v)
			}
			res.Rots = append(res.Rots, tmp)
		case delta.MatchString(line):
			// order is DELTA J, K, JK, delta J, K
			// in MHz
			fields := strings.Fields(line)
			f, _ := strconv.ParseFloat(fields[len(fields)-1],
				64)
			res.Deltas = append(res.Deltas, f)
		case phi.MatchString(line):
			// order is PHI J, K, JK, KJ, phi j, jk, k
			// in Hz
			// may need this in delta too
			line := strings.ReplaceAll(line, "D", "E")
			fields := strings.Fields(line)
			f, _ := strconv.ParseFloat(fields[len(fields)-1],
				64)
			res.Phis = append(res.Phis, f)
		case strings.Contains(line, "INT COORD TYPE") &&
			!geom && res.Requil == nil:
			geom = true
			skip++
		case geom && !strings.Contains(line, "LINEAR"):
			if line == "" {
				geom = false
				continue
			}
			fields := strings.Fields(line)
			e, _ := strconv.ParseFloat(fields[2], 64)
			a, _ := strconv.ParseFloat(fields[4], 64)
			res.Requil = append(res.Requil, e)
			res.Ralpha = append(res.Ralpha, a)
		case icn.MatchString(line):
			// Torsions do not appear in r(equil|alpha) part so
			// neglect here as well
			match := atom.FindAllStringSubmatch(line, -1)
			fmt.Fprintf(&buf, "%s(", gparams[len(match)])
			for l, p := range match {
				fmt.Fprintf(&buf, "%s%s",
					strings.TrimSpace(p[2]),
					strings.TrimSpace(p[1]))
				if l < len(match)-1 {
					fmt.Fprint(&buf, "-")
				}
			}
			fmt.Fprint(&buf, ")")
			res.Rhead = append(res.Rhead, buf.String())
			buf.Reset()
		case strings.Contains(line, "FERMI RESONANCE   "):
			skip += 3
			if strings.Contains(line, "TYPE 1") {
				fermi1 = true
			} else {
				fermi2 = true
			}
		case (fermi1 || fermi2) && line == "":
			// just set them both instead of checking
			fermi1 = false
			fermi2 = false
		case fermi1:
			fields := strings.Fields(line)
			key := fields[3]
			fermiMap[key] = append(fermiMap[key],
				fmt.Sprintf("2v_%s", fields[1]))
		case fermi2:
			fields := strings.Fields(line)
			key := fields[3]
			fermiMap[key] = append(fermiMap[key],
				fmt.Sprintf("v_%s+v_%s", fields[1],
					fields[2]))
		}
		// TODO option for BZA and/or BZS
		// TODO option for D in addition to DELTA
	}
	sorter := make([]string, 0, len(fermiMap))
	for k := range fermiMap {
		sorter = append(sorter, k)
	}
	sort.Strings(sorter)
	for _, k := range sorter {
		for _, r := range fermiMap[k] {
			fmt.Fprintf(&buf, "%s=", r)
		}
		fmt.Fprintf(&buf, "v_%s", k)
		res.Fermi = append(res.Fermi, buf.String())
		buf.Reset()
	}
	return res
}
