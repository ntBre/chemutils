package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Table struct {
	Caption   string
	Alignment string
	Header    string
	Body      string
}

func (t *Table) Texify(part string) string {
	// TODO hacky - should just pass t.Header or t.Body, or any
	// string for that matter, but can't figure out how in the
	// template
	var (
		ranger []string
		str    strings.Builder
	)
	subscript := regexp.MustCompile(`v(_[0-9]+)`)
	switch part {
	case "BODY":
		ranger = strings.Split(t.Body, "\n")
	case "HEAD":
		ranger = strings.Split(t.Header, "\n")
	}
	for i, line := range ranger {
		fields := strings.Fields(line)
		fmt.Fprint(&str, strings.Join(fields, " & "), `\\`)
		if i < len(ranger)-1 {
			str.WriteString("\n")
		}
	}
	return strings.ReplaceAll(
		subscript.ReplaceAllString(str.String(), "$$\\nu$1$$"),
		"<", "$\\angle$")
}

const (
	texTemplate = `\begin{table}[ht]
\centering
\caption{<.Caption>}
\begin{tabular}{<.Alignment>}
<if .Header ->
<.Texify "HEAD">
\hline
<end ->
<.Texify "BODY">
\end{tabular}
\end{table}

`

	plainTemplate = `{{.Caption}}:
{{if .Header -}}
{{.Header}}
{{end -}}
{{.Body}}
`
)
