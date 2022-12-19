package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/byuoitav/qsys-download/handlers"
)

func main() {
	port := ":8013"

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// action endpoints
	write := router.Group("/write")
	write.POST("/:address/power/:power", handlers.SetPower)
	write.POST("/:address/volume/:volume", handlers.SetVolume)
	write.POST("/:address/volume/mute/:mute", handlers.SetMute)
	write.POST("/:address/display/:blank", handlers.SetBlank)
	write.POST("/:address/input/:port", handlers.SetInput)

	// status endpoints
	read := router.Group("/read")
	read.GET("/:address/power", handlers.GetPower)
	read.GET("/:address/volume", handlers.GetVolume)
	read.GET("/:address/volume/mute", handlers.GetMute)
	read.GET("/:address/input", handlers.GetInput)

	server := &http.Server{
		Addr:           port,
		MaxHeaderBytes: 1021 * 10,
	}

	router.Run(server.Addr)
}

func test() {

	filename := "DefLeppardPourSomeSugarOnMelyrics.mp3"
	filepath := "tmp_audio/" + filename

	coreIP := "10.5.34.167"
	url := coreIP + "/api/v0/cores/self/media/Audio/" + filename

	downloadFile(filepath, url)
}
