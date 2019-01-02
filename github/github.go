package github

import (
	"net/http"
	"time"

	"github.com/google/go-github/github"
)

const (
	// default: 30, max: 100
	// https://developer.github.com/v3/#pagination
	perPage = 100
)

type RepositoriesServiceInterface interface {
	GetReleaseByTag(owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error)
	GetCommit(owner, repo, sha string) (*github.RepositoryCommit, *github.Response, error)
	ListReleases(owner, repo string, opt *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error)
	ListTags(owner string, repo string, opt *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error)
}

// Client represents a wrapper of GitHub API client
type Client struct {
	repositories RepositoriesServiceInterface
}

type Release struct {
	Name      string
	Tag       string
	CreatedAt time.Time
}

// NewClient creates new Client object
func NewClient(httpClient *http.Client) *Client {
	return &Client{
		repositories: github.NewClient(httpClient).Repositories,
	}
}

// MakeReleasesMap makes map of tag name and release
func MakeReleasesMap(releases []*github.RepositoryRelease) map[string]*github.RepositoryRelease {
	result := map[string]*github.RepositoryRelease{}

	for _, release := range releases {
		result[*release.TagName] = release
	}

	return result
}

// GetRelease returns release metadata of the given tag
func (c *Client) GetRelease(owner, repo, tag string) (*github.RepositoryRelease, error) {
	release, _, err := c.repositories.GetReleaseByTag(owner, repo, tag)
	if err != nil {
		return nil, err
	}

	return release, nil
}

// GetTag returns commit metadata of the given tag
func (c *Client) GetTagCommit(owner, repo, tag string) (*github.RepositoryCommit, error) {
	commit, _, err := c.repositories.GetCommit(owner, repo, tag)
	if err != nil {
		return nil, err
	}

	return commit, nil
}

// ListReleasesNew retrieves all tags and releases of the given repository
func (c *Client) ListReleasesNew(owner, repo string) ([]*Release, error) {
	releases, err := c.ListReleases(owner, repo)
	if err != nil {
		return []*Release{}, err
	}

	releasesMap := MakeReleasesMap(releases)

	tags, err := c.ListTags(owner, repo)
	if err != nil {
		return []*Release{}, err
	}

	rs := []*Release{}

	for _, tag := range tags {
		var release *Release

		if r, ok := releasesMap[*tag.Name]; ok {
			var name string

			if r.Name == nil {
				name = ""
			} else {
				name = *r.Name
			}

			t := *r.CreatedAt

			release = &Release{
				Name:      name,
				Tag:       *tag.Name,
				CreatedAt: time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()),
			}
		} else {
			release = &Release{
				Tag: *tag.Name,
			}
		}

		rs = append(rs, release)
	}

	return rs, nil
}

// ListReleases lists all releases of the given repository
func (c *Client) ListReleases(owner, repo string) ([]*github.RepositoryRelease, error) {
	allReleases := []*github.RepositoryRelease{}

	listOpts := &github.ListOptions{
		PerPage: perPage,
	}

	for {
		releases, resp, err := c.repositories.ListReleases(owner, repo, listOpts)
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

	listOpts := &github.ListOptions{
		PerPage: perPage,
	}

	for {
		tags, resp, err := c.repositories.ListTags(owner, repo, listOpts)
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
