package cmd

import (
	"fmt"

	"github.com/dtan4/ghrls/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.String())
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
