package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var (
	Log *logrus.Logger
)

func init() {
	Log = logrus.New()
	Log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.SetOutput(logFile)
	} else {
		Log.Error("无法打开日志文件 ", err)
	}

	// 设置日志输出到文件和控制台
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	Log.SetOutput(multiWriter)

	Log.SetLevel(logrus.DebugLevel)
}
