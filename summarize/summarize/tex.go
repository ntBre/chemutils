package main

type Table struct {
	Caption   string
	Alignment string
	Header    string
	Body      string
}

func NewTable(caption string) *Table {
	return &Table{Caption: caption}
}

const (
	tabTemplate = `\begin{table}[ht]
\centering
\caption{{{.Caption}}}
\begin{tabular}{{{.Alignment}}}
{{.Header}}
\hline
{{.Body}}
\end{tabular}
\end{table}
`
)
