package github

import (
	"net/http"

	"github.com/google/go-github/github"
)

// Client represents a wrapper of GitHub API client
type Client struct {
	client *github.Client
}

// NewClient creates new Client object
func NewClient(httpClient *http.Client) *Client {
	return &Client{
		client: github.NewClient(httpClient),
	}
}

// ListReleases lists all releases of the given repository
func (c *Client) ListReleases(owner, repo string) ([]*github.RepositoryRelease, error) {
	allReleases := []*github.RepositoryRelease{}

	listOpts := &github.ListOptions{}

	for {
		releases, resp, err := c.client.Repositories.ListReleases(owner, repo, listOpts)
		if err != nil {
			return []*github.RepositoryRelease{}, err
		}

		allReleases = append(allReleases, releases...)

		if resp.NextPage == 0 {
			break
		}

		listOpts.Page = resp.NextPage
	}

	return allReleases, nil
}

// ListTags lists all tags of the given repository
func (c *Client) ListTags(owner, repo string) ([]*github.RepositoryTag, error) {
	allTags := []*github.RepositoryTag{}

	listOpts := &github.ListOptions{}

	for {
		tags, resp, err := c.client.Repositories.ListTags(owner, repo, listOpts)
		if err != nil {
			return []*github.RepositoryTag{}, err
		}

		allTags = append(allTags, tags...)

		if resp.NextPage == 0 {
			break
		}

		listOpts.Page = resp.NextPage
	}

	return allTags, nil
}
