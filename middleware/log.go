package middleware

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/gitran-com/gitran-server/config"
)

func Logger() gin.HandlerFunc {
	// 写入文件
	src, err := os.OpenFile(config.GinLogFile+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Warnf("Cannot open gin log file '%s' : %v", config.GinLogFile, err)
	}
	// 实例化
	logger := logrus.New()
	//设置输出
	logger.Out = src
	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		config.GinLogFile+".%Y%m%d.log",
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(config.GinLogFile+".log"),
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	logger.AddHook(lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: config.TimeFormat,
	}))

	return func(c *gin.Context) {
		//开始时间
		startTime := time.Now()
		//处理请求
		c.Next()
		//结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		//请求方式
		reqMethod := c.Request.Method
		//请求路由
		reqUrl := c.Request.RequestURI
		//状态码
		statusCode := c.Writer.Status()
		//请求ip
		clientIP := c.ClientIP()
		// 日志格式
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUrl,
		}).Info()
	}
}
