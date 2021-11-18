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
	initFile(g)
	initTran(g)
	initVote(g)
	initPin(g)
	initBuild(g)
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
	gg.GET("/:user_id", controller.GetUser)
	gg.GET("/:user_id/projects", controller.ListUserProj)
}

func initProj(g *gin.RouterGroup) {
	gg := g.Group("/projects")
	gg.GET("/:uri", middleware.TryAuthUser(), middleware.MustGetProj(), middleware.MustGetProjRole(), controller.GetProj)
	gg.GET("/:uri/status", middleware.MustGetProj(), controller.ProjStatWS)
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

func initFile(g *gin.RouterGroup) {
	gg := g.Group("/projects/:uri/files", middleware.MustAuthUser(), middleware.MustAuthProjViewer())
	gg.GET("", controller.ListProjFiles)
	gg.GET("/:file_id", controller.GetProjFile)
}

func initTran(g *gin.RouterGroup) {
	gg := g.Group("/projects/:uri/translations", middleware.MustAuthUser(), middleware.MustAuthProjViewer())
	gg.GET("/:code/:sent_id", controller.ListSentTrans)
	gg.POST("/:code/:sent_id", middleware.MustAuthProjContributor(), controller.PostTran)
	gg.DELETE("/:tran_id", middleware.MustAuthProjCommitterOrTranCommitter(), controller.DelTran)
}

func initPin(g *gin.RouterGroup) {
	gg := g.Group("/projects/:uri/pins", middleware.MustAuthUser(), middleware.MustAuthProjCommitter())
	gg.POST("/:tran_id", controller.PinTran)
	gg.DELETE("/:tran_id", controller.UnpinTran)
}

func initVote(g *gin.RouterGroup) {
	gg := g.Group("/projects/:uri/votings", middleware.MustAuthUser(), middleware.MustAuthProjContributor())
	gg.POST("/:tran_id/likes", controller.AddLikes)
	gg.POST("/:tran_id/unlikes", controller.AddUnlikes)
	gg.DELETE("/:tran_id", controller.DelVote)
}

func initBuild(g *gin.RouterGroup) {
	gg := g.Group("/projects/:uri/builds", middleware.MustAuthUser(), middleware.MustAuthProjContributor())
	gg.POST("", controller.NewBuild)
}

func initLang(g *gin.RouterGroup) {
	gg := g.Group("/languages")
	gg.GET("", controller.ListLangs)
	gg.GET("/:code", controller.GetLang)
}
