package templategeneration

const ProjectMarkdownTemplate2 = ``
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
const PackageMarkdownTemplate2 = ``
const PackageMarkdownTemplate = `
[Back](../README.md)
# Package {{.PackageName}}

## Index
| Service Name | Method |
| - | - | {{range $appName, $Diagrams := .SequenceDiagrams}}{{range $Diagram := $Diagrams}}
| {{$appName}} | [{{$Diagram.EndpointNameWithoutSpaces}}](#{{$Diagram.AppName}}-{{$Diagram.EndpointNameWithoutSpaces}}) |{{end}}{{end}}
{{range $appName, $Diagrams := .DatabaseModel}}| {{$appName}} | [Database](#Database-{{$appName}}) |{{end}}

---
{{range $appName, $Diagrams := .SequenceDiagrams}}
{{$first := true}}
{{range $Diagram := $Diagrams}}
{{if $first}}
## {{$Diagram.AppName}}
{{$Diagram.AppComment}}
{{end}}
{{$first = false}}


## {{$Diagram.AppName}} {{$Diagram.EndpointName}}

{{$Diagram.EndpointComment}}

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
{{end}}
---
{{end}}

{{range $appName, $Diagrams := .DatabaseModel}}
## Database {{$appName}}
![alt text]({{$Diagrams.OutputFileName__}}.svg)
{{end}}

`
