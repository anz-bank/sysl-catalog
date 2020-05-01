package catalog

const ProjectMarkdownTemplate = `
# {{.Title}}

| Package |
----|{{range $Package := .AlphabeticalRows}}
[{{$Package.PackageName}}]({{$Package.PackageName}}/{{$Package.OutputFile}})|{{end}}

Integration diagram:

![alt text]({{.RootLevelIntegrationDiagram.Img}})

Integration diagram with end point analysis:

![alt text]({{.RootLevelIntegrationDiagramEPA.Img}})
`

const PackageMarkdownTemplate = `
[Back](../README.md)
# Package {{.PackageName}}

## Service Index
| Service Name | Method |
----|----{{range $appName, $Diagrams := .SequenceDiagrams}}{{range $Diagram := $Diagrams}}
{{$appName}} | [{{$Diagram.EndpointNameWithoutSpaces}}](#{{$Diagram.AppName}}-{{$Diagram.EndpointNameWithoutSpaces}}) |{{end}}{{end}}

## Database Index
| Database Name |
----|
{{range $appName, $Diagrams := .DatabaseModel}}| [{{$appName}}](#Database-{{$appName}}) |{{end}}

[Types](#Types)

## Integration diagram

![alt text]({{.Integration.Img}})

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
![alt text]({{.Img}})

### Request types
{{range $DataModelDiagram := $Diagram.InputDataModel}}
{{$DataModelDiagram.TypeComment}}
![alt text]({{$DataModelDiagram.Img}})
{{end}}

### Response types
{{range $DataModelDiagram := $Diagram.OutputDataModel}}
{{$DataModelDiagram.TypeComment}}
![alt text]({{$DataModelDiagram.Img}})
{{end}}
{{end}}
---
{{end}}

{{range $appName, $Diagrams := .DatabaseModel}}
## Database {{$appName}}
{{$Diagrams.AppComment}}
![alt text]({{$Diagrams.Img}})
{{end}}

## Types
App Name | Diagram | Comment
------------|----------------|------------{{range $typeName, $Diagrams := .Types}}
{{$Diagrams.AppName}}.{{$typeName}} | ![alt text]({{$Diagrams.Img}}) | Comment {{$Diagrams.TypeComment}}|{{end}}

`

//{{$Diagrams.AppName}}.{{$typeName}} | {{$Diagrams.Img}} | Comment {{$Diagrams.TypeComment}}{{end}}
