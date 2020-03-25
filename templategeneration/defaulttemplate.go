package templategeneration

const ProjectMarkdownTemplate = `
# {{.Title}}

| Package |
| - | {{range $Package := .AlphabeticalRows}}
[{{$Package.PackageName}}]({{$Package.PackageName}}/{{$Package.OutputFile}})|{{end}}

Integration diagram:

![alt text]({{.RootLevelIntegrationDiagram.OutputFileName__}})

Integration diagram with end point analysis:

![alt text]({{.RootLevelIntegrationDiagramEPA.OutputFileName__}})
`

const PackageMarkdownTemplate = `
[Back](../README.md)
# Package {{.PackageName}}

## Index
| AppName | Endpoint |
| - | - | {{range $Diagram := .SequenceDiagrams}}
| {{$Diagram.AppName}} | [{{$Diagram.EndpointName}}](#{{$Diagram.AppName}}%20{{$Diagram.EndpointName}}) |{{end}}]



---
{{range $Diagram := .SequenceDiagrams}}
## {{$Diagram.AppName}} {{$Diagram.EndpointName}}

### Sequence Diagram
![alt text]({{.OutputFileName__}}.svg)

### Parameter types
{{range $DataModelDiagram := $Diagram.InputDataModel}}
![alt text]({{$DataModelDiagram.OutputFileName__}}.svg)
{{end}}

### Return types
{{range $DataModelDiagram := $Diagram.OutputDataModel}}
![alt text]({{$DataModelDiagram.OutputFileName__}}.svg)
{{end}}
---
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
