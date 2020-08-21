package catalog

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/anz-bank/sysl/pkg/mod"
	"github.com/anz-bank/sysl/pkg/sysl"
)

// CopySyslModCache copies ALL the contents of the sysl module cache directory
// (typically found in ~/.sysl/github.com) and outputs it to the specified folder
func CopySyslModCache(targetDir string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	syslCacheDir := homeDir + "/.sysl/github.com"
	cpSpecsToOutputDir := exec.Command("cp", "-r", syslCacheDir, targetDir)
	if err := cpSpecsToOutputDir.Run(); err != nil {
		return err
	}
	return nil
}

func IsOpenAPIFile(filePath string) bool {
	fileExt := path.Ext(filePath)
	if fileExt == ".yaml" || fileExt == ".yml" || fileExt == ".json" {
		return true
	}
	return false
}

func GetImportPathAndVersion(app *sysl.Application) (importPath string, version string, err error) {
	specURL, ok := app.Attrs["redoc-spec"]
	if ok {
		// Fetch the OpenAPI spec file into the cached ~/.sysl directory
		remoteFileMod, err := mod.Retrieve(specURL.GetS(), "")
		if err != nil {
			return "", "", err
		}
		importPath = specURL.GetS()
		version = remoteFileMod.Version
	} else {
		importPath = app.SourceContext.GetFile()
		version = app.SourceContext.GetVersion()
	}
	return importPath, version, nil
}

// BuildSpecURL takes a filepath and version and builds a URL to the cached spec file
// It also trims . prefixes and adds a / so that the URL is relative
func BuildSpecURL(filePath string, version string) string {
	// Append the version tag to the repo name
	if version != "" {
		filePath = AppendVersion(filePath, version)
	}
	filePath = strings.TrimPrefix(filePath, ".")
	if !strings.HasPrefix(filePath, "/") {
		filePath = "/" + filePath
	}
	return filePath
}

// AppendVersion takes in a remote file import path and a version, and appends the version tag to the repo name separated by the '@' char
// e.g for an input
// - remoteFilePath github.com/anz-bank/sysl-examples/demos/grocerystore/grocerystore.sysl
// - version v0.0.0-c63b9e92813a
// it returns /github.com/anz-bank/sysl-examples@v0.0.0-c63b9e92813a/demos/grocerystore/grocerystore.sysl
func AppendVersion(remoteFilePath string, version string) string {
	names := strings.FieldsFunc(remoteFilePath, func(c rune) bool {
		return c == '/'
	})
	if len(names) < 3 {
		return ""
	}
	names[2] = names[2] + "@" + version
	return path.Join(names...)
}

// GetRemoteFromGit gets the URL to the git remote
// e.g github.com/myorg/somerepo/
func GetRemoteFromGit() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	out, err := cmd.CombinedOutput()
	if err != nil {

		return "", fmt.Errorf("error getting git remote: is sysl-catalog running in a git repo? %w", err)
	}
	return StripExtension(string(out)), nil
}

// StripExtension removes spaces and suffixes
func StripExtension(input string) string {
	noExt := strings.TrimSuffix(input, path.Ext(input))
	noSpace := strings.TrimSpace(noExt)
	return noSpace
}

// BuildGithubRawURL gets the base URL for raw content hosted on github.com or Github Enterprise
// For github.com it takes in https://github.com/anz-bank/sysl-catalog and returns https://raw.githubusercontent.com/anz-bank/sysl-catalog/master/
// For Github Enterprise it takes in https://github.myorg.com/anz-bank/sysl-catalog and returns https://github.myorg.com/raw/anz-bank/sysl-catalog/master/
func BuildGithubRawURL(repoURL string) (gitURL string) {
	url, err := url.Parse(repoURL)
	if err != nil {
		panic(err)
	}
	switch url.Host {
	case "github.com":
		url.Host = "raw.githubusercontent.com"
		url.Path = url.Path + "/master/"
		gitURL = url.String()
	default:
		// Handles github enterprise which uses a different URL scheme for raw files
		url.Path = "raw" + url.Path + "/master/"
		gitURL = url.String()
	}
	return gitURL
}

// BuildGithubBlobURL creates a root URL for github blob
// it will not work for non github links.
func BuildGithubBlobURL(repoURL string) string {
	url, err := url.Parse(repoURL)
	if err != nil {
		panic(err)
	}
	url.Path = path.Join(url.Path, "/blob/master/")
	return url.String()
}
