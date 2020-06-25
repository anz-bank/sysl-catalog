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
	t.Parallel()
	expected := "https://raw.githubusercontent.com/anz-bank/sysl-catalog/master/pkg/catalog/test/simple.yaml"
	url := BuildSpecURL(&sysl.SourceContext{File: "pkg/catalog/test/simple.yaml"})
	assert.Equal(t, expected, url)
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
