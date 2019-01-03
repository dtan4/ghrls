package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/dtan4/ghrls/github"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get REPOSITORY TAG",
	Short: "Describe release information",
	Long: `Describe release information

Example:

$ ghrls get kubernetes/kubernetes v1.5.2
Tag:         v1.5.2
Commit:      08e099554f3c31f6e6f07b448ab3ed78d0520507
Name:        v1.5.2
Author:      saad-ali
CreatedAt:   2017-01-12 13:51:15 +0900 JST
PublishedAt: 2017-01-12 16:25:50 +0900 JST
URL:         https://github.com/kubernetes/kubernetes/releases/tag/v1.5.2
Assets:      https://github.com/kubernetes/kubernetes/releases/download/v1.5.2/kubernetes.tar.gz

See [kubernetes-announce@](https://groups.google.com/forum/#!forum/kubernetes-announce) and [CHANGELOG](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG.md#v152) for details.
` +
		"SHA256 for `kubernetes.tar.gz`: `67344958325a70348db5c4e35e59f9c3552232cdc34defb8a0a799ed91c671a3`" +
		`
Additional binary downloads are linked in the [CHANGELOG](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG.md#downloads-for-v152).
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		timezone := time.Local
		client := github.NewClient(rootOpts.GitHubToken)

		return RunGet(os.Stderr, os.Stdout, args, client, timezone)
	},
}

func RunGet(stdout, stderr io.Writer, args []string, client github.ClientInterface, timezone *time.Location) error {
	if len(args) != 2 {
		return fmt.Errorf("Please specify repository <user/name> and tag.")
	}

	ss := strings.Split(args[0], "/")
	if len(ss) != 2 {
		return fmt.Errorf("Invalid repository name: %s", args[0])
	}
	owner, repo := ss[0], ss[1]

	tag := args[1]

	t, err := client.DescribeRelease(owner, repo, tag)
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			return fmt.Errorf("%s/%s@%s : not found", owner, repo, tag)
		}
		return err
	}

	w := tabwriter.NewWriter(stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "Tag:\t"+t.Name)
	fmt.Fprintln(w, "Commit:\t"+t.Release.Commit)
	fmt.Fprintln(w, "Name:\t"+t.Release.Name)
	fmt.Fprintln(w, "Author:\t"+t.Release.Author)
	fmt.Fprintln(w, "CreatedAt:\t"+t.Release.CreatedAt.In(timezone).String())
	fmt.Fprintln(w, "PublishedAt:\t"+t.Release.PublishedAt.In(timezone).String())
	fmt.Fprintln(w, "URL:\t"+t.Release.URL)

	if len(t.Release.ArtifactURLs) > 0 {
		fmt.Fprintln(w, "Artifacts:\t"+t.Release.ArtifactURLs[0])

		for _, url := range t.Release.ArtifactURLs[1:] {
			fmt.Fprintln(w, "\t"+url)
		}
	}

	w.Flush()

	if t.Release.Body != "" {
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, t.Release.Body)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(getCmd)
}
