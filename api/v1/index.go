package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/controller"
)

func Init(g *gin.RouterGroup) {
	g.GET("/ping", controller.PingV1)
	authInit(g)
	userInit(g)
}
