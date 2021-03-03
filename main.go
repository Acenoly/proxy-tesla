package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"tesla/config"
	"tesla/controllers"
	_ "tesla/utils"
)

func main() {

	// Disable Console Color, you don't need console color when writing the logs to file.
	// 禁用控制台日志颜色,日志写到文件的时候,不需要打开控制台日志颜色
	gin.DisableConsoleColor()
	// Logging to a file.  新建日志文件,得到文件结构,文件结构实现了写出器Writer接口
	f, _ := os.Create("/data/gin.log")
	//io.MultiWriter(多写出器方法)创建一个写出器, 将传入的多个写出器追加为一个写出器数组, 得到的写出器实现了Writer接口, 它会将需要写出的数据写出到每个写出器, 就像Unix命令tee,会将数据写入文件的同时打印到标准输出
	//配置Gin默认日志写出器为得到的多写出器
	gin.DefaultWriter = io.MultiWriter(f)
	// Use the following code if you need to write the logs to file and console at the same time.
	// 使用下面的代码,将日志写入文件的同时,也输出到控制台
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	//init config instance
	router := gin.Default()
	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe

	router.GET("/api/auth", controllers.AuthController)

	router.GET("/api/traffic", controllers.TrafficController)

	router.Run(":" + config.AppConfig.AppPort)

}
