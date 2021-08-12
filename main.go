package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	v1 "github.com/gitran-com/gitran-server/api/v1"
	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/middleware"
	"github.com/gitran-com/gitran-server/model"
	"github.com/gitran-com/gitran-server/service"
)

func main() {
	//解析命令行参数
	flag.Parse()
	//读取解析配置文件
	if err := config.Init(); err != nil {
		return
	}
	//数据库&语言文件初始化
	if err := model.Init(); err != nil {
		return
	}
	//服务初始化
	if err := service.Init(); err != nil {
		return
	}
	//API初始化
	var r *gin.Engine
	if config.IsDebug {
		r = gin.Default()
		r.Use(debugCORS())
	} else {
		log.Infof("writing log in file...")
		r = gin.New()
		r.Use(middleware.Logger())
		r.Use(cors.Default())
	}
	r.Use(sessions.Sessions(config.APP.Name, cookie.NewStore([]byte(config.APP.SessionSecret))))
	api := r.Group(config.APP.APIPrefix + "/api")
	{
		apiv1 := api.Group("/v1")
		v1.Init(apiv1)
	}
	if err := r.Run(config.APP.Addr); err != nil {
		log.Fatalf("server run error : %v\n", err.Error())
	}
}

func debugCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, DELETE, PUT, PATCH")
		c.Writer.Header().Set("Content-Type", "*; charset=utf-8")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
