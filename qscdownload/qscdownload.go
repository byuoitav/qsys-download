package qscdownload

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Baerer struct {
	Token string `json:"token"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type QSCresponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

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

func login(coreIP string) (baerer string, err error) {
	loginInfo, err := json.Marshal(getCredentials())
	fmt.Println("Logging in with: ", string(loginInfo))
	if err != nil {
		return "", err
	}
	payload := strings.NewReader(string(loginInfo))
	//fmt.Println(loginInfo, payload)
	url := "http://" + coreIP + "/api/v0/logon"

	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var token Baerer
	err = json.Unmarshal(body, &token)
	if err != nil {
		fmt.Println(string(body), err)
		return "", err
	}

	//todo: error check for bad username/password
	fmt.Println("Loggin Complete. Bearer Token: ", token.Token)
	return token.Token, nil
}

func logout(coreIP string, token string) (err error) {
	fmt.Println("Logging out Bearer: ", token)
	url := "http://" + coreIP + "/api/v0/logon"
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	baerer := "Baerer " + token
	fmt.Println(token)
	req.Header.Add("Authorization", baerer)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Completed logging out Bearer: ", token)
	fmt.Println("Logout response: ", string(body))
	return nil
}

func getCredentials() (credentials Login) {
	f, err := os.ReadFile("config/qsc_login.cfg")
	if err != nil {
		fmt.Println("Error openiong qsc_login.cfg. Proceeding without QSC login credentials")
		return
	}
	err = json.Unmarshal(f, &credentials)
	if err != nil {
		fmt.Println("Error unmarshaling qsc_login.cfg. Please ensure file is in json format")
		return
	}
	return credentials
}

func DownloadFile(filepath string, coreIP string, downloadFilePath string) (size int64, err error) {
	fmt.Println(filepath, coreIP, downloadFilePath)

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
	url := "http://" + coreIP + "/api/v0/cores/self/media/" + downloadFilePath
	fmt.Println("Attermpting to login without credentials")
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	fmt.Println("Attermpting to login without credentials response: ", resp)
	defer resp.Body.Close()

	//Login if needed
	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Attermpting to login with credentials")
		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		token, err := login(coreIP)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		bearer := "Bearer " + token
		req.Header.Add("Authorization", bearer)
		req.Header.Add("Accept", "application/octet-stream")

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		defer resp.Body.Close()
		fmt.Println("Attempting to login with credentials response: ", resp)
		logout(coreIP, token)

	}

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	// Write the body to file
	size, e := io.Copy(out, resp.Body)
	if err != nil {
		return 0, e
	}

	return size, nil
}
