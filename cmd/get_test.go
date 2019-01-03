package cmd

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/dtan4/ghrls/github"
)

type fakeClientForGet struct {
	Tag *github.Tag
	Err error
}

func (c fakeClientForGet) DescribeRelease(owner, repo, tag string) (*github.Tag, error) {
	if c.Err != nil {
		return nil, c.Err
	}

	return c.Tag, nil
}

func (c fakeClientForGet) ListTagsAndReleases(owner, repo string) ([]*github.Tag, error) {
	return []*github.Tag{}, nil
}

func TestRunTag_success(t *testing.T) {
	gmt, err := time.LoadLocation("Europe/London")
	if err != nil {
		t.Fatal(err)
	}

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	testcases := []struct {
		args     []string
		timezone *time.Location
		want     string
	}{
		{
			args:     []string{"owner/repo", "v1.5.2"},
			timezone: gmt,
			want: "" +
				"Tag:         v1\n" +
				"Commit:      856abeb2b507fc1db16dcaea938775ff938a5355\n" +
				"Name:        v1\n" +
				"Author:      dtan4\n" +
				"CreatedAt:   2018-12-13 00:30:24 +0000 GMT\n" +
				"PublishedAt: 2018-12-14 00:30:24 +0000 GMT\n" +
				"URL:         https://github.com/owner/repo/releases/tag/v1\n" +
				"Artifacts:   https://github.com/owner/repo/releases/download/v1/darwin.tar.gz\n" +
				"\n" +
				"The quick brown fox jumps over the lazy dog\n",
		},
		{
			args:     []string{"owner/repo", "v1.5.2"},
			timezone: jst,
			want: "" +
				"Tag:         v1\n" +
				"Commit:      856abeb2b507fc1db16dcaea938775ff938a5355\n" +
				"Name:        v1\n" +
				"Author:      dtan4\n" +
				"CreatedAt:   2018-12-13 09:30:24 +0900 JST\n" +
				"PublishedAt: 2018-12-14 09:30:24 +0900 JST\n" +
				"URL:         https://github.com/owner/repo/releases/tag/v1\n" +
				"Artifacts:   https://github.com/owner/repo/releases/download/v1/darwin.tar.gz\n" +
				"\n" +
				"The quick brown fox jumps over the lazy dog\n",
		},
	}

	client := fakeClientForGet{
		Tag: &github.Tag{
			Name: "v1",
			Release: &github.Release{
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
		},
	}

	for _, tc := range testcases {
		stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)

		if err := RunGet(stdout, stderr, tc.args, client, tc.timezone); err != nil {
			t.Errorf("want: no error, got: %#v", err)
		}

		if stdout.String() != tc.want {
			t.Errorf("stderr want:\n%q\ngot:\n%q", tc.want, stdout.String())
		}
	}
}

func TestRunGet_invalidArgs(t *testing.T) {
	testcases := []struct {
		args []string
		want string
	}{
		{
			args: []string{},
			want: "Please specify repository <user/name> and tag.",
		},
		{
			args: []string{"owner/repo"},
			want: "Please specify repository <user/name> and tag.",
		},
		{
			args: []string{"owner/repo", "v1", "foobar"},
			want: "Please specify repository <user/name> and tag.",
		},
		{
			args: []string{"owner repo", "v1"},
			want: "Invalid repository name: owner repo",
		},
	}

	gmt, err := time.LoadLocation("Europe/London")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testcases {
		stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)

		err := RunGet(stdout, stderr, tc.args, fakeClientForGet{}, gmt)

		if err == nil {
			t.Error("want: error, got: nil")
		}

		if err.Error() != tc.want {
			t.Errorf("error want: %q, got: %q", tc.want, err.Error())
		}
	}
}

func TestRunGet_notfound(t *testing.T) {
	testcases := []struct {
		args []string
		err  error
		want string
	}{
		{
			args: []string{"owner/repo", "v1"},
			err:  fmt.Errorf("404 Not Found"),
			want: "owner/repo@v1 : not found",
		},
		{
			args: []string{"owner/repo", "v1"},
			err:  fmt.Errorf("unexpected error"),
			want: "unexpected error",
		},
	}

	gmt, err := time.LoadLocation("Europe/London")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testcases {
		stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
		client := fakeClientForGet{
			Err: tc.err,
		}

		err := RunGet(stdout, stderr, tc.args, client, gmt)

		if err == nil {
			t.Error("want: error, got: nil")
		}

		if err.Error() != tc.want {
			t.Errorf("error want: %q, got: %q", tc.want, err.Error())
		}
	}
}
