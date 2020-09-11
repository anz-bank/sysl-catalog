package catalog

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/anz-bank/pkg/mod"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/spf13/afero"
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
func GetImportPathAndVersion(app *sysl.Application) (importPath string, version string, err error) {
	specURL, ok := app.Attrs["redoc-spec"]
	if ok {
		homeDir, _ := os.UserHomeDir()
		// Fetch the OpenAPI spec file into the cached ~/.sysl directory
		err2 := mod.Config("github", mod.GoModulesOptions{}, mod.GitHubOptions{CacheDir: homeDir + "/.sysl", AccessToken: os.Getenv("SYSL_GITHUB_TOKEN"), Fs: afero.NewOsFs()}) // Setup sysl module in Github mode
		if err2 != nil {
			return "", "", err2
		}
		name, ver := mod.ExtractVersion(specURL.GetS())
		remoteFileMod, err := mod.Retrieve(name, ver)
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
	// Remove trailing @tag or @branchname
	repoPath, _ := mod.ExtractVersion(filePath)
	// Append the version tag to the repo name
	if version != "" {
		repoPath = AppendVersion(repoPath, version)
	}
	repoPath = strings.TrimPrefix(repoPath, ".")
	if !strings.HasPrefix(repoPath, "/") {
		repoPath = "/" + repoPath
	}
	return repoPath
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
