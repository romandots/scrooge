package server

import (
	"github.com/gin-gonic/gin"
	"scrooge/config"
	"scrooge/utils"
)

type Connection interface{}

func InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	gin.SetMode(gin.DebugMode)

	return router
}

func InitRoutes(router *gin.Engine) {
	router.GET("/", welcome)
}

func Run(router *gin.Engine) error {
	err := router.Run(config.Socket)
	if err == nil {
		utils.Info("Server is running on port " + config.Socket)
	}

	return err
}
