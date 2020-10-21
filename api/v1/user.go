package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/controller"
	"github.com/wzru/gitran-server/middleware"
)

func userInit(g *gin.RouterGroup) {
	gg := g.Group("/users")
	gg.GET("/:login", controller.GetUser)
	gg.Use(middleware.AuthJWT())
	{
		gg.PUT("/:login", controller.UpdateUser)
		gg.POST("/projects", controller.CreateProj)
	}
}
