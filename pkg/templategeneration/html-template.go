package templategeneration

const ProjectHTMLTemplate = `
<h1 id="-title-">{{.Title}}</h1>
<table>
  <tr>
    <th> Packages </th>
  </tr>
{{range $Package := .AlphabeticalRows}}<tr><td>
<a href="{{$Package.PackageName}}/{{$Package.OutputFile}}">{{$Package.PackageName}}</a></td></tr>{{end}}</p>
</table>
<h2>Integration diagram:</h2>
<p><img src="{{.RootLevelIntegrationDiagram.OutputFileName__}}" alt="alt text"></p>
<h2>Integration diagram with end point analysis:</h2>
<p><img src="{{.RootLevelIntegrationDiagramEPA.OutputFileName__}}" alt="alt text"></p>
`

const PackageHTMLTemplate = `
<p><a href="../index.html">Back</a></p>
<h1 id="package-packagename-">Package {{.PackageName}}</h1>
<h2 id="index">Index</h2>
<p>| Service Name | Method |
| - | - | {{range $Diagram := .SequenceDiagrams}}
| {{$Diagram.AppName}} | <a href="#{{$Diagram.AppName}}-{{$Diagram.EndpointName}}">{{$Diagram.EndpointName}}</a> |{{end}}]</p>
<hr>
<p>{{range $Diagram := .SequenceDiagrams}}</p>
<h2 id="-diagram-appname-diagram-endpointname-">{{$Diagram.AppName}} {{$Diagram.EndpointName}}</h2>
<h3 id="sequence-diagram">Sequence Diagram</h3>
<p><img src="{{.OutputFileName__}}.svg" alt="alt text"></p>
<h3 id="request-types">Request types</h3>
<p>{{range $DataModelDiagram := $Diagram.InputDataModel}}
<img src="{{$DataModelDiagram.OutputFileName__}}.svg" alt="alt text">
{{end}}</p>
<h3 id="response-types">Response types</h3>
<p>{{range $DataModelDiagram := $Diagram.OutputDataModel}}
<img src="{{$DataModelDiagram.OutputFileName__}}.svg" alt="alt text"></p>
<h2 id="-end-">{{end}}</h2>
<p>{{end}}</p>


`
