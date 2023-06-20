package main

import (
	"jwt-authentication-golang/controllers"
	"jwt-authentication-golang/database"
	"jwt-authentication-golang/middlewares"
	"jwt-authentication-golang/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Initialize Database
	database.Connect(config.DBConStr)
	database.Migrate()

	// Initialize Router
	router := initRouter()
	//router.Run(":8080")
	router.Run(config.ServerAddress)
}

func initRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/token", controllers.GenerateToken)
		api.POST("/user/register", controllers.RegisterUser)
		secured := api.Group("/secured").Use(middlewares.Auth())
		{
			secured.GET("/ping", controllers.Ping)
			secured.POST("/chat", controllers.Call)
		}
	}
	return router
}
