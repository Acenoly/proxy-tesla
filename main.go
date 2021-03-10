package main

import (
	"github.com/gin-gonic/gin"
	"tesla/config"
	"tesla/controllers"
	_ "tesla/utils"
)

func main() {
	//init config instance
	router := gin.Default()
	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe

	router.GET("/api/auth", controllers.AuthController)
	router.GET("/api/traffic", controllers.TrafficController)
	router.POST("/api/kick", controllers.KickController)

	router.Run(":" + config.AppConfig.AppPort)

}