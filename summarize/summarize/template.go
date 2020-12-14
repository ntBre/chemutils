package main

type Table struct {
	Caption   string
	Alignment string
	Header    string
	Body      string
}

const (
	texTemplate = `\begin{table}[ht]
\centering
\caption{<.Caption>}
\begin{tabular}{<.Alignment>}
<if .Header ->
<.Header>
\hline
<end ->
<.Body>
\end{tabular}
\end{table}

`

	plainTemplate = `{{.Caption}}:
{{if .Header -}}
{{.Header}}
{{end -}}
{{.Body}}`
)
