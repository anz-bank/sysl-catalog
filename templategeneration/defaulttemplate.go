package templategeneration

const IndexMarkdownTemplate = `
# {{.Title}}
| Package |
| - | {{range $Package := .AlphabeticalRows}}
[{{$Package.PackageName}}]({{$Package.PackageName}}/{{$Package.OutputFile}})|{{end}}
`

const AppMarkdownTemplate = `
[Back](../README.md)
| AppName | Endpoint OutputFileName__ |
| - | - | {{range $Diagram := .SequenceDiagrams}}
| {{$Diagram.AppName}} | [{{$Diagram.OutputFileName__}}]({{$Diagram.OutputFileName__}}.md) |{{end}}
`

const EmbededSvgTemplate = `
[Back](README.md)

![alt text]({{.OutputFileName__}}.svg)

`
