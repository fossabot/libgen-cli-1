package cmd

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"

	"github.com/binodsh/libgen"
	"github.com/binodsh/libgen-cli/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(SearchCommand)

	SearchCommand.Flags().StringP("sort-by", "", "", "sort by (--sort-by author)")
	SearchCommand.Flags().StringP("sort-order", "", "", "sort order (--sort-order desc)")
}

//SearchCommand search the book
var SearchCommand = &cobra.Command{
	Use:   "search",
	Short: "search the book",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		searchStr := strings.Join(args, " ")
		searchOptions := libgen.SearchOptions{}
		if strings.TrimSpace(searchStr) != "" {
			searchOptions.Query = searchStr
		} else {
			fmt.Println("invalid search query")
			return
		}

		sortBy, _ := cmd.Flags().GetString("sort-by")
		sortMode, _ := cmd.Flags().GetString("sort-mode")
		if sortBy != "" {
			searchOptions.SortBy = sortBy

			if sortMode != "" {
				searchOptions.SortMode = sortMode
			}
		}

		books, err := libgen.Search(searchOptions)

		if err != nil {
			fmt.Println("error while searching book")
			return
		}

		var titles []string
		for _, book := range books {
			titles = append(titles, book.Title)
		}

		//promtui template
		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "\U0001F336 {{ .Title | cyan }}",
			Inactive: "  {{ .Title | cyan }}",
			Selected: "\U0001F336 {{ .Name | red | cyan }}",
			Details: `
				--------- Book Details ----------
				{{ "Name:" | faint }} {{ .Title }}
				{{ "@author:" | faint }} {{ .Author }}    {{ "@publisher:" | faint }} {{ .Publisher }}    {{ "@year:" | faint }} {{ .Year }}
				{{ "@extension:" | faint }} {{ .Extension }}    {{ "@pages:" | faint }} {{ .Pages }}`,
		}

		prompt := promptui.Select{
			Label:        "Select book to download",
			Items:        books,
			Templates:    templates,
			HideSelected: true,
		}

		i, _, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		err = internal.Download(books[i])
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}
