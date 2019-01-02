package github

import (
	"reflect"
	"testing"
	"time"

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
	tag_v1_13_2_beta_0 := "v1.13.2-beta.0"
	tag_v1_13_1 := "v1.13.1"
	release_v1_13_1 := "v1.13.1"

	return []*github.RepositoryRelease{
		&github.RepositoryRelease{
			TagName:   &tag_v1_13_2_beta_0,
			Name:      nil,
			CreatedAt: &github.Timestamp{time.Date(2018, 12, 14, 0, 30, 24, 0, time.UTC)},
		},
		&github.RepositoryRelease{
			TagName:   &tag_v1_13_1,
			Name:      &release_v1_13_1,
			CreatedAt: &github.Timestamp{time.Date(2018, 12, 13, 0, 30, 24, 0, time.UTC)},
		},
	}, &github.Response{}, nil
}

func (s fakeRepositoriesService) ListTags(owner string, repo string, opt *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error) {
	tag_v1_13_2_beta_1 := "v1.13.2-beta.1"
	tag_v1_13_2_beta_0 := "v1.13.2-beta.0"
	tag_v1_13_1 := "v1.13.1"

	return []*github.RepositoryTag{
		&github.RepositoryTag{
			Name: &tag_v1_13_2_beta_1,
		},
		&github.RepositoryTag{
			Name: &tag_v1_13_2_beta_0,
		},
		&github.RepositoryTag{
			Name: &tag_v1_13_1,
		},
	}, &github.Response{}, nil
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

func TestListTagsAndReleases(t *testing.T) {
	c := &Client{
		repositories: fakeRepositoriesService{},
	}

	owner := "owner"
	repo := "repo"

	want := []*Tag{
		&Tag{
			Name:    "v1.13.2-beta.1",
			Release: nil,
		},
		&Tag{
			Name: "v1.13.2-beta.0",
			Release: &Release{
				Name:      "",
				CreatedAt: time.Date(2018, 12, 14, 0, 30, 24, 0, time.UTC),
			},
		},
		&Tag{
			Name: "v1.13.1",
			Release: &Release{
				Name:      "v1.13.1",
				CreatedAt: time.Date(2018, 12, 13, 0, 30, 24, 0, time.UTC),
			},
		},
	}

	got, err := c.ListTagsAndReleases(owner, repo)
	if err != nil {
		t.Errorf("want no error, got: %#v", err)
	}

	if len(got) != len(want) {
		t.Errorf("want: %d items, got: %d items", len(want), len(got))
	}

	for i, g := range got {
		if !reflect.DeepEqual(*g, *want[i]) {
			t.Errorf("want: %#v, got: %#v", *want[i], *g)
		}
	}
}
