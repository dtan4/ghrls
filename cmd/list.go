package cmd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dtan4/ghrls/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: doList,
}

func doList(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Please specify repository <user/name>.")
	}

	ss := strings.Split(args[0], "/")
	if len(ss) != 2 {
		return fmt.Errorf("Invalid repository name: %s", args[0])
	}
	owner, repo := ss[0], ss[1]

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

	fmt.Println("=== tags")

	tags, err := client.ListTags(owner, repo)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		fmt.Println(*tag.Name)
	}

	fmt.Println("")
	fmt.Println("=== releases")

	releases, err := client.ListReleases(owner, repo)
	if err != nil {
		return err
	}

	for _, release := range releases {
		fmt.Println(*release.TagName, *release.CreatedAt)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(listCmd)
}
