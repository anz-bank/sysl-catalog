//+build integration

package catalog

import (
	"os"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/pkg/mod"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Note the tests in this folder require SYSL_GITHUB_TOKEN to be set

func TestGetImportPathAndVersion(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)
	cacheDir := homeDir + "/.sysl/"
	err = mod.Config("github", mod.GoModulesOptions{}, mod.GitHubOptions{CacheDir: cacheDir, Fs: afero.NewMemMapFs()}) // Setup sysl module in Github mode
	assert.NoError(t, err)
	importPath := "github.com/cuminandpaprika/syslmod/specs/brokenOpenAPI.yaml"
	attrs := map[string]*sysl.Attribute{
		"redoc-spec": {
			Attribute: &sysl.Attribute_S{
				S: importPath,
			},
		},
	}
	app := &sysl.Application{Attrs: attrs}
	result, ver, err := GetImportPathAndVersion(app)
	assert.NoError(t, err)
	assert.Equal(t, importPath, result)
	assert.Equal(t, "v0.0.0-3db1c953643b", ver)
}

func TestGetImportPathAndVersionBranch(t *testing.T) {
	importPath := "github.com/cuminandpaprika/syslmod/specs/brokenOpenAPI.yaml@develop"
	attrs := map[string]*sysl.Attribute{
		"redoc-spec": {
			Attribute: &sysl.Attribute_S{
				S: importPath,
			},
		},
	}
	app := &sysl.Application{Attrs: attrs}
	result, ver, err := GetImportPathAndVersion(app)
	assert.NoError(t, err)
	assert.Equal(t, "github.com/cuminandpaprika/syslmod/specs/brokenOpenAPI.yaml@develop", result)
	assert.Equal(t, "v0.0.0-3db1c953643b", ver)
}

func TestGetImportPathAndVersionNonExistentFile(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)
	cacheDir := homeDir + "/.sysl/"
	err = mod.Config("github", mod.GoModulesOptions{}, mod.GitHubOptions{CacheDir: cacheDir, Fs: afero.NewMemMapFs()}) // Setup sysl module in Github mode
	assert.NoError(t, err)
	importPath := "github.com/cuminandpaprika/syslmod/nonexistent.yaml"
	attrs := map[string]*sysl.Attribute{
		"redoc-spec": {
			Attribute: &sysl.Attribute_S{
				S: importPath,
			},
		},
	}
	app := &sysl.Application{Attrs: attrs}
	_, _, err = GetImportPathAndVersion(app)
	assert.Error(t, err)
}

func TestCreateRedocFromAttribute(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)
	cacheDir := homeDir + "/.sysl/"
	err = mod.Config("github", mod.GoModulesOptions{Root: ""}, mod.GitHubOptions{CacheDir: cacheDir, Fs: afero.NewMemMapFs()}) // Setup sysl module in Github mode
	assert.NoError(t, err)
	appName := "myAppName"
	fileName := "github.com/cuminandpaprika/syslmod/specs/brokenOpenAPI.yaml"
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
