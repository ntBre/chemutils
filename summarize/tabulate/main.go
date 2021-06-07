/*
1. find list of files
   - sort and order as you want the columns to be
2. run summarize on each and capture result
   - accept list of files from stdin so you can pipe in
3. put summarize output into nice tables

Example table:

|     | filename1 | filename2 | ... | filenameN |
|-----+-----------+-----------+-----+-----------|
| w_1 |           |           |     |           |
| w_2 |           |           |     |           |
| w_3 |           |           |     |           |
|-----+-----------+-----------+-----+-----------|
| v_1 |           |           |     |           |
| v_2 |           |           |     |           |
| v_3 |           |           |     |           |
*/
package main

import (
	"fmt"
	"os"

	"bufio"

	"github.com/ntBre/chemutils/summarize"
)

type Res struct {
	Name string
	Val  *summarize.Result
}

// tabulate -cols "v1 v2 v3 v4 v5"

/*
   input file like this?
 V1  /home/brent/chem/c2h4/sic/010/v1/freqs/spectro2.out
 V2  /home/brent/chem/c2h4/sic/010/v2/freqs/spectro2.out
 V3  /home/brent/chem/c2h4/sic/010/v3/freqs/spectro2.out
 V4  /home/brent/chem/c2h4/sic/010/v4/freqs/spectro2.out
 V5  /home/brent/chem/c2h4/sic/010/v5/freqs/spectro2.out

*/

// assume we have this
var (
	cols = []string{"V1", "V2", "V3", "V4", "V5"}
)

func main() {
	// TODO cols flag for custom col labels
	// - also read from input?
	// - just take input file instead of stdin?
	// - either
	scanner := bufio.NewScanner(os.Stdin)
	results := make([]Res, 0)
	var width int
	for scanner.Scan() {
		filename := scanner.Text()
		if l := len(filename); l > width {
			width = l
		}
		results = append(results,
			Res{
				Name: filename,
				Val:  summarize.Spectro(filename),
			})
	}
	if cols != nil {
		width = 0
		for _, c := range cols {
			if l := len(c); l > width {
				width = l
			}
		}
	}
	if width < 8 {
		width = 8
	}
	width += 2
	floatfmt := fmt.Sprintf("%%%d.1f", width)
	strngfmt := fmt.Sprintf("%%%ds", width)
	fmt.Printf("%8s", "Mode")
	if cols == nil {
		for _, r := range results {
			fmt.Printf(strngfmt, r.Name)
		}
	} else {
		for c := range results {
			fmt.Printf(strngfmt, cols[c])
		}
	}
	fmt.Print("\n")
	for i, _ := range results[0].Val.Harm {
		fmt.Printf("%8s", fmt.Sprintf("w_{%d}", i+1))
		for _, r := range results {
			fmt.Printf(floatfmt, r.Val.Harm[i])
		}
		fmt.Print("\n")
	}
}
