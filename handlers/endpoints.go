package endpoints

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DeviceManager struct {
	Log *zap.Logger
}

func (d *DeviceManager) RunHTTPServer(router *gin.Engine, port string) error {
	d.Log.Info("registering http endpoints")

	// action endpoints
	route := router.Group("/api/v1")
	route.GET("/:address/download/:file", d.downloadFile)

	// status endpoints
	route.GET("/:address/download-status/:file", d.getDownloadStatus)

	server := &http.Server{
		Addr:           port,
		MaxHeaderBytes: 1021 * 10,
	}

	d.Log.Info("running http server")
	router.Run(server.Addr)

	return fmt.Errorf("http server stopped")
}
