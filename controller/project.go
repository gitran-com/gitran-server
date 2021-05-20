package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/constant"
	"github.com/wzru/gitran-server/model"
	"github.com/wzru/gitran-server/service"
)

func ListProj(ctx *gin.Context) {
	service.ListProj(ctx)
}

func GetProj(ctx *gin.Context) {
	service.GetProj(ctx)
}

func CreateUserProj(ctx *gin.Context) {
	service.CreateUserProj(ctx)
}

func CreateOrgProj(ctx *gin.Context) {
	service.CreateOrgProj(ctx)
}

func ListUserPubProj(ctx *gin.Context) {
	service.ListUserPubProj(ctx)
}

func ListAuthUserProj(ctx *gin.Context) {
	service.ListAuthUserProj(ctx)
}

func ListUserProjBrch(ctx *gin.Context) {
	proj := ctx.Keys["project"].(*model.Project)
	if proj.Type == constant.TypeGitURL {
		service.ListGitRepoBrch(ctx)
	} else if proj.Type == constant.TypeGithub {
		service.ListGithubRepoBrch(ctx)
	} else {
		//TODO
	}
	// service.ListUserProjBrch(ctx)
}
