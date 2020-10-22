package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/controller"
	"github.com/wzru/gitran-server/middleware"
)

func projInit(g *gin.RouterGroup) {
	gg := g.Group("/projects")
	gg.GET("/:owner/:project", controller.GetProj)
	gg.Use(middleware.AuthJWT())
	{
		// gg.PUT("/:login", controller.UpdateProj)
		gg.POST("", controller.CreateProj)
	}
}
