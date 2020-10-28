package catalog

import (
	"bytes"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var (
	ofTypeSymbol = regexp.MustCompile(`(?m)(?:<:)(?:.*)`)
)

// CreateMarkdown is a wrapper function that also converts output markdown to html if in server mode
func (p *Generator) CreateMarkdown(t *template.Template, outputFileName string, i interface{}) error {
	var buf bytes.Buffer
	if err := t.Execute(&buf, i); err != nil {
		return err
	}
	if err := p.Fs.MkdirAll(path.Dir(outputFileName), os.ModePerm); err != nil {
		return err
	}
	f2, err := p.Fs.Create(outputFileName)
	if err != nil {
		return err
	}
	out := buf.Bytes()
	var converted bytes.Buffer
	if p.Format == "html" && !p.DisableCss {
		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		)
		if err := md.Convert(out, &converted); err != nil {
			return errors.Wrap(err, "Error converting markdown to html:")
		}
		raw := converted.String()
		raw = strings.ReplaceAll(raw, "README.md", p.OutputFileName)
		out = []byte(header + raw + style + endTags)
	}
	if _, err = f2.Write(out); err != nil {
		return err
	}
	return nil
}
