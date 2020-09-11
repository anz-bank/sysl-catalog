package catalog

import (
	"testing"

	"github.com/alecthomas/assert"
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
		{"Handles @develop", "github.com/myorg/test/spec.json@develop", true},
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
		{"AppendVersionWithBranchTag", args{filePath: "github.com/anz-bank/sysl-examples/demos/grocerystore/grocerystore.sysl@develop", version: "v0.0.0-c63b9e92813a"}, "/github.com/anz-bank/sysl-examples@v0.0.0-c63b9e92813a/demos/grocerystore/grocerystore.sysl", false},
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
