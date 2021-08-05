package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/controller"
	"github.com/gitran-com/gitran-server/middleware"
)

//Init 初始化路由
func Init(g *gin.RouterGroup) {
	g.GET("/ping", controller.PingV1)
	g.GET("/test", controller.Test)
	initAuth(g)
	initUser(g)
	initProj(g)
	initLang(g)
}

func initAuth(g *gin.RouterGroup) {
	gg := g.Group("/auth")
	gg.POST("/login", controller.Login)
	gg.POST("/register", controller.Register)
	gg.POST("/refresh", controller.Refresh)
	if config.Github.Enable {
		gg.GET("/github", controller.AuthGithub)
		gg.GET("/github/login", controller.AuthGithubLogin)
		gg.GET("/github/import", middleware.AuthUserJWT(), controller.AuthGithubImport)
		gg.GET("/github/tokens", middleware.AuthUserJWT(), controller.GetGithubTokens)
		gg.GET("/github/repos/:id", middleware.AuthUserJWT(), controller.GetGithubRepos)
	}
}

func initUser(g *gin.RouterGroup) {
	gg := g.Group("/users")
	gg.GET("/:id", controller.GetUser)
	gg.GET("/:id/projects", controller.ListUserProj)
	gg.Use(middleware.AuthUserJWT())
	{
		// gg.PUT("/:username", controller.UpdateUser)
	}
}

func initProj(g *gin.RouterGroup) {
	gg := g.Group("/projects")
	gg.GET("/:id", controller.GetProj)
	gg.Use(middleware.AuthUserJWT())
	{
		gg.POST("", controller.CreateUserProj)
	}
	gg.Use(middleware.AuthUserProjJWT())
	{
		//Repo Branch
		gg.GET("/:id/branches", controller.ListProjBrch)

		//Project Config
		gg.GET("/:id/configs", controller.ListUserProjCfg)
		gg.POST("/:id/configs", controller.CreateUserProjCfg)
		gg.POST("/:id/configs/save", controller.SaveUserProjCfg)

		//Branch Rule
		gg.GET("/:id/configs/:config_id/rules", controller.ListUserProjBrchRule)
		gg.POST("/:id/configs/:config_id/rules", controller.CreateUserProjBrchRule)
	}
}

func initLang(g *gin.RouterGroup) {
	gg := g.Group("/languages")
	gg.GET("", controller.GetLangs)
	gg.GET("/:id", controller.GetLang)
}
