package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/controller"
	"github.com/wzru/gitran-server/middleware"
)

func projInit(g *gin.RouterGroup) {
	gg := g.Group("/projects")
	gg.GET("/:owner", controller.ListProj)
	gg.GET("/:owner/:project", controller.GetProj)
	gg.GET("/:owner/:project/configs", controller.GetUserProjCfg)
	gg.Use(middleware.AuthUserJWT())
	{
		gg.POST("", controller.CreateUserProj)
	}
	gg.Use(middleware.AuthUserProjJWT())
	{
		gg.POST("/:owner/:project/configs", controller.CreateUserProjCfg)
		// gg.POST("/:owner/:project/configs/:config_id/rules", controller.CreateUserProjBrchRule)
		// gg.GET("/:owner/:project/branches", controller.ListUserProjBrch)
	}
}
