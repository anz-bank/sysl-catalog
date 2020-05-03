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
{{if ne $Diagram.AppComment ""}}
- {{$Diagram.AppComment}}
{{end}}
{{end}}
{{$first = false}}


## {{$Diagram.AppName}} {{$Diagram.EndpointName}}

{{$Diagram.EndpointComment}}

### Sequence Diagram
{{.Img}}

### Request types
{{range $DataModelDiagram := $Diagram.InputDataModel}}
{{if ne $DataModelDiagram.TypeComment ""}}
- {{$DataModelDiagram.TypeComment}}
{{end}}
{{$DataModelDiagram.Img}}
{{end}}

### Response types
{{range $DataModelDiagram := $Diagram.OutputDataModel}}
{{if ne $DataModelDiagram.TypeComment ""}}
- {{$DataModelDiagram.TypeComment}}
{{end}}
{{$DataModelDiagram.Img}}
{{end}}
{{end}}
---
{{end}}

{{range $appName, $Diagrams := .DatabaseModel}}
## Database {{$appName}}
{{if ne $Diagrams.AppComment ""}}
- {{$Diagrams.AppComment}}
{{end}}
{{$Diagrams.Img}}
{{end}}

## Types
<table>
<tr>
<th>App Name</th>
<th>Diagram</th>
<th>Comment</th>
<th>Full Diagram</th>
</tr>

{{range $typeName, $Diagrams := .Types}}
<tr>
<td>

{{$Diagrams.Simple.AppName}}.<br>{{$typeName}}
</td>
<td>

{{$Diagrams.Simple.Img}}
</td>
<td> 

{{if ne $Diagrams.Simple.TypeComment ""}}<details closed><summary>Comment</summary><br>{{$Diagrams.Simple.TypeComment}}</details>{{end}} 
</td>
<td>

{{$Diagrams.Full.Link}}
</td>
</tr>{{end}}
</table>
`
