package templategeneration

const ProjectMarkdownTemplate = `
# {{.Title}}
Integration diagram:

![alt text]({{.RootLevelIntegrationDiagram.AppName}}.svg)

| Package |
| - | {{range $Package := .AlphabeticalRows}}
[{{$Package.PackageName}}]({{$Package.PackageName}}/{{$Package.OutputFile}})|{{end}}
`

const PackageMarkdownTemplate = `
[Back](../README.md)
# Package {{.PackageName}}

## Sequence Diagrams
| AppName | Endpoint OutputFileName__ |
| - | - | {{range $Diagram := .SequenceDiagrams}}
| {{$Diagram.AppName}} | [{{$Diagram.OutputFileName__}}]({{$Diagram.OutputFileName__}}.md) |{{end}}
`

const EmbededSvgTemplate = `
[Back](README.md)

![alt text]({{.OutputFileName__}}.svg)

`