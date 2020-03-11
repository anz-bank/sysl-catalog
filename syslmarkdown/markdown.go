package syslmarkdown

const Index = `| Service | Method |
| - |:-:|{{range $key, $val := .Apps}}{{range $endpointName, $endpoint := $val.Endpoints}}
|{{$key}}|{{$endpointName}}|{{end}}{{end}}
`
