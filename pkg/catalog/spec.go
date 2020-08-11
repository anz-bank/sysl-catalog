package catalog

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

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

func IsOpenAPIFile(source *sysl.SourceContext) bool {
	importPath := source.GetFile()
	fileExt := path.Ext(importPath)
	if fileExt == ".yaml" || fileExt == ".json" {
		return true
	}
	return false
}

// BuildSpecURL takes a source context reference and builds an raw git URL for it
// It handles sourceContext paths which are from remote repos as well as in the same repo
func BuildSpecURL(source *sysl.SourceContext) (string, error) {
	filePath := source.GetFile()
	filePath = strings.TrimPrefix(filePath, ".")
	if !strings.HasPrefix(filePath, "/") {
		filePath = "/" + filePath
	}
	return filePath, nil
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
