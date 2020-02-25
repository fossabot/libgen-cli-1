package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/binodsh/libgen"
	"github.com/spf13/cobra"
)

var SearchStr string

var rootCmd = &cobra.Command{
	Use:   "libgen",
	Short: "libgen is cli tool to search book in libgen",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if strings.TrimSpace(SearchStr) != "" {
			books := libgen.SearchBookByTitle(SearchStr)

			for _, book := range books {
				fmt.Printf("Author: %s\nTitle: %s\nYear: %s\n\n", book.Author, book.Title, book.Year)
			}
		}
	},
}

func Execute() {
	rootCmd.Flags().StringVarP(&SearchStr, "search", "s", "", "search string")
	rootCmd.MarkFlagRequired("search")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
