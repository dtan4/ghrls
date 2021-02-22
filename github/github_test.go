package github

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-github/v33/github"
)

type fakeRepositoriesService struct{}

func (s fakeRepositoriesService) GetReleaseByTag(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error) {
	tagName := "v1"
	body := "The quick brown fox jumps over the lazy dog"
	assetURL := "https://github.com/owner/repo/releases/download/v1/darwin.tar.gz"
	login := "dtan4"
	name := "v1"
	htmlURL := "https://github.com/owner/repo/releases/tag/v1"

	return &github.RepositoryRelease{
		Assets: []*github.ReleaseAsset{
			&github.ReleaseAsset{
				BrowserDownloadURL: &assetURL,
			},
		},
		Author: &github.User{
			Login: &login,
		},
		Body:        &body,
		CreatedAt:   &github.Timestamp{time.Date(2018, 12, 13, 0, 30, 24, 0, time.UTC)},
		HTMLURL:     &htmlURL,
		Name:        &name,
		PublishedAt: &github.Timestamp{time.Date(2018, 12, 14, 0, 30, 24, 0, time.UTC)},
		TagName:     &tagName,
	}, &github.Response{}, nil
}

func (s fakeRepositoriesService) GetCommit(ctx context.Context, owner, repo, sha string) (*github.RepositoryCommit, *github.Response, error) {
	commitSHA := "856abeb2b507fc1db16dcaea938775ff938a5355"

	return &github.RepositoryCommit{
		SHA: &commitSHA,
	}, &github.Response{}, nil
}

func (s fakeRepositoriesService) ListReleases(ctx context.Context, owner, repo string, opt *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error) {
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

func (s fakeRepositoriesService) ListTags(ctx context.Context, owner string, repo string, opt *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error) {
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

func TestNewClient(t *testing.T) {
	testcases := []struct {
		accessToken string
	}{
		{accessToken: ""},
		{accessToken: "dummyaccesstoken"},
	}

	for _, tc := range testcases {
		c := NewClient(tc.accessToken)

		if c == nil {
			t.Error("want: object, got: nil")
		}
	}
}

func TestDescribeRelease(t *testing.T) {
	c := &Client{
		repositories: fakeRepositoriesService{},
	}

	owner := "owner"
	repo := "repo"
	tag := "v1"

	want := &Tag{
		Name: "v1",
		Release: &Release{
			ArtifactURLs: []string{
				"https://github.com/owner/repo/releases/download/v1/darwin.tar.gz",
			},
			Author:      "dtan4",
			Body:        "The quick brown fox jumps over the lazy dog",
			Commit:      "856abeb2b507fc1db16dcaea938775ff938a5355",
			CreatedAt:   time.Date(2018, 12, 13, 0, 30, 24, 0, time.UTC),
			Name:        "v1",
			PublishedAt: time.Date(2018, 12, 14, 0, 30, 24, 0, time.UTC),
			URL:         "https://github.com/owner/repo/releases/tag/v1",
		},
	}

	got, err := c.DescribeRelease(context.Background(), owner, repo, tag)
	if err != nil {
		t.Errorf("want no error, got: %#v", err)
	}

	if !reflect.DeepEqual(*got, *want) {
		t.Errorf("want: %#v, got: %#v", *want, *got)
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

	got, err := c.ListTagsAndReleases(context.Background(), owner, repo)
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
