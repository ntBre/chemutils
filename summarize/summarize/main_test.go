package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"bytes"

	"os"

	"os/exec"

	"github.com/ntBre/chemutils/summarize"
)

func TestColPrint(t *testing.T) {
	datum0 := []float64{1, 2, 3}
	datum1 := []float64{4, 5, 6}
	got := colPrint("%4.1f", true, datum0, datum1)
	want := "    1 1.0 4.0\n    2 2.0 5.0\n    3 3.0 6.0"
	if got != want {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

func TestPrintResult(t *testing.T) {
	tmp := *cm
	*cm = true
	defer func() {
		*cm = tmp
	}()
	initConst()
	res := summarize.SpectroFile("testfiles/spectro.out")
	truefile := "testfiles/summary.txt"
	var got bytes.Buffer
	printAll(&got, res)
	want, _ := ioutil.ReadFile(truefile)
	if !reflect.DeepEqual(got.Bytes(), want) {
		t.Errorf("got:\n%v, wanted:\n%v\n", got.String(), string(want))
		out, _ := os.Create("testfiles/problem.out")
		defer out.Close()
		got.WriteTo(out)
	}
	// (diff "testfiles/problem.out" "testfiles/summary.txt")
}

func TestTex(t *testing.T) {
	// inverting the meaning of short
	if !testing.Short() {
		t.Skip("skipping TestTex")
	}
	dir := "testfiles/tex"
	os.MkdirAll(dir, 0755)
	hold := *tex
	*tex = true
	initConst()
	texname := filepath.Join(dir, "test.tex")
	texfile, _ := os.Create(texname)
	defer func() {
		*tex = hold
		texfile.Close()
	}()
	res := summarize.SpectroFile("testfiles/spectro.out")
	fmt.Fprint(texfile, "\\documentclass{article}\n\\begin{document}\n\n")
	printAll(texfile, res)
	fmt.Fprint(texfile, "\\end{document}\n")
	pdfcmd := exec.Command("pdflatex", "test.tex")
	pdfcmd.Dir = dir
	pdfcmd.Run()
	zathura := exec.Command("zathura", "test.pdf")
	zathura.Dir = dir
	zathura.Run()
}

func TestEqnify(t *testing.T) {
	tests := []struct {
		give string
		want string
		end  bool
	}{
		{
			give: "1	r(C_2 - C_3)",
			want: `S_{1} &= &r(\text{C}_2 - \text{C}_3)\\`,
		},
		{
			give: "2	r(C_1 - C_2) + r(C_1 - C_3)",
			want: `S_{2} &= &\frac{1}{\sqrt{2}}[r(\text{C}_1 - \text{C}_2)` +
				` + r(\text{C}_1 - \text{C}_3)]\\`,
		},
		{
			give: "4	<(H_4 - C_2 - C_1) + <(H_5 - C_3 - C_1)",
			want: `S_{4} &= &\frac{1}{\sqrt{2}}[` +
				`\angle(\text{H}_4 - \text{C}_2 - \text{C}_1)` +
				` + \angle(\text{H}_5 - \text{C}_3 - \text{C}_1)]\\`,
		},
		{
			give: "8	t(H_4 - C_2 - C_1 - C_3) - t(H_5 - C_3 - C_1 - C_2)",
			want: `S_{8} &= &\frac{1}{\sqrt{2}}[` +
				`\tau(\text{H}_4 - \text{C}_2 - \text{C}_1 - \text{C}_3)` +
				` - \tau(\text{H}_5 - \text{C}_3 - \text{C}_1 - ` +
				`\text{C}_2)]\\`,
		},
	}
	for _, test := range tests {
		got := Eqnify(test.give, test.end)
		if got != test.want {
			t.Errorf("got\n%v\nwanted\n%v\n", got, test.want)
		}
	}
}
