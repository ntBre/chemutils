// Package summarize is for summarizing output from quantum chemistry
// programs
package summarize

import (
	"bufio"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// FreqReport gathers harmonic, anharmonic, and resonance-corrected
// frequencies from a spectro  output file for reporting
func Spectro(filename string, nfreqs int) (zpt float64,
	harm, fund, corr []float64,
	rotABC [][]float64,
	deltas, phis []float64) {

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
		harmFund bool
		rot      bool
	)
	freq := regexp.MustCompile(`[0-9]+`)
	delta := regexp.MustCompile(`(?i)delta (J|(JK)|K)`)
	phi := regexp.MustCompile(`(?i)phi (J|(JK)|(KJ)|K)`)
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
				harm = append(harm, h)
				fund = append(fund, f)
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
			if freq.MatchString(fields[0]) {
				state, _ := strconv.Atoi(fields[0])
				if state == 1 {
					zpt, _ = strconv.ParseFloat(fields[1], 64)
					freqs--
				} else if state <= nfreqs+1 {
					f, _ := strconv.ParseFloat(fields[2], 64)
					corr = append(corr, f)
					freqs--
				}
			}
		case strings.Contains(line, "NON-DEG(Vt)"):
			skip += 3
			rot = true
		case rot:
			// order is A0 -> An
			// in cm-1
			// could skip 3 more here to get BZS too
			rot = false
			fields := strings.Fields(line)
			tmp := make([]float64, 0, 3)
			for f := range fields {
				v, _ := strconv.ParseFloat(fields[f], 64)
				tmp = append(tmp, v)
			}
			rotABC = append(rotABC, tmp)
		case delta.MatchString(line):
			// order is DELTA J, K, JK, delta J, K
			// in MHz
			fields := strings.Fields(line)
			f, _ := strconv.ParseFloat(fields[len(fields)-1], 64)
			deltas = append(deltas, f)
		case phi.MatchString(line):
			// order is PHI J, K, JK, KJ, phi j, jk, k
			// in Hz
			// may need this in delta too
			line := strings.ReplaceAll(line, "D", "E")
			fields := strings.Fields(line)
			f, _ := strconv.ParseFloat(fields[len(fields)-1], 64)
			phis = append(phis, f)
		}
		// TODO geometry parameters
		// presumably vibrationally averaged coordinates
		// - pretty sure R(EQUIL), but what are R(G) and R(ALPHA)?
	}
	// put freqs in decreasing order
	sort.Sort(sort.Reverse(sort.Float64Slice(harm)))
	sort.Sort(sort.Reverse(sort.Float64Slice(fund)))
	sort.Sort(sort.Reverse(sort.Float64Slice(corr)))
	return
}
