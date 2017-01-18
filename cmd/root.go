package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	SilenceErrors: true,
	SilenceUsage:  true,
	Use:           "ghrls",
	Short:         "List & Describe GitHub Releases",
}

var rootOpts = struct {
	GitHubToken string
}{}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if os.Getenv("GITHUB_TOKEN") != "" {
		rootOpts.GitHubToken = os.Getenv("GITHUB_TOKEN")
	}
}
