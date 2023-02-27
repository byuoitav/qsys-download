package qscdownload

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DownloadFile(filepath string, url string) (size int64, err error) {
	fmt.Println(filepath, url)

	//Check if directory exists and create it if not
	dirs := strings.Split(filepath, "/")

	exists, err := exists(dirs[0])
	fmt.Println(exists)

	if !exists {
		fmt.Println(dirs[0] + " folder does not exist, creating dir")
		if err := os.Mkdir(dirs[0], os.ModePerm); err != nil {
			fmt.Println("make dir failed")
			return 0, err
		}
	}

	// Download the file from the QSC core
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

	// Write the body to file
	size, e := io.Copy(out, resp.Body)
	if err != nil {
		return 0, e
	}

	return size, nil
}
