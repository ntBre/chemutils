package summarize

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
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
	Be     []float64
	Lin    bool
}

// FreqReport gathers harmonic, anharmonic, and resonance-corrected
// frequencies from a spectro  output file for reporting
func Spectro(r io.Reader) *Result {
	res := new(Result)
	fermiMap := make(map[string][]string)
	scanner := bufio.NewScanner(r)
	var (
		line     string
		fields   []string
		skip     int
		corr     bool
		holdFreq float64
		holdZPT  float64
		pos      int
		state    []string
		good     bool
		one      bool
		harmFund bool
		rot      bool
		geom     bool
		fermi1   bool
		fermi2   bool
		buf      bytes.Buffer
		// this is cute
		gparams = []string{"", "", "r", "<"}
	)
	res.Lin = true
	freq := regexp.MustCompile(`[0-9]+`)
	delta := regexp.MustCompile(`(?i)delta (J|(JK)|K)`)
	phi := regexp.MustCompile(`(?i)phi (J|(JK)|(KJ)|K)`)
	icn := regexp.MustCompile(`\([ 0-9]+\)\s+(BOND|ANGLE)`)
	atom := regexp.MustCompile(`([0-9]+)\(([A-Za-z ]+)\)`)
	for scanner.Scan() {
		line = scanner.Text()
		fields = strings.Fields(line)
		switch {
		case skip > 0:
			skip--
		case strings.Contains(line, "BAND CENTER ANALYSIS"):
			skip += 3
			harmFund = true
		case harmFund && len(line) > 1:
			if strings.Contains(line, "DUNHAM") ||
				strings.Contains(line, "VIBRATIONAL ENERGY AND") {
				harmFund = false
			} else if freq.MatchString(fields[0]) && len(fields) > 2 {
				h, _ := strconv.ParseFloat(fields[1], 64)
				f, _ := strconv.ParseFloat(fields[2], 64)
				res.Harm = append(res.Harm, h)
				res.Fund = append(res.Fund, f)
			}
		case strings.Contains(line, "STATE NO.") &&
			!strings.Contains(line, "SPECTRUM"):
			skip += 2
			corr = true
		case corr && strings.Contains(line, "*******************"):
			corr = false
			good = true
		case corr && strings.Contains(line, "NON-DEG (Vs)"):
			// need to grab freq, then check if the next line
			// contains more state info before appending
			holdZPT, _ = strconv.ParseFloat(fields[1], 64)
			holdFreq, _ = strconv.ParseFloat(fields[2], 64)
			pos, _ = strconv.Atoi(fields[0])
			state = nil
			good = true
			one = false
			for _, f := range fields[6:] {
				// if there's a 2 or a second one it's bad
				if f == "2" || (one && f == "1") {
					good = false
					break
				} else if f == "1" {
					one = true
				}
				state = append(state, f)
			}
		case corr && good && strings.Contains(line, "DEGEN   (Vt)"):
			for _, f := range fields[3:] {
				// if there's a 2 or a second one it's bad
				if f == "2" || (one && f == "1") {
					good = false
					break
				} else if f == "1" {
					one = true
				}
				state = append(state, f)
			}
		case corr && good && strings.Contains(line, "DEGEN   (Vl)"):
		case corr && good && len(fields) > 0:
			for _, f := range fields {
				if f == "2" || (one && f == "1") {
					good = false
					break
				} else if f == "1" {
					one = true
				}
				state = append(state, f)
			}
		case corr && len(fields) == 0 && good:
			if !one {
				res.ZPT = holdZPT
			} else if pos < 2 {
				continue
			} else {
				for pos-2 >= len(res.Corr) {
					res.Corr = append(res.Corr, 0)
				}
				res.Corr[pos-2] = holdFreq
			}
		case strings.Contains(line, "NON-DEG(Vt)"):
			for _, f := range fields[2:] {
				if f == "2" || (one && f == "1") {
					good = false
					break
				} else if f == "1" {
					one = true
				}
				state = append(state, f)
			}
		case good && strings.Contains(line, "BZA"):
			rot = true
			res.Lin = false
		case good && res.Lin && strings.Contains(line, "BZS"):
			rot = true
		case rot && good:
			state = nil
			rot = false
			one = false
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
		case fermi1 && !strings.Contains(line, "NOT FOUND"):
			fields := strings.Fields(line)
			key := fields[3]
			fermiMap[key] = append(fermiMap[key],
				fmt.Sprintf("2v_%s", fields[1]))
		case fermi2 && !strings.Contains(line, "NOT FOUND"):
			fields := strings.Fields(line)
			key := fields[3]
			fermiMap[key] = append(fermiMap[key],
				fmt.Sprintf("v_%s+v_%s", fields[1],
					fields[2]))
		case strings.Contains(line, "Be ="):
			fields := strings.Fields(line)
			v, _ := strconv.ParseFloat(fields[2], 64)
			res.Be = append(res.Be, v)
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
