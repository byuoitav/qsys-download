package qscdownload

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(filepath string, url string) (size int64, err error) {
	//Todo: Check for file path folder exists and create it if not on system
	fmt.Println(filepath, url)
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	size, e := io.Copy(out, resp.Body)
	if err != nil {
		return 0, e
	}

	return size, nil
}
