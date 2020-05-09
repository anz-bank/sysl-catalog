// server.go: implements http handler interface so that Generator struct can be used directly as a handler
package catalog

import (
	"html"
	"net/http"
	"path"

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
	var (
		bytes []byte
		file  string
		err   error
	)
	defer func() {
		if _, err := w.Write(bytes); err != nil {
			p.Log.Info(err)
		}
	}()
	if p.Fs == nil {
		p.Update(p.Module)
	}
	request := r.RequestURI
	switch path.Ext(request) {
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
		bytes, err = afero.ReadFile(p.Fs, request)
		if err != nil {
			p.Log.Info(err)
		}
		return
	case "":
		request += "index.html"
	}
	bytes, err = afero.ReadFile(p.Fs, request)
	if err != nil {
		p.Log.Info(err)
	}
	file = string(bytes)
	if !p.LiveReload {
		return
	}
	switch p.Format {
	case "html":
		bytes = []byte(file + script)
	default:
		bytes = []byte(
			header +
				`<pre style="word-wrap: break-word; white-space: pre-wrap;">` +
				html.EscapeString(file) +
				`</pre>` + script + endTags)
	}
}
