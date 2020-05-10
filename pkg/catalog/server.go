// server.go: implements http handler interface so that Generator struct can be used directly as a handler
package catalog

import (
	"encoding/base64"
	"fmt"
	"html"
	"net/http"
	"path"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/spf13/afero"
)

// Update loads another Sysl module into a project and runs
func (p *Generator) Update(m *sysl.Module, errs ...error) *Generator {
	p.Fs = afero.NewMemMapFs()
	p.Module = m
	for _, err := range errs {
		if p.errs == nil {
			p.errs = make([]error, 0, len(errs))
		}
		if err != nil {
			p.errs = append(p.errs, err)
		}
	}
	if len(p.errs) != 0 {
		return p
	}
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
	if len(p.errs) > 0 {
		bytes = convertToEscapedHTML(fmt.Sprintln(p.errs))
		p.errs = []error{}
		return
	}
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
	case ".ico":
		bytes, err = base64.StdEncoding.DecodeString(favicon)
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
		bytes = convertToEscapedHTML(file)
	}
}

func convertToEscapedHTML(file string) []byte {
	return []byte(
		header +
			`<pre style="word-wrap: break-word; white-space: pre-wrap;">` +
			html.EscapeString(file) +
			`</pre>` + script + endTags)
}
