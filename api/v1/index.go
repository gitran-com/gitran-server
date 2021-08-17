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
	initUsers(g)
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
		//GitHub登录回调
		gg.GET("/github/login", controller.AuthGithubLogin)
		//GitHub引入repo回调
		gg.GET("/github/import", middleware.MustAuthUser(), controller.AuthGithubImport)
		//新注册GitHub用户
		gg.POST("/github/new", middleware.MustAuthNewGithubUser(), controller.NewGithubUser)
		//获得所有GitHub仓库
		gg.GET("/github/repos", middleware.MustAuthUser(), controller.GetGithubRepos)
	}
}

func initUser(g *gin.RouterGroup) {
	gg := g.Group("/user")
	gg.Use(middleware.MustAuthUser())
	{
		gg.GET("", controller.GetMe)
		gg.PUT("", controller.EditMe)
		gg.GET("/projects", controller.GetMyProjects)
	}
}

func initUsers(g *gin.RouterGroup) {
	gg := g.Group("/users")
	gg.GET("/:id", controller.GetUser)
	gg.GET("/:id/projects", controller.ListUserProj)
}

func initProj(g *gin.RouterGroup) {
	gg := g.Group("/projects")
	gg.GET("/:uri", middleware.TryAuthUser(), middleware.MustGetProjRole(), controller.GetProj)
	gg.Use(middleware.MustAuthUser())
	{
		gg.POST("", controller.CreateUserProj)
		gg.GET("/:uri/existed", controller.ProjExisted)
	}
	gg.Use(middleware.MustAuthProjAdmin())
	{
		//Repo Branch
		gg.GET("/:uri/branches", controller.ListProjBrch)
		//Project Config
		gg.GET("/:uri/config", controller.GetProjCfg)
		gg.PUT("/:uri/config", controller.UpdateProjCfg)
		gg.GET("/:uri/config/preview", controller.PreviewProjCfg)
	}
}

func initLang(g *gin.RouterGroup) {
	gg := g.Group("/languages")
	gg.GET("", controller.ListLangs)
	gg.GET("/:code", controller.GetLang)
}
