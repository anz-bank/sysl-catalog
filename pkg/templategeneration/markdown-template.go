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
| Service Name | Method |
| - | - | {{range $Diagram := .SequenceDiagrams}}
| {{$Diagram.AppName}} | [{{$Diagram.EndpointName}}](#{{$Diagram.AppName}}-{{$Diagram.EndpointName}}) |{{end}}]



---
{{range $Diagram := .SequenceDiagrams}}
## {{$Diagram.AppName}} {{$Diagram.EndpointName}}

### Sequence Diagram
![alt text]({{.OutputFileName__}}.svg)

### Request types
{{range $DataModelDiagram := $Diagram.InputDataModel}}
![alt text]({{$DataModelDiagram.OutputFileName__}}.svg)
{{end}}

### Response types
{{range $DataModelDiagram := $Diagram.OutputDataModel}}
![alt text]({{$DataModelDiagram.OutputFileName__}}.svg)
{{end}}
---
{{end}}

`
