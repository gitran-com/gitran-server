package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/controller"
	"github.com/wzru/gitran-server/middleware"
)

//Init 初始化路由
func Init(g *gin.RouterGroup) {
	g.GET("/ping", controller.PingV1)
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
		//Project Config
		gg.GET("/:id/configs", controller.ListUserProjCfg)
		gg.POST("/:id/configs", controller.CreateUserProjCfg)
		gg.PUT("/:id/configs", controller.SaveUserProjCfg)
		// gg.GET("/:owner/:project/configs/:config_id", controller.GetUserProjCfg)

		//Branch Rule
		gg.GET("/:id/configs/:config_id/rules", controller.ListUserProjBrchRule)
		gg.POST("/:id/configs/:config_id/rules", controller.CreateUserProjBrchRule)
		// gg.PUT("/:owner/:project/configs/:config_id/rules", controller.SaveUserProjBrchRule)
		// gg.GET("/:owner/:project/configs/:config_id/rules/:rule_id", controller.GetUserProjBrchRule)
		// gg.GET("/:owner/:project/branches", controller.ListUserProjBrch)
	}
}

func initLang(g *gin.RouterGroup) {
	gg := g.Group("/languages")
	gg.GET("", controller.GetLangs)
	gg.GET("/:language_id", controller.GetLang)
}
