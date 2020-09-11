package catalog

import (
	"os"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/pkg/mod"
	"github.com/anz-bank/sysl/pkg/sysl"
)

func TestIsOpenAPIFile(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{"Handles empty string", "", false},
		{"Handles .sysl", "github.com/myorg/test/spec.sysl", false},
		{"Handles .oas.yaml", "github.com/myorg/test/spec.oas.yaml", true},
		{"Handles .yaml", "github.com/myorg/test/spec.yaml", true},
		{"Handles .yml", "github.com/myorg/test/spec.yml", true},
		{"Handles .json", "github.com/myorg/test/spec.json", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsOpenAPIFile(tt.arg)
			if got != tt.want {
				t.Errorf("IsOpenAPIFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildSpecURL(t *testing.T) {
	type args struct {
		filePath string
		version  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Simple", args{filePath: "./pkg/catalog/test/simple.yaml"}, "/pkg/catalog/test/simple.yaml", false},
		{"NoDot", args{filePath: "/pkg/catalog/test/simple.yaml"}, "/pkg/catalog/test/simple.yaml", false},
		{"AppendForwardSlash", args{filePath: "pkg/catalog/test/simple.yaml"}, "/pkg/catalog/test/simple.yaml", false},
		{"AppendVersion", args{filePath: "github.com/anz-bank/sysl-examples/demos/grocerystore/grocerystore.sysl", version: "v0.0.0-c63b9e92813a"}, "/github.com/anz-bank/sysl-examples@v0.0.0-c63b9e92813a/demos/grocerystore/grocerystore.sysl", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildSpecURL(tt.args.filePath, tt.args.version)
			if got != tt.want {
				t.Errorf("BuildSpecURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Note this test might require SYSL_GITHUB_TOKEN to be set
func TestGetImportPathAndVersion(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)
	cacheDir := homeDir + "/.sysl/"
	err = mod.Config("github", nil, &cacheDir, nil) // Setup sysl module in Github mode
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
func TestStripExtensionSSH(t *testing.T) {
	t.Parallel()
	url := "https://github.com/anz-bank/sysl-catalog.git"
	expected := "https://github.com/anz-bank/sysl-catalog"
	stripped := StripExtension(url)
	assert.Equal(t, expected, stripped)
}
func TestStripExtensionNone(t *testing.T) {
	t.Parallel()
	url := "https://github.com/anz-bank/sysl-catalog"
	stripped := StripExtension(url)
	assert.Equal(t, url, stripped)
}
func TestBuildGitURLGithubPublic(t *testing.T) {
	t.Parallel()
	repoPath := "https://github.com/anz-bank/sysl-catalog"
	result := BuildGithubRawURL(repoPath)
	assert.Equal(t, "https://raw.githubusercontent.com/anz-bank/sysl-catalog/master/", result)
}

func TestBuildGitURLGithubEnterprise(t *testing.T) {
	t.Parallel()
	repoPath := "https://github.myorg.com/anz-bank/sysl-catalog"
	result := BuildGithubRawURL(repoPath)
	assert.Equal(t, "https://github.myorg.com/raw/anz-bank/sysl-catalog/master/", result)
}

func TestBuildGithubBlobURL(t *testing.T) {
	t.Parallel()
	assert.Equal(t,
		"github.com/user/repo/blob/master",
		BuildGithubBlobURL("github.com/user/repo"),
	)
}
