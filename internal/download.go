package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/binodsh/libgen"
	"github.com/cheggaaa/pb"
)

//Download downlods books from libgen mirror
func Download(book libgen.BookInfo) error {
	fmt.Println("downloading book - " + book.Title)
	downloadInfo, _ := libgen.GetDownloadInfo(book.ID)

	filename := strconv.FormatInt(book.ID, 10) + "-" + book.Title + "." + book.Extension

	req, err := http.NewRequest("GET", downloadInfo.DowloadLink, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept-Encoding", "*")
	req.Header.Add("referrer", downloadInfo.DownloadPageURL)

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if r.StatusCode == http.StatusOK {
		fileSize := r.ContentLength
		bar := pb.Full.Start64(fileSize)

		//save file to libgen folder under the current directory
		currentDir, _ := os.Getwd()
		currentDir += "/libgen"
		out, err := createFile(currentDir, filename)
		if err != nil {
			return err
		}

		_, err = io.Copy(out, bar.NewProxyReader(r.Body))
		if err != nil {
			return err
		}

		bar.Finish()
		out.Close()
		r.Body.Close()
		fmt.Println("Book saved to " + currentDir + "/" + filename)
	} else {
		return fmt.Errorf("unable to reach mirror %v: HTTP %v", req.Host, r.StatusCode)
	}

	return nil

}

func createFile(outputDir, filename string) (*os.File, error) {
	outputDir = strings.TrimSuffix(outputDir, "/")
	err := os.MkdirAll(outputDir, 0755)

	if err != nil {
		return nil, fmt.Errorf("error while creating file: " + err.Error())
	}

	var out *os.File
	out, err = os.Create(fmt.Sprintf("%s/%s", outputDir, filename))
	if err != nil {
		return nil, fmt.Errorf("error while creating file: " + err.Error())
	}

	return out, nil
}
