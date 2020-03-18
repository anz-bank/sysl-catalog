package templategeneration

const IndexMarkdownTemplate = `
# {{.Title}}
| Package |
| - | {{range $Package := .AlphabeticalRows}}
{{$Package.PackageName}}{{end}}
`

const AppMarkdownTemplate = `
[Back](../README.md)
| AppName | Endpoint OutputFileName |
| - | - | {{range $Diagram := .SequenceDiagrams}}
{{$Diagram.OutputFileName}} | {{$Diagram.OutputFileName}} {{end}}
`

const EmbededSvgTemplate = `
[Back](README.md)

![alt text]({{.OutputFileName}}.svg)

`
