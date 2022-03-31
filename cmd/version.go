package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string = "undefined"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "displays command version",
	Long:  `displays command version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
