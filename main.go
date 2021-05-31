package main

import (
	"github.com/gin-gonic/gin"
	"tesla/config"
	"tesla/controllers"
	"tesla/globalvar"
	"tesla/utils"
	_ "tesla/utils"
	"github.com/robfig/cron"
)

func main() {
	//init config instance
	router := gin.Default()
	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe

	//初始化
	globalvar.InitGlov()

	c := cron.New()
	err := c.AddFunc("*/1 * * * * *", func() {
		controllers.UploadToKafka()
	})
	if err != nil {
		utils.Log.WithField("err", err).Error("start error")
		return
	}
	err = c.AddFunc("*/30 * * * * *", func() {
		controllers.UploadWebLock()
	})
	if err != nil {
		utils.Log.WithField("err", err).Error("start error")
		return
	}
	err = c.AddFunc("* */59 * * * *", func() {
		controllers.UploadWebLock()
	})
	if err != nil {
		utils.Log.WithField("err", err).Error("start error")
		return
	}
	c.Start()

	router.GET("/api/auth", controllers.AuthController)
	router.GET("/api/traffic", controllers.TrafficController)
	router.POST("/api/kick", controllers.KickController)

	router.Run(":" + config.AppConfig.AppPort)

}