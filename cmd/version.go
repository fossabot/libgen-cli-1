package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of libgen cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("libgen cli v0.0.1")
	},
}
