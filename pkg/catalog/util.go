package catalog

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/afero"

	"github.com/anz-bank/sysl/pkg/sysl"
)

func sanitiseOutputName(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "/", "")
}

// AlphabeticalEndpoints sorts a map of Applications alphabetically
func AlphabeticalApps(m map[string]*sysl.Application) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// AlphabeticalEndpoints sorts a map of Applications alphabetically
func AlphabeticalTypes(m map[string]*sysl.Type) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// AlphabeticalEndpoints sorts a map of endpoints alphabetically
func AlphabeticalEndpoints(m map[string]*sysl.Endpoint) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// AlphabeticalEndpoints sorts a map of Package alphabetically
func AlphabeticalPackage(m map[string]*Package) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// GenerateMarkdown generates markdown for any object.
func GenerateMarkdown(outputdir, fileName string, object interface{}, t *template.Template, fs afero.Fs) error {
	var buf bytes.Buffer
	if err := t.Execute(&buf, object); err != nil {
		return err
	}
	fs.MkdirAll(outputdir, os.ModePerm)
	return afero.WriteFile(fs, filepath.Join(outputdir, fileName), buf.Bytes(), os.ModePerm)
}

// LoadMarkdownTemplates loads string markdown templates into slices of template objects.
func LoadMarkdownTemplates(markdowns ...string) ([]*template.Template, error) {
	templates := make([]*template.Template, 0, len(markdowns))
	for _, markdownString := range markdowns {
		newTemplate, err := template.New("").Parse(markdownString)
		if err != nil {
			return nil, err
		}
		templates = append(templates, newTemplate)
	}
	return templates, nil
}

// GetAppPackageName returns the package and app name of any sysl application
func GetAppPackageName(app *sysl.Application) (string, string) {
	appName := strings.Join(app.Name.GetPart(), "")
	packageName := appName
	if attr := app.GetAttrs()["package"]; attr != nil {
		packageName = attr.GetS()
	}
	return packageName, appName
}
func NewTypeRef(appName, typeName string) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_TypeRef{
			TypeRef: &sysl.ScopedRef{
				Ref: &sysl.Scope{Appname: &sysl.AppName{
					Part: []string{appName},
				},
					Path: []string{appName, typeName},
				},
			},
		},
	}
}
