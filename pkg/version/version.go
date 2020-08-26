package version

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type Tag struct {
	Name string    `json:"name"`
	SHA  string    `json:"sha"`
	Date time.Time `json:"date"`
}

// ListTags fetches git repo tags via GitHub API. It return tags in order.
// Could support pagination. Return up to 30 by default
func ListTags(owner, repo string) (*[]Tag, error) {
	client := github.NewClient(nil)
	githubToken := os.Getenv("SYSL_GITHUB_TOKEN")
	if githubToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubToken},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		client = github.NewClient(tc)
	}

	ctx := context.Background()

	// rawtags, _, err := client.Repositories.ListTags(ctx, owner, repo, &github.ListOptions{PerPage: 100})
	rawtags, _, err := client.Repositories.ListTags(ctx, owner, repo, nil)
	if err != nil {
		return nil, err
	}

	tags := make([]Tag, 0, len(rawtags))
	for _, tag := range rawtags {
		sha := tag.GetCommit().GetSHA()
		rc, _, err := client.Repositories.GetCommit(ctx, owner, repo, sha)
		if err != nil {
			return nil, err
		}

		tags = append(tags, Tag{
			Name: tag.GetName(),
			SHA:  sha,
			Date: rc.GetCommit().GetAuthor().GetDate(),
		})
	}
	return &tags, nil
}

// GitListTags fetches git repo tags via git. It returns disordered tags.
func GitListTags(url string) (*[]Tag, error) {
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return nil, err
	}

	iter, err := r.Tags()
	if err != nil {
		return nil, err
	}

	tags := make([]Tag, 0)
	err = iter.ForEach(func(tag *plumbing.Reference) error {
		sha := tag.Hash().String()
		h, err := r.ResolveRevision(plumbing.Revision(sha))
		if err != nil {
			return err
		}
		if h == nil {
			return errors.New("revision not found")
		}

		obj, err := r.Object(plumbing.AnyObject, *h)
		if err != nil {
			return err
		}

		o := obj.(*object.Commit)

		tags = append(tags, Tag{
			Name: strings.TrimPrefix(tag.Name().String(), "refs/tags/"),
			SHA:  sha,
			Date: o.Author.When,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &tags, nil
}
