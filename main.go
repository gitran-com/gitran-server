package main

import (
	"flag"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "github.com/wzru/gitran-server/api/v1"
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/middleware"
	"github.com/wzru/gitran-server/model"
)

func main() {
	//解析命令行参数
	flag.Parse()
	//读取解析配置文件
	if err := config.Init(); err != nil {
		return
	}
	//数据库初始化
	if err := model.Init(); err != nil {
		return
	}
	//defer model.DB.Close()//无需Close
	//API初始化
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(cors.Default())
	api := r.Group(config.APP.APIPrefix + "/api")
	{
		apiv1 := api.Group("/v1")
		{
			v1.Init(apiv1)
		}
	}
	if err := r.Run(config.APP.Addr); err != nil {
		log.Fatalf("Server run error : %v\n", err)
	}
}
