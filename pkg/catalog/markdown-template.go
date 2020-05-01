package catalog

const ProjectMarkdownTemplate = `
# {{.Title}}

| Package |
----|{{range $Package := .AlphabeticalRows}}
[{{$Package.PackageName}}]({{$Package.PackageName}}/{{$Package.OutputFile}})|{{end}}

Integration diagram:

{{.RootLevelIntegrationDiagram.Img}}

Integration diagram with end point analysis:

{{.RootLevelIntegrationDiagramEPA.Img}}
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

{{.Integration.Img}}

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
{{.Img}}

### Request types
{{range $DataModelDiagram := $Diagram.InputDataModel}}
{{$DataModelDiagram.TypeComment}}
{{$DataModelDiagram.Img}}
{{end}}

### Response types
{{range $DataModelDiagram := $Diagram.OutputDataModel}}
{{$DataModelDiagram.TypeComment}}
{{$DataModelDiagram.Img}}
{{end}}
{{end}}
---
{{end}}

{{range $appName, $Diagrams := .DatabaseModel}}
## Database {{$appName}}
{{$Diagrams.AppComment}}
{{$Diagrams.Img}}
{{end}}

## Types
<table>
<tr>
<th>App Name</th>
<th>Diagram</th>
<th>Comment</th>
</tr>
<tr>{{range $typeName, $Diagrams := .Types}}
<td>{{$Diagrams.AppName}}.{{$typeName}} </td>
<td> {{$Diagrams.Img}}</td>
<td>  {{if ne $Diagrams.TypeComment ""}}<details closed><summary>Comment</summary><br>{{$Diagrams.TypeComment}}</details>{{end}} </td></tr>{{end}}
</table>
`

//{{$Diagrams.AppName}}.{{$typeName}} |<br> {{$Diagrams.Img}} <br>| Comment {{$Diagrams.TypeComment}}{{end}}
