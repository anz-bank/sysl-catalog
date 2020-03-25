package templategeneration

const ProjectMarkdownTemplate = `
# {{.Title}}

| Package |
| - | {{range $Package := .AlphabeticalRows}}
[{{$Package.PackageName}}]({{$Package.PackageName}}/{{$Package.OutputFile}})|{{end}}

Integration diagram:

![alt text]({{.RootLevelIntegrationDiagram.AppName}}.svg)
`

const PackageMarkdownTemplate = `
[Back](../README.md)
# Package {{.PackageName}}

## Index
| AppName | Endpoint |
| - | - | {{range $Diagram := .SequenceDiagrams}}
| {{$Diagram.AppName}} | [{{$Diagram.EndpointName}}](#{{$Diagram.AppName}}{{$Diagram.EndpointName}}) |{{end}}]

{{range $Diagram := .SequenceDiagrams}}
## {{$Diagram.AppName}} {{$Diagram.EndpointName}}

### Parameter types
{{range $DataModelDiagram := $Diagram.InputDataModel}}
![alt text]({{$DataModelDiagram.OutputFileName__}}.svg)
{{end}}
### Sequence Diagram
![alt text]({{.OutputFileName__}}.svg)
{{end}}

`

const EmbededSvgTemplate = `
[Back](README.md)

![alt text]({{.OutputFileName__}}.svg)

`

// [Back](../README.md)
// # Package {{.PackageName}}

// ## Data Model Diagrams

// ![alt text]({{.PackageName}}_datamodel.svg)

// ## Sequence Diagrams
// | AppName | Endpoint |
// | - | - | {{range $Diagram := .SequenceDiagrams}}
// | {{$Diagram.AppName}} | [{{$Diagram.OutputFileName__}}]({{$Diagram.OutputFileName__}}.md) |{{end}}
