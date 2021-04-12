package catalog

import (
	"path"
	"strings"

	"github.com/joshcarp/gop/gop"

	"github.com/anz-bank/pkg/mod"
	"github.com/anz-bank/sysl/pkg/sysl"
)

func IsOpenAPIFile(filePath string) bool {
	name, _ := mod.ExtractVersion(filePath)
	fileExt := path.Ext(name)
	if fileExt == ".yaml" || fileExt == ".yml" || fileExt == ".json" {
		return true
	}
	return false
}

// GetImportPathAndVersion takes a Sysl Application and returns an import and version string
// The import path is of format github.com/org/repo/myfile.yaml
// it CAN contain version tags e.g github.com/org/repo/myfile.yaml@develop or github.com/org/repo/myfile.yaml@v1.1.0
// The version string is of format version-12digitCommitSHA e.g v0.0.0-c63b9e92813a
func GetImportPathAndVersion(retr gop.Retriever, app *sysl.Application) (importPath string, version string, err error) {
	specURL, ok := app.Attrs["redoc-spec"]
	if ok {
		if err != nil {
			return "", "", err
		}
		_, _, err = retr.Retrieve(specURL.GetS())
		if err != nil {
			return "", "", err
		}
		_, ver := mod.ExtractVersion(specURL.GetS())

		importPath = specURL.GetS()
		version = ver
	} else {
		for _, ctx := range app.SourceContexts {
			if IsOpenAPIFile(ctx.GetFile()) {
				importPath = ctx.GetFile()
				version = ctx.GetVersion()
				break
			}
		}
	}
	return importPath, version, nil
}

// BuildSpecURL takes a filepath and version and builds a URL to the cached spec file
// It also trims . prefixes and adds a / so that the URL is relative
func BuildSpecURL(filePath string, version string) string {
	repoPath := strings.TrimPrefix(filePath, ".")
	if !strings.HasPrefix(repoPath, "/") {
		repoPath = "/" + repoPath
	}
	return repoPath
}
