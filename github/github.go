package github

import (
	"context"
	"net/http"
	"time"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

const (
	// default: 30, max: 100
	// https://developer.github.com/v3/#pagination
	perPage = 100
)

type Release struct {
	ArtifactURLs []string
	Author       string
	Body         string
	Commit       string
	CreatedAt    time.Time
	Name         string
	PublishedAt  time.Time
	URL          string
}

type Tag struct {
	Name    string
	Release *Release
}

type RepositoriesServiceInterface interface {
	GetReleaseByTag(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error)
	GetCommit(ctx context.Context, owner, repo, sha string) (*github.RepositoryCommit, *github.Response, error)
	ListReleases(ctx context.Context, owner, repo string, opt *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error)
	ListTags(ctx context.Context, owner string, repo string, opt *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error)
}

type ClientInterface interface {
	DescribeRelease(ctx context.Context, owner, repo, tag string) (*Tag, error)
	ListTagsAndReleases(ctx context.Context, owner, repo string) ([]*Tag, error)
}

// Client represents a wrapper of GitHub API client
type Client struct {
	repositories RepositoriesServiceInterface
}

// NewClient creates new Client object
func NewClient(accessToken string) *Client {
	var hc *http.Client

	if accessToken == "" {
		hc = nil
	} else {
		ts := oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: accessToken,
		})
		hc = oauth2.NewClient(oauth2.NoContext, ts)
	}

	return &Client{
		repositories: github.NewClient(hc).Repositories,
	}
}

// DescribeRelease returns detail of the given release
func (c *Client) DescribeRelease(ctx context.Context, owner, repo, tag string) (*Tag, error) {
	release, err := c.getRelease(ctx, owner, repo, tag)
	if err != nil {
		return nil, err
	}

	commit, err := c.getTagCommit(ctx, owner, repo, tag)
	if err != nil {
		return nil, err
	}

	var body, name string

	if release.Body == nil {
		body = ""
	} else {
		body = *release.Body
	}

	if release.Name == nil {
		name = ""
	} else {
		name = *release.Name
	}

	artifactURLs := []string{}

	for _, asset := range release.Assets {
		artifactURLs = append(artifactURLs, *asset.BrowserDownloadURL)
	}

	createdAt := *release.CreatedAt
	publishedAt := *release.PublishedAt

	return &Tag{
		Name: *release.TagName,
		Release: &Release{
			ArtifactURLs: artifactURLs,
			Author:       *release.Author.Login,
			Body:         body,
			Commit:       *commit.SHA,
			CreatedAt: time.Date(
				createdAt.Year(),
				createdAt.Month(),
				createdAt.Day(),
				createdAt.Hour(),
				createdAt.Minute(),
				createdAt.Second(),
				createdAt.Nanosecond(),
				createdAt.Location(),
			),
			Name: name,
			PublishedAt: time.Date(
				publishedAt.Year(),
				publishedAt.Month(),
				publishedAt.Day(),
				publishedAt.Hour(),
				publishedAt.Minute(),
				publishedAt.Second(),
				publishedAt.Nanosecond(),
				publishedAt.Location(),
			),
			URL: *release.HTMLURL,
		},
	}, nil
}

func (c *Client) getRelease(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, error) {
	release, _, err := c.repositories.GetReleaseByTag(ctx, owner, repo, tag)
	if err != nil {
		return nil, err
	}

	return release, nil
}

func (c *Client) getTagCommit(ctx context.Context, owner, repo, tag string) (*github.RepositoryCommit, error) {
	commit, _, err := c.repositories.GetCommit(ctx, owner, repo, tag)
	if err != nil {
		return nil, err
	}

	return commit, nil
}

// ListTagsAndReleases retrieves all tags and releases of the given repository
func (c *Client) ListTagsAndReleases(ctx context.Context, owner, repo string) ([]*Tag, error) {
	tags, err := c.listTags(ctx, owner, repo)
	if err != nil {
		return []*Tag{}, err
	}

	releases, err := c.listReleases(ctx, owner, repo)
	if err != nil {
		return []*Tag{}, err
	}

	releasesMap := map[string]*github.RepositoryRelease{}

	for _, release := range releases {
		releasesMap[*release.TagName] = release
	}

	ts := []*Tag{}

	for _, t := range tags {
		var tag *Tag

		if r, ok := releasesMap[*t.Name]; ok {
			var name string

			if r.Name == nil {
				name = ""
			} else {
				name = *r.Name
			}

			createdAt := *r.CreatedAt

			tag = &Tag{
				Name: *t.Name,
				Release: &Release{
					Name: name,
					CreatedAt: time.Date(
						createdAt.Year(),
						createdAt.Month(),
						createdAt.Day(),
						createdAt.Hour(),
						createdAt.Minute(),
						createdAt.Second(),
						createdAt.Nanosecond(),
						createdAt.Location(),
					),
				},
			}
		} else {
			tag = &Tag{
				Name: *t.Name,
			}
		}

		ts = append(ts, tag)
	}

	return ts, nil
}

// ListReleases lists all releases of the given repository
func (c *Client) listReleases(ctx context.Context, owner, repo string) ([]*github.RepositoryRelease, error) {
	allReleases := []*github.RepositoryRelease{}

	listOpts := &github.ListOptions{
		PerPage: perPage,
	}

	for {
		releases, resp, err := c.repositories.ListReleases(ctx, owner, repo, listOpts)
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
func (c *Client) listTags(ctx context.Context, owner, repo string) ([]*github.RepositoryTag, error) {
	allTags := []*github.RepositoryTag{}

	listOpts := &github.ListOptions{
		PerPage: perPage,
	}

	for {
		tags, resp, err := c.repositories.ListTags(ctx, owner, repo, listOpts)
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
