package catalog

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

## Service Index
| Service Name | Method |
| - | - | {{range $appName, $Diagrams := .SequenceDiagrams}}{{range $Diagram := $Diagrams}}
| {{$appName}} | [{{$Diagram.EndpointNameWithoutSpaces}}](#{{$Diagram.AppName}}-{{$Diagram.EndpointNameWithoutSpaces}}) |{{end}}{{end}}


## Database Index
| Database Name |
| - |
{{range $appName, $Diagrams := .DatabaseModel}}| [{{$appName}}](#Database-{{$appName}}) |{{end}}

## Integration diagram

![alt text]({{.Integration.OutputFileName__}})
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
![alt text]({{.OutputFileName__}})

### Request types
{{range $DataModelDiagram := $Diagram.InputDataModel}}
{{$DataModelDiagram.TypeComment}}
![alt text]({{$DataModelDiagram.OutputFileName__}})
{{end}}

### Response types
{{range $DataModelDiagram := $Diagram.OutputDataModel}}
{{$DataModelDiagram.TypeComment}}
![alt text]({{$DataModelDiagram.OutputFileName__}})
{{end}}
{{end}}
---
{{end}}

{{range $appName, $Diagrams := .DatabaseModel}}
## Database {{$appName}}
{{$Diagrams.AppComment}}
![alt text]({{$Diagrams.OutputFileName__}})
{{end}}

{{range $appName, $Diagrams := .Types}}
## Database {{$appName}}
{{$Diagrams.AppComment}}
![alt text]({{$Diagrams.OutputFileName__}})
{{end}}

`
