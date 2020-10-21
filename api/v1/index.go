package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/controller"
	"github.com/wzru/gitran-server/middleware"
)

func Init(g *gin.RouterGroup) {
	g.GET("/ping", controller.PingV1)
	g.POST("/login", controller.Login)
	g.POST("/register", controller.Register)
	if config.Github.Enable {
		// g.POST("/auth/github", controller.Login)
		g.POST("/auth/github/callback", controller.AuthGithubCallback)
	}
	g.Group("/").Use(middleware.AuthJWT())
	{
		g.POST("/refresh", controller.Refresh)
		g.GET("/users/:name", controller.GetUser)
		g.PUT("/users/:name", controller.UpdateUser)
		g.POST("/projects", controller.CreateProj)
	}
}
