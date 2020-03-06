package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
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
	Args:  cobra.MinimumNArgs(1),
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

		//show spinner while searching the book
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		s.Suffix = " searching for '" + searchStr + "'"
		books, err := libgen.Search(searchOptions)
		s.Stop()

		if err != nil {
			fmt.Println("error while searching book")
			return
		}

		if len(books) == 0 {
			fmt.Println("Sorry, no books found")
			return
		}

		//promtui template
		templates := &promptui.SelectTemplates{
			Label:    "{{ . | red | bold}}",
			Active:   "➡ {{ .Title | cyan | bold }}",
			Inactive: "  {{ .Title | cyan | faint}}",
			Selected: "➡️ {{ .Title | red | cyan }}",
			Details: `
				--------- Book Details ----------
				{{ "Title:" | faint }} {{if gt (len .Title) 60}} {{ (slice .Title 0 60) }}.... {{else}} {{.Title}} {{end}}
				{{ "@author:" | faint }} {{if gt (len .Author) 60}} {{ (slice .Author 0 60) }}.... {{else}} {{.Author}} {{end}}    {{ "@publisher:" | faint }} {{if gt (len .Publisher) 60}} {{ (slice .Publisher 0 60) }}.... {{else}} {{.Publisher}} {{end}}    {{ "@year:" | faint }} {{ .Year }}
				{{ "@extension:" | faint }} {{ .Extension }}    {{ "@pages:" | faint }} {{ .Pages }}`,
			Help: `{{ "Use the arrow keys to navigate: ↓ ↑ → ← & Hit Enter to download the book." | faint}}`,
		}

		prompt := promptui.Select{
			Label:        "Select a book to download",
			Items:        books,
			Templates:    templates,
			HideSelected: true,
			Size:         10,
		}

		i, _, err := prompt.Run()

		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		err = internal.Download(books[i])
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}
