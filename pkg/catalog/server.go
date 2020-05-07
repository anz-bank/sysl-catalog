// server.go: implements http handler interface so that Generator struct can be used directly as a handler
package catalog

import (
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
func (p *Generator) ServerSettings(server, liveReload, imageTags bool) *Generator {
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
		file = strings.ReplaceAll(file, "<body>", `<body><script src="/livereload.js?port=6900&mindelay=10&v=2" data-no-instant defer></script>`)
	}
	w.Write([]byte(file))
}
