package catalog

import (
	"bytes"
	"text/template"

	"github.com/sirupsen/logrus"
)

const RedocPage = `<!DOCTYPE html>
<html>
  <head>
    <title>ReDoc</title>
    <!-- needed for adaptive design -->
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">

    <!--
    ReDoc doesn't change outer page styles
    -->
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <redoc spec-url='{{.SpecURL}}'></redoc>
    <script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"> </script>
  </body>
</html>
`

// Redoc is the struct passed to the string template RedocPage
type Redoc struct {
	SpecURL string
}

// BuildRedoc generates a slice of bytes containing the HTML contents of a redoc page
// It takes in the url where the spec is located
func BuildRedoc(specURL string) []byte {
	buf := new(bytes.Buffer)
	t := template.Must(template.New("redoc").Parse(RedocPage))
	err := t.Execute(buf, Redoc{
		SpecURL: specURL,
	})
	if err != nil {
		logrus.Error(err)
	}
	return buf.Bytes()
}
