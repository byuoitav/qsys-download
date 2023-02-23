package endpoints

import (
	"fmt"
	"net/http"
	"os"

	qsc "github.com/byuoitav/qsc-download/qscdownload"
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
	fmt.Println("downloadFile Top")
	d.Log.Debug("downloading file", zap.String("file", context.Param("file")), zap.String("address", context.Param("address")))

	power := context.Param("file")
	address := context.Param("address")
	fmt.Println(power, address)
	fmt.Println((context.Query("key")))          //This is how to get keys in the URL after the ?
	fmt.Println((context.PostForm("form data"))) //This is how to get keys as form data

	//context.String(http.StatusOK, "Some String")

	filename := "DefLeppardPourSomeSugarOnMelyrics.mp3"
	filepath := "tmp_audio/" + filename

	coreIP := "10.5.34.167"
	url := coreIP + "/api/v0/cores/self/media/Audio/" + filename

	qsc.DownloadFile(filepath, url)
	// if err != nil {
	// 	d.Log.Warn("could not download file", zap.Error(err))
	// 	context.JSON(http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// err := actions.SetPower(context, context.Param("address"), power)
	// if err != nil {
	// 	d.Log.Warn("failed to set power", zap.Error(err))
	// 	context.JSON(http.StatusInternalServerError, err.Error())
	// 	return
	//}

	d.Log.Debug("successfully set power", zap.String("power", context.Param("power")), zap.String("address", context.Param("address")))
	context.JSON(http.StatusOK, 1)
}

func (d *DeviceManager) getDownloadStatus(context *gin.Context) {
	d.Log.Debug("downloading file", zap.String("file", context.Param("file")), zap.String("address", context.Param("address")))
	//power := context.Param("file")
	// if err != nil {
	// 	d.Log.Warn("could not download file", zap.Error(err))
	// 	context.JSON(http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// err := actions.SetPower(context, context.Param("address"), power)
	// if err != nil {
	// 	d.Log.Warn("failed to set power", zap.Error(err))
	// 	context.JSON(http.StatusInternalServerError, err.Error())
	// 	return
	//}

	d.Log.Debug("successfully set power", zap.String("power", context.Param("power")), zap.String("address", context.Param("address")))
	context.JSON(http.StatusOK, 1)
}
