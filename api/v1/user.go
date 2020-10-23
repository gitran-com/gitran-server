package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/controller"
	"github.com/wzru/gitran-server/middleware"
)

func userInit(g *gin.RouterGroup) {
	gg := g.Group("/users")
	gg.GET("/:username", controller.GetUser)
	gg.Use(middleware.AuthJWT())
	{
		gg.PUT("/:username", controller.UpdateUser)
		gg.POST("/projects", controller.CreateUserProj)
	}
}
