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
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		return err
	}

	commit, err := client.GetTagCommit(owner, repo, tag)
	if err != nil {
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
