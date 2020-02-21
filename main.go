package main

import (
	"fmt"

	"github.com/binodsh/libgen"
)

func main() {
	fmt.Println("hello world")

	books := libgen.SearchBookByTitle("nepal")

	for _, book := range books {
		fmt.Printf("Author: %s\nTitle: %s\nDownload Link: %s\n\n", book.Author, book.Title, book.DownloadLink)
	}
}
