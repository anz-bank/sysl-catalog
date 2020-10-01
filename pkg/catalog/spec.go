package catalog

import (
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/anz-bank/pkg/mod"
	"github.com/anz-bank/sysl/pkg/env"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/joshcarp/gop/gop/cli"
	"github.com/spf13/afero"
)

// CopySyslCacheDir copies ALL the contents of the sysl module cache directory
// (typically found in ~/.sysl/github.com) and outputs it to the specified folder
func CopySyslCacheDir(targetDir string) error {
	syslCacheDir := env.SYSL_CACHE.Value()
	cpSpecsToOutputDir := exec.Command("cp", "-r", syslCacheDir+"/.", targetDir)
	if err := cpSpecsToOutputDir.Run(); err != nil {
		return err
	}
	return nil
}

func IsOpenAPIFile(filePath string) bool {
	name, _ := mod.ExtractVersion(filePath)
	fileExt := path.Ext(name)
	if fileExt == ".yaml" || fileExt == ".yml" || fileExt == ".json" {
		return true
	}
	return false
}

func initRetriever(fs afero.Fs) (*cli.Retriever, error) {
	tokensEnv := env.SYSL_TOKENS.Value() // Expects the token to be in form gita.com:<tokena>,gitb.com:<tokenb>
	var hostTokens []string
	var cache, proxy string
	if tokensEnv != "" {
		hostTokens = strings.Split(tokensEnv, ",")
	}
	tokenmap := make(map[string]string, len(hostTokens))
	for _, e := range hostTokens {
		arr := strings.Split(e, ":")
		if len(arr) < 2 {
			return nil, fmt.Errorf("SYSL_TOKENS env var is invalid, should be in form `gita.com:<tokena>,gitb.com:<tokenb>`")
		}
		tokenmap[arr[0]] = arr[1]
	}
	if moduleFlag := env.SYSL_MODULES.Value(); moduleFlag != "" && moduleFlag != "false" && moduleFlag != "off" {
		cache = env.SYSL_CACHE.Value()
		proxy = env.SYSL_PROXY.Value()
	}
	retriever := cli.Default(fs, cache, proxy, tokenmap)
	return &retriever, nil
}

// GetImportPathAndVersion takes a Sysl Application and returns an import and version string
// The import path is of format github.com/org/repo/myfile.yaml
// it CAN contain version tags e.g github.com/org/repo/myfile.yaml@develop or github.com/org/repo/myfile.yaml@v1.1.0
// The version string is of format version-12digitCommitSHA e.g v0.0.0-c63b9e92813a
func GetImportPathAndVersion(app *sysl.Application, fs afero.Fs) (importPath string, version string, err error) {
	specURL, ok := app.Attrs["redoc-spec"]
	if ok {
		retriever, err := initRetriever(fs)
		if err != nil {
			return "", "", err
		}
		_, _, err = retriever.Retrieve(specURL.GetS())
		if err != nil {
			return "", "", err
		}
		_, ver := mod.ExtractVersion(specURL.GetS())

		importPath = specURL.GetS()
		version = ver
	} else {
		importPath = app.SourceContext.GetFile()
		version = app.SourceContext.GetVersion()
	}
	return importPath, version, nil
}

// BuildSpecURL takes a filepath and version and builds a URL to the cached spec file
// It also trims . prefixes and adds a / so that the URL is relative
func BuildSpecURL(filePath string, version string) string {
	_, ver := mod.ExtractVersion(filePath)
	if ver == "" {
		filePath = filePath + "@master"
	}
	repoPath := strings.TrimPrefix(filePath, ".")
	if !strings.HasPrefix(repoPath, "/") {
		repoPath = "/" + repoPath
	}
	return repoPath
}
