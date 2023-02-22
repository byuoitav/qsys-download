package main

import (
	"net/http"

	"github.com/byuoitav/qsys-download/endpoints"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
)

func main() {
	var port, logLevel string
	pflag.StringVarP(&port, "port", "p", "8013", "port for microservice to av-api communication")
	pflag.StringVarP(&logLevel, "log", "l", "Info", "Initial log level")
	pflag.Parse()

	port = ":" + port

	manager := endpoints.DeviceManager{
		Log: buildLogger(logLevel),
	}

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "good",
		})
	})

	manager.RunHTTPServer(router, port)
}

// func test() {

// 	filename := "DefLeppardPourSomeSugarOnMelyrics.mp3"
// 	filepath := "tmp_audio/" + filename

// 	coreIP := "10.5.34.167"
// 	url := coreIP + "/api/v0/cores/self/media/Audio/" + filename

// 	downloadFile(filepath, url)
// }
