package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"scrooge/cache"
	"scrooge/config"
	"scrooge/postgres"
	"scrooge/server"
	"scrooge/telegram"
	"scrooge/utils"
)

var conn *pgxpool.Pool
var router *gin.Engine

func main() {
	config.InitConfig()
	cache.InitRedis()
	err := postgres.InitPool()
	if err != nil {
		utils.Error("Failed to connect to the database")
		panic(err)
	}

	go telegram.StartBot()

	router = server.InitRouter()
	server.InitRoutes(router)
	err = server.Run(router)
	if err != nil {
		utils.Error("Failed to start the server")
		panic(err)
	}
}
