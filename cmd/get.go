package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/dtan4/ghrls/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
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
	RunE: doGet,
}

func doGet(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("Please specify repository <user/name> and tag.")
	}

	ss := strings.Split(args[0], "/")
	if len(ss) != 2 {
		return fmt.Errorf("Invalid repository name: %s", args[0])
	}
	owner, repo := ss[0], ss[1]

	tag := args[1]

	var httpClient *http.Client

	if rootOpts.GitHubToken == "" {
		httpClient = nil
	} else {
		ts := oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: rootOpts.GitHubToken,
		})
		httpClient = oauth2.NewClient(oauth2.NoContext, ts)
	}

	client := github.NewClient(httpClient)

	release, err := client.GetRelease(owner, repo, tag)
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			return fmt.Errorf("%s/%s@%s : not found", owner, repo, tag)
		}
		return err
	}

	commit, err := client.GetTagCommit(owner, repo, tag)
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			return fmt.Errorf("%s/%s@%s : not found", owner, repo, tag)
		}
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "Tag:\t"+*release.TagName)
	fmt.Fprintln(w, "Commit:\t"+*commit.SHA)

	if release.Name == nil {
		fmt.Fprintln(w, "Name:\t")
	} else {
		fmt.Fprintln(w, "Name:\t"+*release.Name)
	}

	fmt.Fprintln(w, "Author:\t"+*release.Author.Login)
	fmt.Fprintln(w, "CreatedAt:\t"+release.CreatedAt.Local().String())
	fmt.Fprintln(w, "PublishedAt:\t"+release.PublishedAt.Local().String())
	fmt.Fprintln(w, "URL:\t"+*release.HTMLURL)

	if len(release.Assets) > 0 {
		fmt.Fprintln(w, "Assets:\t"+*release.Assets[0].BrowserDownloadURL)

		for _, asset := range release.Assets[1:] {
			fmt.Fprintln(w, "\t"+*asset.BrowserDownloadURL)
		}
	} else {
		fmt.Fprintln(w, "Assets:\t")
	}

	w.Flush()

	if release.Body != nil {
		fmt.Println("")
		fmt.Println(*release.Body)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(getCmd)
}
