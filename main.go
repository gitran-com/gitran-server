package main

import (
	"flag"

	v1 "github.com/WangZhengru/gitran-be/api/v1"
	"github.com/WangZhengru/gitran-be/config"
	"github.com/WangZhengru/gitran-be/middleware"
	"github.com/WangZhengru/gitran-be/model"
	"github.com/gin-gonic/gin"
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
	api := r.Group(config.APP.APIPrefix + "/api")
	{
		apiv1 := api.Group("/v1")
		{
			v1.Init(apiv1)
		}
	}
	r.Run(config.APP.Addr)
	//fmt.Println("MAIN RETURN")
}
