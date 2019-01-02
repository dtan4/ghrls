package github

import (
	"testing"

	"github.com/google/go-github/github"
)

type fakeRepositoriesService struct{}

func (s fakeRepositoriesService) GetReleaseByTag(owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error) {
	return &github.RepositoryRelease{}, &github.Response{}, nil
}

func (s fakeRepositoriesService) GetCommit(owner, repo, sha string) (*github.RepositoryCommit, *github.Response, error) {
	return &github.RepositoryCommit{}, &github.Response{}, nil
}

func (s fakeRepositoriesService) ListReleases(owner, repo string, opt *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error) {
	return []*github.RepositoryRelease{}, &github.Response{}, nil
}

func (s fakeRepositoriesService) ListTags(owner string, repo string, opt *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error) {
	return []*github.RepositoryTag{}, &github.Response{}, nil
}

func TestGetRelease(t *testing.T) {
	c := &Client{
		repositories: fakeRepositoriesService{},
	}

	owner := "owner"
	repo := "repo"
	tag := "tag"

	if _, err := c.GetRelease(owner, repo, tag); err != nil {
		t.Errorf("want no error, got %#v", err)
	}
}

func TestGetTagCommit(t *testing.T) {
	c := &Client{
		repositories: fakeRepositoriesService{},
	}

	owner := "owner"
	repo := "repo"
	tag := "tag"

	if _, err := c.GetTagCommit(owner, repo, tag); err != nil {
		t.Errorf("want no error, got %#v", err)
	}
}

func TestListReleases(t *testing.T) {
	c := &Client{
		repositories: fakeRepositoriesService{},
	}

	owner := "owner"
	repo := "repo"

	if _, err := c.ListReleases(owner, repo); err != nil {
		t.Errorf("want no error, got %#v", err)
	}
}

func TestListTags(t *testing.T) {
	c := &Client{
		repositories: fakeRepositoriesService{},
	}

	owner := "owner"
	repo := "repo"

	if _, err := c.ListTags(owner, repo); err != nil {
		t.Errorf("want no error, got %#v", err)
	}
}
