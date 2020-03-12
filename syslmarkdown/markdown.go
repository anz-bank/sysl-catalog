package syslmarkdown

const Index = `| Service | EndpointName |
| - |:-:|{{range $key, $val := .Apps}}{{range $endpointName, $endpoint := $val.Endpoints}}
|{{$key}}|{{$endpointName}}{{end}}{{end}}
`
