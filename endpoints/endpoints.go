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

//Info needed for upload:
//BOX API key - store as env or file
//BOX Folder ID - store as env or file
//File Name for BOX - get from QSC
//Q-Sys File Name and path - get from QSC
//Q-Sys Core IP Address - get from QSC

func (d *DeviceManager) RunHTTPServer(router *gin.Engine, port string) error {
	d.Log.Info("registering http endpoints")

	// action endpoints
	route := router.Group("/api/v1")
	route.PUT("/:address/download/:file", d.downloadFile)

	server := &http.Server{
		Addr:           port,
		MaxHeaderBytes: 1021 * 10,
	}

	d.Log.Info("running http server")
	router.Run(server.Addr)

	return fmt.Errorf("http server stopped")
}
