package version

import (
	"context"
	"encoding/json"

	"github.com/google/go-github/v32/github"
)

func GitTagsJson(org, repo string) (string, error) {
	client := github.NewClient(nil)

	tags, _, err := client.Repositories.ListTags(context.Background(), org, repo, nil)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(tags[0].Commit.Author)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
