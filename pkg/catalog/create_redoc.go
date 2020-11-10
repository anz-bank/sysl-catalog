package catalog

import (
	"bytes"
	"os"
	"path"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/ghodss/yaml"
)

const RedocPage = `<!DOCTYPE html>
<html>
  <head>
    <title>ReDoc</title>
    <!-- needed for adaptive design -->
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
<div id='redoc-container'></div>
<script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"></script>
<script>
Redoc.init({{.}}
,{}, document.getElementById('redoc-container'))
</script>
  </body>
</html>
`

// Redoc is the struct passed to the string template RedocPage
type Redoc struct {
	SpecURL string
}

// CreateRedoc registers a file that needs to be created when either:
// - The @redoc-spec attribute has been set
// - The source context has an extension suggesting it is an OpenAPI file
func (p *Generator) CreateRedoc(app *sysl.Application, appName string) string {
	if app == nil || appName == "" {
		return ""
	}
	importPath, _, err := GetImportPathAndVersion(p.Retriever, app)
	if err != nil {
		p.Log.Error(err)
		return ""
	}
	if !IsOpenAPIFile(importPath) {
		return ""
	}
	appName = strings.ReplaceAll(appName, " :: ", "_")
	redocOutputPath, _ := CreateFileName(p.CurrentDir, appName+".redoc.html")
	redocOutputPath = path.Join(p.OutputDir, redocOutputPath)
	var c []byte
	if p.Retriever != nil {
		c, _, _ = p.Retriever.Retrieve(importPath)
	} else {
		return ""
	}
	var buf bytes.Buffer
	js, _ := yaml.YAMLToJSON(c)
	if err := p.Redoc.Execute(&buf, string(js)); err != nil {
		return ""
	}

	link, _ := CreateFileName("", appName+".redoc.html")
	_ = p.Fs.MkdirAll(path.Dir(redocOutputPath), os.ModePerm)
	file, err := p.Fs.Create(redocOutputPath)
	if err != nil {
		p.Log.Error("error creating redoc file: ", err)
		return ""
	}
	_, err = file.Write(buf.Bytes())
	if err != nil {
		p.Log.Error("error writing redoc: ", err)
		return ""
	}
	return link
}
