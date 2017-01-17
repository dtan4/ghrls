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

var (
	headers = []string{
		"TAG",
		"TYPE",
		"CREATEDAT",
		"NAME",
	}
)

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

	tags, err := client.ListTags(owner, repo)
	if err != nil {
		return err
	}

	releases, err := client.ListReleases(owner, repo)
	if err != nil {
		return err
	}

	releasesMap := github.MakeReleasesMap(releases)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	for _, tag := range tags {
		ss := []string{*tag.Name}

		if release, ok := releasesMap[*tag.Name]; ok {
			ss = append(ss, "TAG+RELEASE", release.CreatedAt.Local().String())

			if release.Name == nil {
				ss = append(ss, "")
			} else {
				ss = append(ss, *release.Name)
			}
		} else {
			ss = append(ss, "TAG", "", "")
		}

		fmt.Fprintln(w, strings.Join(ss, "\t"))
	}

	w.Flush()

	return nil
}

func init() {
	RootCmd.AddCommand(listCmd)
}
