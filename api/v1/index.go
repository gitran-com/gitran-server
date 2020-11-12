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
	initGithub(g)
}

func initGithub(g *gin.RouterGroup) {
	gg := g.Group("/github")
	if config.Github.Enable {
		gg.Use(middleware.AuthUserJWT(), middleware.AuthUserGithubJWT())
		{
			gg.GET("/repos", controller.ListGithubRepo)
		}
	}
}

func initAuth(g *gin.RouterGroup) {
	gg := g.Group("/auth")
	gg.POST("/login", controller.Login)
	gg.POST("/register", controller.Register)
	gg.POST("/refresh", controller.Refresh)
	if config.Github.Enable {
		gg.POST("/github/login", controller.AuthGithubLogin)
		gg.POST("/github/register", controller.AuthGithubRegister)
		gg.POST("/github/login/callback", controller.AuthGithubLoginCallback)
		gg.Use(middleware.AuthUserJWT())
		{
			gg.POST("/github/bind", controller.AuthGithubBind)
			gg.POST("/github/bind/callback", controller.AuthGithubBindCallback)
			gg.POST("/github/import", controller.AuthGithubImport)
			gg.POST("/github/import/callback", controller.AuthGithubImportCallback)
		}
	}
}

func initLang(g *gin.RouterGroup) {
	gg := g.Group("/languages")
	gg.GET("", controller.GetLangs)
	gg.GET("/:language_id", controller.GetLang)
}

func initProj(g *gin.RouterGroup) {
	gg := g.Group("/projects")
	gg.GET("/:owner/:project", controller.GetProj) //获取特定项目
	gg.Use(middleware.AuthUserJWT())
	{
		gg.POST("", controller.CreateUserProj) //新建用户项目
	}
	gg.Use(middleware.AuthUserJWT(), middleware.AuthUserProjJWT())
	{
		//Project Config
		gg.GET("/:owner/:project/configs", controller.ListUserProjCfg)
		gg.POST("/:owner/:project/configs", controller.CreateUserProjCfg)
		gg.PUT("/:owner/:project/configs", controller.SaveUserProjCfg)
		//Branch Rule
		gg.GET("/:owner/:project/configs/:config_id/rules", controller.ListUserProjBrchRule)
		gg.POST("/:owner/:project/configs/:config_id/rules", controller.CreateUserProjBrchRule)
	}
}

func initUser(g *gin.RouterGroup) {
	pubGrp := g.Group("/users")
	pubGrp.Use(middleware.GetUser())
	{
		pubGrp.GET("/:username", controller.GetUser)
		pubGrp.GET("/:username/projects", controller.ListUserPubProj)
	}
	prvGrp := g.Group("/user")
	prvGrp.Use(middleware.AuthUserJWT())
	{

		prvGrp.GET("/projects", controller.ListAuthUserProj)
		prvGrp.POST("/projects", controller.CreateUserProj)
		// prvGrp.PUT("/:username", controller.UpdateUser)
	}
}
