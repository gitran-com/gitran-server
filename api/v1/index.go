package v1

import (
	"github.com/WangZhengru/gitran-be/controller"
	"github.com/gin-gonic/gin"
)

func Init(g *gin.RouterGroup) {
	g.POST("/login", controller.Login)
	g.POST("/register", controller.Register)
}
