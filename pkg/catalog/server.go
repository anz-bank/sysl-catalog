// server.go: implements http handler interface so that Generator struct can be used directly as a handler
package catalog

import (
	"encoding/base64"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslwrapper"
	"github.com/spf13/afero"
)

// Update loads another Sysl module into a project and runs
func (p *Generator) Update(m *sysl.Module, errs ...error) *Generator {
	//p.Fs = afero.NewMemMapFs()
	p.errs = []error{}
	for _, err := range errs {
		if err != nil {
			p.errs = append(p.errs, err)
			// Clear generated files since we only want to display an error
			p.GeneratedFiles = nil
			p.Fs = afero.NewMemMapFs()
		}
	}

	if len(p.errs) == 0 {
		p.RootModule = m
		p.GeneratedFiles = make(map[string][]byte)
		p.Mapper = syslwrapper.MakeAppMapper(m)
		p.Mapper.IndexTypes()
		p.Mapper.ConvertTypes()
		if p.RootModule != nil && len(p.ModuleAsMacroPackage(p.RootModule)) <= 1 && !p.CustomTemplate {
			p.StartTemplateIndex = 1 // skip the MacroPackageProject
		} else {
			p.StartTemplateIndex = 0
		}
		p.Run()
	}

	return p
}

// ServerSettings sets the server settings, this should be set before using as http handler
func (p *Generator) ServerSettings(disableCSS, liveReload, imageTags bool) *Generator {
	p.DisableCss = disableCSS
	p.LiveReload = liveReload
	p.OutputDir = "/"
	p.Server = true
	p.Fs = afero.NewMemMapFs()
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
	defer func() {
		if len(p.errs) > 0 {
			bytes = convertToEscapedHTML(fmt.Sprintln(p.errs))
		}
	}()
	request := r.URL.Path
	if p.RootModule == nil && path.Ext(request) != ".ico" {
		bytes = convertToHTML(`<img class="blink-image" src="favicon.ico">` + flashing)
		return
	}
	if p.Fs == nil && path.Ext(request) != ".ico" {
		p.Update(p.RootModule)
	}
	switch path.Ext(request) {
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
		unescapedPath, err := url.PathUnescape(request)
		if err != nil {
			p.errs = append(p.errs, err)
		}
		if svg, ok := p.GeneratedFiles[path.Join(unescapedPath)]; ok {
			bytes = svg
			return
		}
		p.errs = []error{}
		if err != nil {
			p.errs = append(p.errs, err)
		}
		p.GeneratedFiles[path.Join(unescapedPath)] = bytes
		return
	case ".ico":
		bytes, err = base64.StdEncoding.DecodeString(favicon)
		if err != nil {
			p.errs = append(p.errs, err)
			p.Log.Info(err)
		}
		return
	case ".yaml", ".yml", ".json", ".sysl":
		bytes, err = afero.ReadFile(afero.NewOsFs(), strings.TrimPrefix(request, "/"))
		if err != nil {
			p.Log.Error(err)
			return
		}
		return
	case "":
		request += "index.html"
	}
	bytes, err = afero.ReadFile(p.Fs, path.Join(p.OutputDir, request))
	if err != nil {
		if len(p.errs) > 0 && p.errs[len(p.errs)-1].Error() == err.Error() {
			return
		}
		p.errs = append(p.errs, err)
		p.Log.Info(err)
		return
	}
	file = string(bytes)
	if !p.LiveReload {
		p.errs = []error{}
		return
	}
	switch p.Format {
	case "html":
		bytes = []byte(file + liveReload)
		if p.DisableCss {
			bytes = convertToEscapedHTML(file)
		}
	default:
		bytes = convertToEscapedHTML(file)
	}
	p.errs = []error{}
}

func convertToEscapedHTML(file string) []byte {
	return []byte(
		header +
			`<pre style="word-wrap: break-word; white-space: pre-wrap;">` +
			html.EscapeString(file) +
			`</pre>` + liveReload + endTags)
}

func convertToHTML(file string) []byte {
	return []byte(header + file + liveReload + endTags)
}
