package main

const IndexMarkdownTemplate = `
| Package | Service Name | EndpointName |
| - | - | - |{{range $PackageName, $App := .}}{{range $Service := $App.EndPoints}}
[{{$App.PackageName}}]({{$App.PackageRelLink}})|{{$Service.AppName}}|[{{$Service.EndpointName}}]({{$Service.Package}}/{{$Service.EndpointName}}.svg.md) |{{end}}{{end}}
`

const AppMarkdownTemplate = `
[Back](../README.md)
| Service | EndpointName |
| - |:-:|
{{range $EndPoints := .}}{{$EndPoints.AppName}}|[{{$EndPoints.EndpointName}}]({{$EndPoints.EndpointName}}.svg.md) |
{{end}}
`

const embededSvgTemplate = `
[Back](README.md)

![alt text]({{.EndpointName}}.svg)

`
