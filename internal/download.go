package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

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

		out, err2 := createOutputFile(filename)
		if err2 != nil {
			return err2
		}

		_, err := io.Copy(out, bar.NewProxyReader(r.Body))
		if err != nil {
			return err
		}

		bar.Finish()
		out.Close()
		r.Body.Close()
	} else {
		return fmt.Errorf("unable to reach mirror %v: HTTP %v", req.Host, r.StatusCode)
	}

	return nil

}

func createOutputFile(filename string) (*os.File, error) {
	var mkErr error
	var out *os.File

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if stat, err := os.Stat(fmt.Sprintf("%s/libgen", wd)); err == nil && stat.IsDir() {
		out, mkErr = os.Create(fmt.Sprintf("%s/libgen/%s", wd, filename))
	} else {
		if err := os.Mkdir(fmt.Sprintf("%s/libgen", wd), 0755); err != nil {
			return nil, err
		}
		out, mkErr = os.Create(fmt.Sprintf("%s/libgen/%s", wd, filename))
	}
	if mkErr != nil {
		return nil, mkErr
	}

	return out, nil
}
