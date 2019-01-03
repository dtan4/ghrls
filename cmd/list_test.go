package cmd

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/dtan4/ghrls/github"
)

type fakeClientForList struct {
	Tags []*github.Tag
	Err  error
}

func (c fakeClientForList) DescribeRelease(owner, repo, tag string) (*github.Tag, error) {
	return &github.Tag{}, nil
}

func (c fakeClientForList) ListTagsAndReleases(owner, repo string) ([]*github.Tag, error) {
	if c.Err != nil {
		return []*github.Tag{}, c.Err
	}

	return c.Tags, nil
}

func TestRunList_success(t *testing.T) {
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
			args:     []string{"owner/repo"},
			timezone: gmt,
			want: "" +
				"TAG               TYPE           CREATEDAT                        NAME\n" +
				"v1.6.0-alpha.0    TAG                                             \n" +
				"v1.5.3-beta.0     TAG                                             \n" +
				"v1.5.2            TAG+RELEASE    2017-01-12 04:51:15 +0000 GMT    v1.5.2\n",
		},
		{
			args:     []string{"owner/repo"},
			timezone: jst,
			want: "" +
				"TAG               TYPE           CREATEDAT                        NAME\n" +
				"v1.6.0-alpha.0    TAG                                             \n" +
				"v1.5.3-beta.0     TAG                                             \n" +
				"v1.5.2            TAG+RELEASE    2017-01-12 13:51:15 +0900 JST    v1.5.2\n",
		},
	}

	client := fakeClientForList{
		Tags: []*github.Tag{
			&github.Tag{
				Name:    "v1.6.0-alpha.0",
				Release: nil,
			},
			&github.Tag{
				Name:    "v1.5.3-beta.0",
				Release: nil,
			},
			&github.Tag{
				Name: "v1.5.2",
				Release: &github.Release{
					CreatedAt: time.Date(2017, 1, 12, 4, 51, 15, 0, time.UTC),
					Name:      "v1.5.2",
				},
			},
		},
	}

	for _, tc := range testcases {
		stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)

		if err := RunList(stdout, stderr, tc.args, client, tc.timezone); err != nil {
			t.Errorf("want: no error, got: %#v", err)
		}

		if stdout.String() != tc.want {
			t.Errorf("stderr want:\n%q\ngot:\n%q", tc.want, stdout.String())
		}
	}
}

func TestRunList_invalidArgs(t *testing.T) {
	testcases := []struct {
		args []string
		want string
	}{
		{
			args: []string{},
			want: "Please specify repository <user/name>.",
		},
		{
			args: []string{"owner", "repo"},
			want: "Please specify repository <user/name>.",
		},
		{
			args: []string{"owner repo"},
			want: "Invalid repository name: owner repo",
		},
	}

	gmt, err := time.LoadLocation("Europe/London")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testcases {
		stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)

		err := RunList(stdout, stderr, tc.args, fakeClientForList{}, gmt)

		if err == nil {
			t.Error("want: error, got: nil")
		}

		if err.Error() != tc.want {
			t.Errorf("error want: %q, got: %q", tc.want, err.Error())
		}
	}
}

func TestRunList_notfound(t *testing.T) {
	testcases := []struct {
		args []string
		err  error
		want string
	}{
		{
			args: []string{"owner/repo"},
			err:  fmt.Errorf("404 Not Found"),
			want: "owner/repo: not found",
		},
		{
			args: []string{"owner/repo"},
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
		client := fakeClientForList{
			Err: tc.err,
		}

		err := RunList(stdout, stderr, tc.args, client, gmt)

		if err == nil {
			t.Error("want: error, got: nil")
		}

		if err.Error() != tc.want {
			t.Errorf("error want: %q, got: %q", tc.want, err.Error())
		}
	}
}
