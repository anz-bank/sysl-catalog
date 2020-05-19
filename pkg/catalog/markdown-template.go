// markdown-template.go: the markdown template used to template the sysl module
package catalog

const ProjectTemplate = `
{{range $name, $link := .Links}}
[{{$name}}]({{$link}})
{{end}}
# {{Base .Title}}

| Package |
----|{{range $val := Packages .Module}}
[{{$val}}]({{$val}}/README.md)|{{end}}

## Integration Diagram
<img src="{{CreateIntegrationDiagram .Module .Title false}}">

## End Point Analysis Integration Diagram
<img src="{{CreateIntegrationDiagram .Module .Title true}}">

`

const MacroPackageProject = `
# {{Base .Title}}

| Package |
----|{{range $val := MacroPackages .Module}}
[{{$val}}]({{$val}}/README.md)|{{end}}

## Integration Diagram
<img src="{{CreateIntegrationDiagram .Module .Title false}}">

## End Point Analysis Integration Diagram
<img src="{{CreateIntegrationDiagram .Module .Title true}}">

`

const NewPackageTemplate = `
[Back](../README.md)
{{$packageName := ModulePackageName .}}

# {{$packageName}}

## Service Index
| Service Name | Method | Source Location |
----|----|----{{$Apps := .Apps}}{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}{{if eq (hasPattern $app.Attrs "ignore") false}}{{$Endpoints := $app.Endpoints}}{{range $endpointName := SortedKeys $Endpoints}}{{$endpoint := index $Endpoints $endpointName}}{{if eq (hasPattern $endpoint.Attrs "ignore") false}}
{{$appName}} | [{{$endpoint.Name}}](#{{$appName}}-{{SanitiseOutputName $endpoint.Name}}) | [{{SourcePath $app}}]({{SourcePath $app}})|  {{end}}{{end}}{{end}}{{end}}

![]({{CreateIntegrationDiagram . $packageName false}})

{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}
{{if eq (hasPattern $app.Attrs "ignore") false}}
{{if ne $appName $packageName}}
# {{$appName}}{{end}}

{{Attribute $app "description"}}
{{range $e := $app.Endpoints}}
{{if eq (hasPattern $e.Attrs "ignore") false}}
## {{$appName}} {{SanitiseOutputName $e.Name}}
{{Attribute $e "description"}}

![]({{CreateSequenceDiagram $appName $e}})

### Request types
{{if eq (len $e.Param) 0}}
No Request types
{{end}}

{{range $param := $e.Param}}
{{Attribute $param.Type "description"}}

![]({{CreateParamDataModel $app $param}})
{{end}}

{{if $e.RestParams}}{{if $e.RestParams.UrlParam}}
{{range $param := $e.RestParams.UrlParam}}
{{$PathDataModel := (CreatePathParamDataModel $appName $param)}}
{{if ne $PathDataModel ""}}
### Path Parameter

![]({{CreatePathParamDataModel $appName $param}})
{{end}}{{end}}{{end}}

{{if $e.RestParams.QueryParam}}
{{range $param := $e.RestParams.QueryParam}}
{{$queryDataModel := (CreateQueryParamDataModel $appName $param)}}
{{if ne $queryDataModel ""}}
### Query Parameter

![]({{$queryDataModel}})
{{end}}{{end}}{{end}}{{end}}

### Response types
{{$responses := false}}
{{range $s := $e.Stmt}}{{$diagram := CreateReturnDataModel  $appName $s $e}}{{if ne $diagram ""}}
{{$responses = true}}
{{$ret := (GetReturnType $e $s)}}{{if $ret }}
{{Attribute $ret "description"}}{{end}}

![]({{$diagram}})
{{end}}{{end}}
{{if eq $responses false}}
No Response Types
{{end}}{{end}}{{end}}{{end}}{{end}}

{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}
{{if hasPattern $app.GetAttrs "db"}}

## Database
{{Attribute $app "description"}}
![]({{GenerateDataModel $app}})
{{end}}{{end}}

### Types

<table>
<tr>
<th>App Name</th>
<th>Diagram</th>
<th>Description</th>
<th>Full Diagram</th>
{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}{{$types := $app.Types}}
{{if ne (hasPattern $app.Attrs "db") true}}
</tr>

{{range $typeName := SortedKeys $types}}{{$type := index $types $typeName}}
<tr>
<td>

{{$appName}}.<br>{{$typeName}}
</td>
<td>

<img src="{{CreateTypeDiagram  $app $typeName $type false}}">
</td>
<td> 

{{if ne (Attribute $type "description") ""}}<details closed><summary>Description</summary><br>{{Attribute $type "description"}}</details>{{end}} 
</td>
<td>

<a href="{{CreateTypeDiagram  $app $typeName $type true}}">Link</a>
</td>
</tr>{{end}}{{end}}{{end}}
</table>

`
