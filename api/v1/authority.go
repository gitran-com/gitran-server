package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/controller"
)

func authInit(g *gin.RouterGroup) {
	g.POST("/login", controller.Login)
	g.POST("/register", controller.Register)
	g.POST("/refresh", controller.Refresh)
}
