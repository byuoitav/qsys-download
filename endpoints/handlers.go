package endpoints

import (
	"net/http"

	//"github.com/byuoitav/clevertouch-control/device/actions"
	"go.uber.org/zap"

	"fmt"

	"github.com/gin-gonic/gin"
)

func (d *DeviceManager) downloadFile(context *gin.Context) {
	fmt.Println("downloadFile Top")
	d.Log.Debug("downloading file", zap.String("file", context.Param("file")), zap.String("address", context.Param("address")))
	power := context.Param("file")
	fmt.Println(power)
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
