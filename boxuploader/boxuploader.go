package boxuploader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FolderInfo struct {
	Entries []Entries `json:"entries"`
}
type Entries struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type BoxErrorCheck struct {
	TheType string `json:"type"`
}

//Todo: Check responses for authorization and return errors appropriately.

func CheckForFolder(room string, token string, parentFolderID string) (folderID string, err error) {
	fmt.Println("boxuploader - CheckForFolder start")

	url := "https://api.box.com/2.0/folders/" + parentFolderID + "/items"

	bearer := "Bearer " + token
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error on response.\n[ERROR] -", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading the response bytes:", err)
		return "", err
	}

	var folder FolderInfo
	err = json.Unmarshal(body, &folder)
	//todo error check validation

	folderID = ""
	for i := range folder.Entries {
		remoteFolder := folder.Entries[i].Name

		fmt.Println(folder.Entries[i].Name)
		if remoteFolder == room {
			folderID = folder.Entries[i].ID
		}
	}

	return folderID, nil
}

func CreateFolder(room string, token string, parentFolderID string) (folderID string, err error) {
	url := "https://api.box.com/2.0/folders"
	method := "POST"

	name := `"name": "` + room + `"`
	parent := `"parent":{ "id":"` + parentFolderID + `"}`

	payload := strings.NewReader(`{
	  ` + name + `,
	  ` + parent + `
	  	}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		//fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//fmt.Println(err)
		return
	}

	//Update folder id from the JSON return body
	var errorCheck Entries
	err = json.Unmarshal(body, &errorCheck)
	//error response validation
	if errorCheck.ID == "error" {
		return "", err
	}

	folderID = errorCheck.ID
	return folderID, nil
}

func UploadFile(filename string, localfilepath string, token string, folderID string, parentFolderID string) (uploaded bool, err error) {

	url := "https://upload.box.com/api/2.0/files/content"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("attributes", `{"name":"`+filename+`", "parent":{"id":"`+folderID+`"}}`)
	file, errFile2 := os.Open(localfilepath)
	if errFile2 != nil {
		fmt.Println(errFile2)
		return false, errFile2
	}

	part2, errFile2 := writer.CreateFormFile("file", filepath.Base(localfilepath))
	if errFile2 != nil {
		fmt.Println(errFile2)
		return false, errFile2
	}

	_, errFile2 = io.Copy(part2, file)
	if errFile2 != nil {
		fmt.Println(errFile2)
		return false, errFile2
	}

	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	defer file.Close()

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return false, err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	uploaded = true
	var errorCheck BoxErrorCheck
	err = json.Unmarshal(body, &errorCheck)
	//error response validation
	if errorCheck.TheType == "error" {
		uploaded = false
		return false, err
	}

	return uploaded, nil
}
