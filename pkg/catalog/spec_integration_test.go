//+build integration

package catalog

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Note the tests in this folder require SYSL_TOKENS to be set

func TestGetImportPathAndVersion(t *testing.T) {
	t.Parallel()
	importPath := "github.com/cuminandpaprika/syslmod/specs/brokenOpenAPI.yaml@master"
	attrs := map[string]*sysl.Attribute{
		"redoc-spec": {
			Attribute: &sysl.Attribute_S{
				S: importPath,
			},
		},
	}
	app := &sysl.Application{Attrs: attrs}
	result, ver, err := GetImportPathAndVersion(app, afero.NewMemMapFs())
	assert.NoError(t, err)
	assert.Equal(t, importPath, result)
	assert.Equal(t, "master", ver)
}

func TestGetImportPathAndVersionBranch(t *testing.T) {
	t.Parallel()
	importPath := "github.com/cuminandpaprika/syslmod/specs/brokenOpenAPI.yaml@develop"
	attrs := map[string]*sysl.Attribute{
		"redoc-spec": {
			Attribute: &sysl.Attribute_S{
				S: importPath,
			},
		},
	}
	app := &sysl.Application{Attrs: attrs}
	result, ver, err := GetImportPathAndVersion(app, afero.NewMemMapFs())
	assert.NoError(t, err)
	assert.Equal(t, "github.com/cuminandpaprika/syslmod/specs/brokenOpenAPI.yaml@develop", result)
	assert.Equal(t, "develop", ver)
}

func TestGetImportPathAndVersionNonExistentFile(t *testing.T) {
	t.Parallel()
	importPath := "github.com/cuminandpaprika/syslmod/nonexistent.yaml@master"
	attrs := map[string]*sysl.Attribute{
		"redoc-spec": {
			Attribute: &sysl.Attribute_S{
				S: importPath,
			},
		},
	}
	app := &sysl.Application{Attrs: attrs}
	_, _, err := GetImportPathAndVersion(app, afero.NewMemMapFs())
	assert.Error(t, err)
}

func TestCreateRedocFromAttribute(t *testing.T) {
	t.Parallel()
	appName := "myAppName"
	fileName := "github.com/cuminandpaprika/syslmod/specs/brokenOpenAPI.yaml@master"
	attrs := map[string]*sysl.Attribute{
		"redoc-spec": {
			Attribute: &sysl.Attribute_S{
				S: fileName,
			},
		},
	}

	app := &sysl.Application{Attrs: attrs}
	gen := Generator{
		RedocFilesToCreate: make(map[string]string),
		Redoc:              true,
		Log:                logrus.New(),
	}
	link := gen.CreateRedoc(app, appName)
	assert.Equal(t, "myappname.redoc.html", link)
}
