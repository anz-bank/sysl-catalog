package main

const IndexMarkdownTemplate = `
# {{.Title}}
| Package |
| - | {{range $Package := .AlphabeticalRows}}
{{$Package.PackageName}}{{end}}
`

const AppMarkdownTemplate = `
[Back](../README.md)
| AppName | Endpoint Name |
| - | - | {{range $Diagram := .SequenceDiagrams}}
{{$Diagram.Name}} | {{$Diagram.Name}} {{end}}
`

const embededSvgTemplate = `
[Back](README.md)

![alt text]({{.Name}}.svg)

`
