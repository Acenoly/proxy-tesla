package utils

import (
	"github.com/sirupsen/logrus"
	"os"
	"tesla/config"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	Log.SetFormatter(&logrus.JSONFormatter{})
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	file, _ := os.OpenFile(config.AppConfig.LogPath, os.O_CREATE|os.O_WRONLY, 0666)
	Log.SetOutput(file)
	//设置最低loglevel
	Log.SetLevel(logrus.InfoLevel)
}
