// server.go: implements http handler interface so that Generator struct can be used directly as a handler
package catalog

import (
	"html"
	"net/http"
	"path"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/spf13/afero"
)

// Update loads another Sysl module into a project and runs
func (p *Generator) Update(m *sysl.Module) *Generator {
	p.Fs = afero.NewMemMapFs()
	p.Module = m
	p.Run()
	return p
}

// ServerSettings sets the server settings, this should be set before using as http handler
func (p *Generator) ServerSettings(disableCSS, liveReload, imageTags bool) *Generator {
	p.DisableCss = disableCSS
	p.LiveReload = liveReload
	p.ImageTags = imageTags
	p.OutputDir = "/"
	return p
}

// ServeHTTP is implements the handler interface
func (p *Generator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var bytes []byte
	var file string
	if p.Fs == nil {
		p.Fs = afero.NewMemMapFs()
		p.Run()
	}
	request := r.RequestURI
	switch path.Ext(request) {
	case "":
		request += "index.html"
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
		bytes, _ = afero.ReadFile(p.Fs, request)
		w.Write(bytes)
		return
	}
	bytes, _ = afero.ReadFile(p.Fs, request)
	file = string(bytes)
	if p.LiveReload {
		if strings.Contains(file, `<body>`) {
			// if its html add the script just after the body tag
			file = strings.ReplaceAll(file, "<body>", `<body>`+script)
		} else {
			// if it's raw html, we can render it but still add the livereload script
			file = header +
				`<pre style="word-wrap: break-word; white-space: pre-wrap;">` +
				html.EscapeString(file) +
				`</pre>` + script + endTags
		}
	}
	w.Write([]byte(file))
}
