package endpoints

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	box "github.com/byuoitav/qsys-download/boxuploader"
	qsc "github.com/byuoitav/qsys-download/qscdownload"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (d *DeviceManager) getBoxAccessToken(context *gin.Context) (token string) {
	f, err := os.ReadFile("../box_api_key.yourmom")
	if err != nil {
		d.Log.Warn("could not open box_api_key.yourmom", zap.Error(err))
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	token = string(f)

	return token
}

func (d *DeviceManager) getBoxFolderID(context *gin.Context) (id string) {
	f, err := os.ReadFile("../box_folder_id.yourmom")
	if err != nil {
		d.Log.Warn("could not open box_folder_id.yourmom", zap.Error(err))
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	id = string(f)

	return id
}

func (d *DeviceManager) downloadFile(context *gin.Context) {
	//Download from QSC
	filename := context.Param("file")
	downloadfilepath := context.PostForm("filePath")
	localfilepath := "../tmp_audio/" + filename

	room := context.PostForm("room")
	coreIP := context.Param("address")

	if filename == "" || downloadfilepath == "" || coreIP == "" {
		context.String(http.StatusOK, "Malformed query. URL should match https://api/v1/core IP address/download/file name.mp3. FilePath must be defined in form data as well.")
	}

	url := "http://" + coreIP + "/api/v0/cores/self/media/" + downloadfilepath

	d.Log.Debug("downloading file from Q-Sys ", zap.String("Storing file: ", localfilepath), zap.String("Core address: ", coreIP), zap.String("Download url: ", url))

	size, err := qsc.DownloadFile(localfilepath, url)
	if err != nil {
		d.Log.Warn("could not download file from Q-Sys", zap.Error(err))
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	d.Log.Debug("Downloaded from Q-Sys " + strconv.FormatInt(int64(size), 10) + " Bytes")

	//Upload to Box
	d.Log.Debug("Uploading to Box ")

	token := d.getBoxAccessToken(context)
	parentFolderID := d.getBoxFolderID(context)

	fmt.Println(token, parentFolderID)
	var folderID string
	folderID, err = box.CheckForFolder(room, token, parentFolderID)
	if err != nil {
		d.Log.Warn("box folder check failed: ", zap.Error(err))
		context.String(http.StatusInternalServerError, err.Error())
		return
	}
	if folderID == "" {
		folderID, err = box.CreateFolder(room, token, parentFolderID)
		if err != nil {
			d.Log.Warn("box create folder failed: ", zap.Error(err))
			context.String(http.StatusInternalServerError, err.Error())
			return
		}
	}
	var uploaded bool
	uploaded, err = box.UploadFile(filename, localfilepath, token, folderID, parentFolderID)
	if err != nil {
		d.Log.Warn("box upload file failed: ", zap.Error(err))
		context.String(http.StatusInternalServerError, err.Error())
		return
	}

	if uploaded {
		//delete local file
		err = os.Remove(localfilepath)
		if err != nil {
			d.Log.Warn("file could not be deleted. file: "+localfilepath, zap.Error(err))
			context.String(http.StatusInternalServerError, err.Error())
		}
	}
	context.String(http.StatusOK, "Downloaded file from QSC and uploaded to Box successful. File size:  "+strconv.FormatInt(int64(size), 10)+" Bytes")
}
