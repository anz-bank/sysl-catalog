package catalog

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/sysl/pkg/sysl"
)

func TestIsOpenAPIFileJSON(t *testing.T) {
	t.Parallel()
	sourceContext := &sysl.SourceContext{File: "github.com/myorg/test/spec.json"}
	result := IsOpenAPIFile(sourceContext)
	assert.True(t, result)
}

func TestIsOpenAPIFileYAML(t *testing.T) {
	t.Parallel()
	sourceContext := &sysl.SourceContext{File: "github.com/myorg/test/spec.yaml"}
	result := IsOpenAPIFile(sourceContext)
	assert.True(t, result)
}

func TestIsOpenAPIFileYAMLOAS(t *testing.T) {
	t.Parallel()
	sourceContext := &sysl.SourceContext{File: "github.com/myorg/test/spec.oas.yaml"}
	result := IsOpenAPIFile(sourceContext)
	assert.True(t, result)
}

func TestIsOpenAPIFileSysl(t *testing.T) {
	t.Parallel()
	sourceContext := &sysl.SourceContext{File: "github.com/myorg/test/spec.sysl"}
	result := IsOpenAPIFile(sourceContext)
	assert.False(t, result)
}

func TestIsOpenAPIFileEmpty(t *testing.T) {
	t.Parallel()
	sourceContext := &sysl.SourceContext{}
	result := IsOpenAPIFile(sourceContext)
	assert.False(t, result)
}

func TestBuildSpecURL(t *testing.T) {
	type args struct {
		source *sysl.SourceContext
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Simple", args{source: &sysl.SourceContext{File: "./pkg/catalog/test/simple.yaml"}}, "/pkg/catalog/test/simple.yaml", false},
		{"NoDot", args{source: &sysl.SourceContext{File: "/pkg/catalog/test/simple.yaml"}}, "/pkg/catalog/test/simple.yaml", false},
		{"AppendForwardSlash", args{source: &sysl.SourceContext{File: "pkg/catalog/test/simple.yaml"}}, "/pkg/catalog/test/simple.yaml", false},
		{"AppendVersion", args{source: &sysl.SourceContext{File: "github.com/anz-bank/sysl-examples/demos/grocerystore/grocerystore.sysl", Version: "v0.0.0-c63b9e92813a"}}, "/github.com/anz-bank/sysl-examples@v0.0.0-c63b9e92813a/demos/grocerystore/grocerystore.sysl", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildSpecURL(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildSpecURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuildSpecURL() = %v, want %v", got, tt.want)
			}
		})
	}
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
